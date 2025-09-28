# Forge Installation Guide

## For Users

### Installation via GitHub CLI

```bash
# Download forge binary
gh release download --repo n-h-n/forge --pattern "forge-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/')" --output forge
chmod +x forge
sudo mv forge /usr/local/bin/

# Verify installation
forge --help
```

### Installation from Source

Build forge from source using Go:

```bash
# Clone the repository
git clone https://github.com/n-h-n/forge.git
cd forge

# Build the binary
go build -o forge .

# Install to system PATH (optional)
sudo mv forge /usr/local/bin/

# Verify installation
forge --help
```

### Generating Forge's Makefile System

1. Add to your project's Makefile:
```makefile
include Makefile.forge
```

2. Generate Makefile.forge:
```bash
forge sync your-project Makefile Makefile.forge
```

3. Install forge tools:
```bash
make .a.forge.sync
```

### Using Semgrep (Alternative to Snyk)

Once forge is installed, you can use Semgrep for static analysis:

```bash
# Install Semgrep
make .d.semgrep

# Scan code for security issues
$(SEMGREP) --config=auto .

# Scan with specific ruleset
$(SEMGREP) --config=p/security .

# Scan specific file types
$(SEMGREP) --config=auto --include="*.py,*.js" .

# Generate SARIF output
$(SEMGREP) --config=auto --output=semgrep-results.sarif .
```

### Using Node.js and npm

Once forge is installed, you can use Node.js and npm:

```bash
# Install Node.js
make .d.nodejs

# Install npm (depends on Node.js)
make .d.npm

# Use Node.js
$(NODE) your_script.js

# Use npm
$(NPM) install
$(NPM) run build
$(NPM) test
```

### Using Python and Python Black

Once forge is installed, you can use Python and Python Black:

```bash
# Install Python
make .d.python

# Install Black (depends on Python)
make .d.black

# Use Python
$(PYTHON) your_script.py

# Use Black for formatting
$(BLACK) your_script.py

# Format all Python files in current directory
$(BLACK) .

# Check formatting without making changes
$(BLACK) --check .
```

## Asset Naming Convention

The forge system expects these asset names:
- `forge-linux-amd64` (Linux x86_64)
- `forge-darwin-amd64` (macOS Intel)
- `forge-darwin-arm64` (macOS Apple Silicon)

## Version Management

- Update `version` variable in `main.go` for new releases
- Use semantic versioning (e.g., v1.0.0, v1.1.0, v2.0.0)
- GitHub Actions will automatically build and release when you push tags
