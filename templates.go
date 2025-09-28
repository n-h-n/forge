package main

// Core template contains the basic forge infrastructure
const coreTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/core
PROJECT_NAME ?=
ifeq ($(PROJECT_NAME),)
  $(error PROJECT_NAME variable not set)
endif
XDG_CACHE_DIR ?= $(HOME)/.cache
FORGE_ROOT ?= $(XDG_CACHE_DIR)/forge
__forge_helper := $(lastword $(MAKEFILE_LIST))
ifeq ($(words $(MAKEFILE_LIST)),1)
  INTERNAL_OP := 1
endif
BUILD_DIR ?= _build
OS = $(shell uname -s)
ARCH = $(shell uname -m)
ifneq ($(findstring NT,$(OS)),)
  $(error Unsupported platform: Forge cannot run natively on Windows. Please use macOS, Linux, WSL, or a virtual machine running a supported OS)
endif
SHA1SUM = $(shell which sha1sum)
SHA256SUM = $(shell which sha256sum)
SHA512SUM = $(shell which sha512sum)
ifeq ($(OS),Darwin)
  SHA1SUM = $(shell which shasum) -a 1
  SHA256SUM = $(shell which shasum) -a 256
  SHA512SUM = $(shell which shasum) -a 512
endif
ifeq ($(SHA1SUM),)
  $(error sha1sum not found)
endif
ifeq ($(SHA256SUM),)
  $(error sha256sum not found)
endif
ifeq ($(SHA512SUM),)
  $(error sha512sum not found)
endif
GIT_SHA = $(shell git rev-parse HEAD)
__forge_sym = $(BUILD_DIR)/forge
GLOBAL_CACHE_DIR = $(__forge_sym)/cache
GLOBAL_TMP_DIR = $(__forge_sym)/tmp/global
TOOLS_DIR = $(__forge_sym)/tools
__project_cache_dir = $(__forge_sym)/projects/$(PROJECT_NAME)
__project_tmp_dir = $(__forge_sym)/tmp/projects/$(PROJECT_NAME)
CACHE_DIR = $(BUILD_DIR)/cache
TMP_DIR = $(BUILD_DIR)/tmp
$(FORGE_ROOT) $(BUILD_DIR):
	mkdir -p $@
$(__forge_sym): | $(FORGE_ROOT) $(BUILD_DIR)
	ln -sfn $(FORGE_ROOT) $@
$(GLOBAL_CACHE_DIR) $(GLOBAL_TMP_DIR) $(TOOLS_DIR): | $(__forge_sym)
	mkdir -p $@
$(__project_cache_dir) $(__project_tmp_dir): | $(__forge_sym)
	mkdir -p $@
$(CACHE_DIR): | $(__project_cache_dir)
	ln -sfn $(realpath $(__project_cache_dir)) $@
$(TMP_DIR): | $(__project_tmp_dir)
	ln -sfn $(realpath $(__project_tmp_dir)) $@
$(__forge_helper): | \
  $(GLOBAL_CACHE_DIR) $(GLOBAL_TMP_DIR) $(TOOLS_DIR) $(CACHE_DIR) $(TMP_DIR)
.a.clean:
	rm -rf $(__project_cache_dir) $(__project_tmp_dir) $(BUILD_DIR)
.PHONY: .a.clean.local
.a.clean.tools:
	rm -rf $(TOOLS_DIR)
.PHONY: .a.clean.tools
.a.clean.cache:
	sudo rm -rf $(GLOBAL_CACHE_DIR)
.PHONY: .a.clean.cache
.a.nuke: .a.clean
	sudo rm -rf $(FORGE_ROOT)
.PHONY: .a.nuke
.a.run:
	$(C)
.PHONY: .a.run
__text_color_red = \033[0;31m
__text_color_green = \033[0;32m
__text_normal = \033[0m
__text_bold = \033[1m

`

// Forge template contains the forge binary management
const forgeTemplate = `# SOURCE: https://github.com/n-h-n/forge/blob/main/internal/component/data/forge
FORGE_VERSION ?= %s
ifeq ($(FORGE_VERSION),)
  $(error FORGE_VERSION variable not set)
endif
__forge_dir = $(TOOLS_DIR)/forge/$(FORGE_VERSION)
__forge_bin = $(__forge_dir)/forge
__forge_kernel = $(shell uname -s | tr "[:upper:]" "[:lower:]")
__forge_arch = $(shell uname -m | sed -e 's/x86_64/amd64/')
__forge_asset = forge-$(__forge_kernel)-$(__forge_arch)
__forge_repo = n-h-n/forge
__forge_marker = $(__forge_dir)/marker.forge
FORGE ?= $(realpath $(__forge_bin))
.a.forge.sync: SYNC_SOURCE = Makefile
.a.forge.sync: $(__forge_marker)
	$(FORGE) sync $(PROJECT_NAME) $(SYNC_SOURCE) $(__forge_helper)
.PHONY: .a.forge.sync
.a.forge.update:
	git ls-remote --tags "https://github.com/$(__forge_repo)" \
	| cut -d/ -f3- \
	| sed -e "s/^v//" \
	| grep -v "[^0-9\.]" \
	| sort -t. -k 1,1n -k 2,2n -k 3,3n \
	| tail -n1 \
	| xargs -ITAG \
	    sed -i.bak -r -e 's/^FORGE_VERSION .+/FORGE_VERSION ?= vTAG/' \
	    $(__forge_helper)
	rm $(__forge_helper).bak
	$(MAKE) .a.forge.sync
.PHONY: .a.forge.update
$(__forge_marker): .d.gh
	rm -rf $(__forge_dir)
	mkdir -p $(__forge_dir)
	$(GH) release download \
	  --repo $(__forge_repo) \
	  --output $(__forge_bin)\
	  --pattern $(__forge_asset) \
	  $(FORGE_VERSION)
	chmod a+x $(__forge_bin)
	touch $@

`
