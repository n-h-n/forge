## Overview

Forge is a Go-based build tool that manages development tools and generates Makefiles. It provides a unified way to manage various development tools like Go, Python, Node.js, npm, jq, Snyk, Semgrep, goimports, GitHub CLI, and Python Black across different projects.

Recreated as a copy of a tool that a former colleague and friend Tadej Borovšak had created.

## Features

- **Tool Management**: Downloads and manages development tools with version control
- **Makefile Generation**: Generates `Makefile.forge` files with tool configurations
- **Cross-platform Support**: Works on macOS and Linux (Windows not supported)
- **Caching**: Uses XDG cache directory for efficient tool storage

## Components Managed

1. **Go** (go1.25.1) - Go toolchain
2. **jq** (jq-1.7.1) - JSON processor
3. **Snyk** (v1.1295.0) - Security vulnerability scanner
4. **Semgrep** (1.75.0) - Static analysis security scanner
5. **goimports** (v0.29.0) - Go import formatter
6. **Python** (3.12.7) - Python toolchain
7. **Node.js** (v20.18.0) - Node.js runtime
8. **npm** (10.8.2) - npm package manager
9. **GitHub CLI** (v2.65.0) - GitHub command line tool
10. **Python Black** (24.10.0) - Python code formatter

## Usage

### Building

```bash
cd forge
go build -o forge .
```

### Commands

```bash
# Show help
./forge --help

# Show version
./forge version

# Sync project and generate Makefile
./forge sync <project-name> <source-makefile> <output-makefile>
```

### Example

```bash
./forge sync hermes Makefile Makefile.forge
```

## Architecture

The tool consists of several Go files:

- `main.go` - Main entry point with CLI setup
- `sync.go` - Core sync functionality and Makefile generation
- `templates.go` - Core and forge templates
- `components.go` - Individual component templates
