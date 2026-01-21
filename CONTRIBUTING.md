# Contributing to mcp-netbird

Thank you for your interest in contributing to mcp-netbird! This document provides guidelines and information for contributors.

## Project Information

**Maintainer:** XNet Inc.  
**Lead Developer:** Joshua S. Doucette  
**License:** Apache License 2.0  
**Repository:** https://github.com/XNet-NGO/mcp-netbird

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and professional in all interactions.

## How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Check if the issue already exists
2. Provide a clear description of the problem
3. Include steps to reproduce (if applicable)
4. Specify your environment (OS, Go version, etc.)

### Suggesting Features

We welcome feature suggestions! Please:
1. Check if the feature has already been requested
2. Provide a clear use case
3. Explain how it would benefit users
4. Consider implementation complexity

### Submitting Pull Requests

1. **Fork the repository** and create a new branch
2. **Make your changes** following our coding standards
3. **Write tests** for new functionality
4. **Update documentation** as needed
5. **Run tests** to ensure everything passes
6. **Submit a pull request** with a clear description

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/XNet-NGO/mcp-netbird.git
cd mcp-netbird

# Install dependencies
go mod download

# Build the project
go build -o mcp-netbird ./cmd/mcp-netbird

# Run tests
go test ./...
```

### Project Structure

```
mcp-netbird/
├── cmd/mcp-netbird/     # Main application entry point
├── tools/               # NetBird API tool implementations
├── docs/                # Documentation
├── .kiro/               # Kiro AI specifications
├── mcpnetbird.go        # Core MCP server logic
├── go.mod               # Go module definition
└── README.md            # Project documentation
```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and types

### Code Organization

- Keep functions focused and single-purpose
- Group related functionality in packages
- Use interfaces for abstraction where appropriate
- Handle errors explicitly and provide context

### Testing

- Write unit tests for new functionality
- Aim for high test coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions
- Run the full test suite before submitting PRs

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## Documentation

### Code Documentation

- Add godoc comments for all exported functions, types, and constants
- Include examples in documentation where helpful
- Keep comments up-to-date with code changes

### README Updates

- Update README.md if you add new features
- Include usage examples for new functionality
- Update configuration documentation as needed

## Commit Messages

Write clear, descriptive commit messages:

```
Short summary (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain what changed and why, not how.

- Bullet points are okay
- Use present tense ("Add feature" not "Added feature")
- Reference issues and PRs where relevant (#123)
```

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat: add support for posture checks

Implement CRUD operations for NetBird posture checks including
OS version checks, geolocation checks, and process checks.

Closes #45

---

fix: handle empty API token gracefully

Previously, empty API tokens would cause a panic. Now we validate
the token and return a descriptive error message.

Fixes #67
```

## Pull Request Process

1. **Create a feature branch** from `main`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** and commit them
   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```

3. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Open a pull request** on GitHub
   - Provide a clear title and description
   - Reference related issues
   - Explain what changed and why
   - Include screenshots if relevant

5. **Address review feedback**
   - Respond to comments
   - Make requested changes
   - Push updates to your branch

6. **Wait for approval** and merge

## Review Process

All pull requests require:
- Code review by a maintainer
- All tests passing
- No merge conflicts
- Updated documentation (if applicable)

Maintainers will:
- Review code for quality and correctness
- Suggest improvements
- Approve or request changes
- Merge approved PRs

## Release Process

Releases are managed by XNet Inc. and Joshua S. Doucette using GoReleaser.

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Platforms

Official releases are built for:
- **Linux**: x86_64, ARM64, Debian packages
- **Windows**: x64
- **macOS**: x64 (Intel), ARM64 (Apple Silicon)

## License

By contributing to mcp-netbird, you agree that your contributions will be licensed under the Apache License 2.0.

All contributions must include the following copyright notice in new files:

```go
// Copyright 2025-2026 XNet Inc.
// Copyright 2025-2026 Joshua S. Doucette
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

## Attribution

This project was originally derived from the MCP Server for Grafana by Grafana Labs. When making substantial changes, please maintain appropriate attribution in the AUTHORS file.

## Getting Help

- **Issues**: Open an issue on GitHub
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact XNet Inc. for private inquiries

## Recognition

Contributors will be recognized in:
- The AUTHORS file
- Release notes
- Project documentation

Thank you for contributing to mcp-netbird!

---

**Maintained by XNet Inc.**  
**Lead Developer: Joshua S. Doucette**  
**Licensed under Apache License 2.0**
