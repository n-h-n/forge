package main

// Go template
const goTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/go
GO_VERSION ?= go1.25.1
GO = env GOENV=$(GO_ENV) GOROOT=$(GO_ROOT) $(GO_ROOT)/bin/go
go_dep_loc = $(shell $(GO) mod download $(1) && $(GO) list -m -f "{{.Dir}}" $(1))
go_nhn_dep_loc = $(call go_dep_loc,github.com/n-h-n/$(1))
go_dep_version = \
  $(shell 2>/dev/null $(GO) mod download $(1) && $(GO) list -m \
    -f "{{if .Replace}}{{.Replace.Version}}{{else}}{{.Version}}{{end}}" $(1) || echo DUMMY)
go_nhn_dep_version = $(call go_dep_version,github.com/n-h-n/$(1))
go_nhn_proto_sym = \
  mkdir -p $(dir $(TMP_DIR)/proto/nhn/$1) \
  && ln -n -f -s '$(call go_nhn_dep_loc,$1)' $(TMP_DIR)/proto/nhn/$1
GO_ENV = $(abspath $(__go_dir)/env)
GO_ROOT = $(realpath $(__go_dir)/go)
__go_dir = $(TOOLS_DIR)/go/$(GO_VERSION)
__go_env = $(__go_dir)/env
__go_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]")
__go_arch = $(shell uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
__go_mod_cache = $(GLOBAL_CACHE_DIR)/go/mod
.d.go: $(__go_env)
.PHONY: .d.go
$(__go_mod_cache):
	mkdir -p $@
$(__go_dir): | $(__go_mod_cache)
	rm -f $(GLOBAL_TMP_DIR)/go.tar.gz $(GLOBAL_TMP_DIR)/go.tar.gz.sha256
	curl \
	  -s \
	  -L \
	  -o $(GLOBAL_TMP_DIR)/go.tar.gz \
	  https://go.dev/dl/$(GO_VERSION).$(__go_kernel)-$(__go_arch).tar.gz
	rm -rf $(GLOBAL_TMP_DIR)/gotmp
	mkdir -p $(GLOBAL_TMP_DIR)/gotmp
	tar -xzf $(GLOBAL_TMP_DIR)/go.tar.gz -C $(GLOBAL_TMP_DIR)/gotmp
	rm -rf $(GLOBAL_TMP_DIR)/go.tar.gz
	mkdir $(GLOBAL_TMP_DIR)/gotmp/packages
	mkdir -p $(@D)
	mv $(GLOBAL_TMP_DIR)/gotmp $@
$(__go_env): | $(__go_dir)
	echo "GOROOT=$(realpath $(__go_dir)/go)" > $@
	echo "GOPATH=$(realpath $(__go_dir)/packages)" >> $@
	echo "GOMODCACHE=$(realpath $(__go_mod_cache))" >> $@
	echo "GOTOOLCHAIN=local" >> $@

`

// jq template
const jqTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/jq
JQ_VERSION ?= jq-1.7.1
JQ = $(realpath $(__jq_bin))
__jq_dir = $(TOOLS_DIR)/jq/$(JQ_VERSION)
__jq_bin = $(__jq_dir)/jq
__jq_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]" | sed -e 's/darwin/macos/')
__jq_arch = $(shell uname -m | sed -e 's/x86_64/amd64/')
__jq_file = jq-$(__jq_kernel)-$(__jq_arch)
__jq_base_url = https://github.com/jqlang/jq/releases/download/$(JQ_VERSION)
__jq_bin_url = $(__jq_base_url)/$(__jq_file)
__jq_checksum_url = $(__jq_base_url)/sha256sum.txt
.d.jq: $(__jq_bin)
.PHONY: .d.jq
$(__jq_bin):
	rm -f $(GLOBAL_TMP_DIR)/$(__jq_file) $(GLOBAL_TMP_DIR)/$(__jq_file).sha256
	curl -s -L -o $(GLOBAL_TMP_DIR)/$(__jq_file) $(__jq_bin_url)
	curl -s -L -o $(GLOBAL_TMP_DIR)/sha256sum.txt $(__jq_checksum_url)
	cd $(GLOBAL_TMP_DIR) && grep -F $(__jq_file) sha256sum.txt | $(SHA256SUM) -c -
	chmod u+x $(GLOBAL_TMP_DIR)/$(__jq_file)
	rm -rf $(__jq_dir)
	mkdir -p $(__jq_dir)
	cp $(GLOBAL_TMP_DIR)/$(__jq_file) $@

`

// Snyk template
const snykTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/snyk
SNYK_VERSION ?= v1.1295.0
SNYK_KILL_SWITCH ?= disabled
SNYK = \
  env \
  GOENV="$(GO_ENV)" \
  GOROOT="$(GO_ROOT)" \
  PATH="$(GO_ROOT)/bin:$(PATH)" \
  $(realpath $(__snyk_bin))
SNYK_TEST = \
  rm -f $(TMP_DIR)/snyk.json \
  && $(SNYK) test --json-file-output=$(TMP_DIR)/snyk.json
SNYK_TEST_FILTER = \
  $(JQ) \
    '[.vulnerabilities[] | select(.fixedIn != [])] | length' \
    $(TMP_DIR)/snyk.json \
  | xargs test 0 -eq
ifeq ($(SNYK_KILL_SWITCH),enabled)
  SNYK = echo -- snyk
endif
__snyk_bin = $(TOOLS_DIR)/snyk/$(SNYK_VERSION)/snyk
__snyk_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]" | sed -e 's/darwin/macos/')
__snyk_arch = $(shell uname -m | sed -e 's/x86_64//' -e 's/arm64/-arm64/')
__snyk_artifact = snyk-$(__snyk_kernel)$(__snyk_arch)
__snyk_artifact_sha = $(__snyk_artifact).sha256
__snyk_base_url = https://github.com/snyk/cli/releases/download/$(SNYK_VERSION)
.d.snyk: $(__snyk_bin) | .d.go .d.jq
ifdef SNYK_TOKEN
	@$(SNYK) auth $(SNYK_TOKEN)
endif
.PHONY: .d.snyk
$(__snyk_bin):
	rm -f $(GLOBAL_TMP_DIR)/snyk-*
	curl \
	  -s \
	  -L \
	  -o $(GLOBAL_TMP_DIR)/$(__snyk_artifact) \
	  $(__snyk_base_url)/$(__snyk_artifact)
	curl \
	  -s \
	  -L \
	  -o $(GLOBAL_TMP_DIR)/$(__snyk_artifact_sha) \
	  $(__snyk_base_url)/$(__snyk_artifact_sha)
	cd $(GLOBAL_TMP_DIR) && $(SHA256SUM) -c $(__snyk_artifact_sha)
	chmod a+x $(GLOBAL_TMP_DIR)/$(__snyk_artifact)
	mkdir -p $(@D)
	mv $(GLOBAL_TMP_DIR)/$(__snyk_artifact) $@

`

// Semgrep template
const semgrepTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/semgrep
SEMGREP_VERSION ?= 1.138.0
SEMGREP = $(realpath $(__semgrep_bin))
__semgrep_dir = $(TOOLS_DIR)/semgrep/$(SEMGREP_VERSION)
__semgrep_bin = $(__semgrep_dir)/bin/semgrep
__semgrep_venv = $(__semgrep_dir)/venv
.d.semgrep: $(__semgrep_bin) | .d.python
.PHONY: .d.semgrep
$(__semgrep_bin): | $(__semgrep_venv)
	mkdir -p $(__semgrep_dir)/bin
	$(__semgrep_venv)/bin/pip install semgrep==$(SEMGREP_VERSION)
	ln -sf $(realpath $(__semgrep_venv))/bin/semgrep $@
$(__semgrep_venv):
	$(PYTHON) -m venv $(__semgrep_venv)
	$(__semgrep_venv)/bin/pip install --upgrade pip

`

// goimports template
const goimportsTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/goimports
GOIMPORTS_VERSION ?= v0.29.0
GOIMPORTS_MOD ?= golang.org/x/tools
GOIMPORTS = env \
  GOENV="$(GO_ENV)" \
  GOROOT="$(GO_ROOT)" \
  PATH="$(GO_ROOT)/bin:$(PATH)" \
  $(realpath $(__goimports_bin))
__goimports_dir = $(TOOLS_DIR)/goimports/$(GOIMPORTS_VERSION)/$(GO_VERSION)
__goimports_bin = $(__goimports_dir)/goimports
.d.goimports: $(__goimports_bin)
.PHONY: .d.goimports
$(__goimports_bin): | .d.go
	mkdir -p $(__goimports_dir)
	cd '$(call go_dep_loc,$(GOIMPORTS_MOD)@$(GOIMPORTS_VERSION))' \
	&& $(GO) build -buildvcs=false -o $(abspath $(__goimports_bin)) ./cmd/goimports

`

// Node.js template
const nodejsTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/nodejs
NODEJS_VERSION ?= v20.18.0
NODE = $(realpath $(__nodejs_bin))
__nodejs_dir = $(TOOLS_DIR)/nodejs/$(NODEJS_VERSION)
__nodejs_bin = $(__nodejs_dir)/bin/node
__nodejs_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]")
__nodejs_arch = $(shell uname -m | sed -e 's/x86_64/x64/' -e 's/aarch64/arm64/')
__nodejs_file = node-$(NODEJS_VERSION)-$(__nodejs_kernel)-$(__nodejs_arch).tar.xz
__nodejs_base_url = https://nodejs.org/dist/$(NODEJS_VERSION)
__nodejs_bin_url = $(__nodejs_base_url)/$(__nodejs_file)
__nodejs_checksum_url = $(__nodejs_base_url)/SHASUMS256.txt
.d.nodejs: $(__nodejs_bin)
.PHONY: .d.nodejs
$(__nodejs_bin):
	rm -f $(GLOBAL_TMP_DIR)/$(__nodejs_file) $(GLOBAL_TMP_DIR)/SHASUMS256.txt
	curl -s -L -o $(GLOBAL_TMP_DIR)/$(__nodejs_file) $(__nodejs_bin_url)
	curl -s -L -o $(GLOBAL_TMP_DIR)/SHASUMS256.txt $(__nodejs_checksum_url)
	cd $(GLOBAL_TMP_DIR) && grep -F $(__nodejs_file) SHASUMS256.txt | $(SHA256SUM) -c -
	rm -rf $(__nodejs_dir)
	mkdir -p $(__nodejs_dir)
	tar -xf $(GLOBAL_TMP_DIR)/$(__nodejs_file) -C $(__nodejs_dir) --strip-components=1
	rm $(GLOBAL_TMP_DIR)/$(__nodejs_file) $(GLOBAL_TMP_DIR)/SHASUMS256.txt

`

// npm template
const npmTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/npm
NPM_VERSION ?= 10.8.2
NPM = $(realpath $(__npm_bin))
__npm_dir = $(TOOLS_DIR)/npm/$(NPM_VERSION)
__npm_bin = $(__npm_dir)/bin/npm
__npm_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]")
__npm_arch = $(shell uname -m | sed -e 's/x86_64/x64/' -e 's/aarch64/arm64/')
__npm_file = npm-$(NPM_VERSION).tgz
__npm_base_url = https://registry.npmjs.org/npm/-/npm-$(NPM_VERSION).tgz
.d.npm: $(__npm_bin) | .d.nodejs
.PHONY: .d.npm
$(__npm_bin): | .d.nodejs
	rm -f $(GLOBAL_TMP_DIR)/$(__npm_file)
	curl -s -L -o $(GLOBAL_TMP_DIR)/$(__npm_file) $(__npm_base_url)
	rm -rf $(__npm_dir)
	mkdir -p $(__npm_dir)
	cd $(__npm_dir) && $(NODE) $(GLOBAL_TMP_DIR)/$(__npm_file) --prefix=$(realpath $(__npm_dir))
	rm $(GLOBAL_TMP_DIR)/$(__npm_file)

`

// GitHub CLI template
const ghTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/gh
GH_VERSION ?= v2.65.0
GH_VERSION_int ?= $(patsubst v%,%,$(GH_VERSION))
GH = $(realpath $(__gh_bin))
gh_token = GITHUB_ACCESS_TOKEN="$$(if [ -z "$$GITHUB_ACCESS_TOKEN" ]; then \
		$(GH) auth token; \
	else \
		echo $$GITHUB_ACCESS_TOKEN; \
	fi)"
gh_token_plain = if [ -z "$$GITHUB_ACCESS_TOKEN" ]; then \
		$(GH) auth token; \
	else \
		echo $$GITHUB_ACCESS_TOKEN; \
	fi
gh_login = if ! $(GH) auth status 2> /dev/null && [ -z "$$GITHUB_ACCESS_TOKEN" ]; then \
		echo "Log in to GitHub CLI (command: make .a.run C='\$$(GH) auth login')"; \
		exit 1; \
	fi
__gh_dir = $(TOOLS_DIR)/gh/$(GH_VERSION)
__gh_bin = $(__gh_dir)/bin/gh
__gh_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]")
__gh_arch = $(shell uname -m | sed -e 's/x86_64/amd64/')
ifeq ($(__gh_kernel),darwin)
  __gh_file = gh_$(GH_VERSION_int)_macOS_$(__gh_arch).zip
else
  __gh_file = gh_$(GH_VERSION_int)_$(__gh_kernel)_$(__gh_arch).tar.gz
endif
__gh_base_url = https://github.com/cli/cli/releases/download/$(GH_VERSION)
__gh_bin_url = $(__gh_base_url)/$(__gh_file)
.d.gh: $(__gh_bin)
.PHONY: .d.gh
$(__gh_bin):
	echo "Installing GitHub CLI (gh)..."
	rm -f $(GLOBAL_TMP_DIR)/$(__gh_file)
	curl -v -L -o $(GLOBAL_TMP_DIR)/$(__gh_file) $(__gh_bin_url)
	rm -rf $(__gh_dir)
	echo "Extracting GitHub CLI (gh)..."
	mkdir -p $(__gh_dir)
	if [ "$(__gh_kernel)" = "darwin" ]; then \
		unzip -q $(GLOBAL_TMP_DIR)/$(__gh_file) -d $(__gh_dir); \
		mv $(__gh_dir)/gh_*/* $(__gh_dir); \
	else \
		tar -xzf $(GLOBAL_TMP_DIR)/$(__gh_file) -C $(__gh_dir); \
		mv $(__gh_dir)/gh_$(GH_VERSION_int)_$(__gh_kernel)_$(__gh_arch)/* $(__gh_dir); \
	fi
	chmod +x $(__gh_bin)
	echo "GitHub CLI (gh) installed successfully."
	rm $(GLOBAL_TMP_DIR)/$(__gh_file)
	touch $@

`

// Python template
const pythonTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/python
PYTHON_VERSION ?= system
PYTHON = $(shell which python3)
.d.python:
	@echo "Using system Python: $(PYTHON)"
	@$(PYTHON) --version
.PHONY: .d.python

`

// Python Black template
const blackTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/black
BLACK_VERSION ?= 24.10.0
BLACK = $(realpath $(__black_bin))
__black_dir = $(TOOLS_DIR)/black/$(BLACK_VERSION)
__black_bin = $(__black_dir)/black
__black_venv = $(__black_dir)/venv
.d.black: $(__black_bin) | .d.python
.PHONY: .d.black
$(__black_bin): | $(__black_venv)
	mkdir -p $(__black_dir)
	$(__black_venv)/bin/pip install black==$(BLACK_VERSION)
	ln -sf $(__black_venv)/bin/black $@
$(__black_venv):
	$(PYTHON) -m venv $(__black_venv)
	$(__black_venv)/bin/pip install --upgrade pip

`
