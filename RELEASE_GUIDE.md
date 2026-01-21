# Release Guide for mcp-netbird

This guide describes the release process for mcp-netbird maintained by XNet Inc. and Joshua S. Doucette.

## Release Platforms

Official releases are built for the following platforms:

### Linux
- **x86_64** (amd64) - Binary and Debian package
- **ARM64** (aarch64) - Binary and Debian package

### Windows
- **x64** (amd64) - ZIP archive

### macOS
- **x64** (Intel) - TAR.GZ archive
- **ARM64** (Apple Silicon) - TAR.GZ archive

### Docker
- **linux/amd64** - Multi-arch Docker image
- **linux/arm64** - Multi-arch Docker image

## Release Process

### Prerequisites

1. **Permissions**: Maintainer access to the repository
2. **Tools**: Git, Go 1.21+, GoReleaser
3. **Credentials**: GitHub token with release permissions

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Step-by-Step Release

#### 1. Prepare the Release

```bash
# Ensure you're on the main branch
git checkout main
git pull origin main

# Run tests
go test ./...

# Verify build works
go build -o mcp-netbird ./cmd/mcp-netbird
```

#### 2. Update Version Information

Update version references in:
- `README.md` (if version is mentioned)
- `CHANGELOG.md` (create if doesn't exist)
- Any version constants in code

#### 3. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag to GitHub
git push origin v1.0.0
```

#### 4. Automated Release

The GitHub Actions workflow will automatically:
1. Run tests
2. Build binaries for all platforms
3. Create Debian packages
4. Generate checksums
5. Create GitHub release with artifacts
6. Build and push Docker images

#### 5. Verify Release

Check the following:
- [ ] GitHub release created
- [ ] All platform binaries present
- [ ] Debian packages available
- [ ] Checksums file generated
- [ ] Docker images pushed to GHCR
- [ ] Release notes formatted correctly

### Manual Release (if needed)

If automated release fails, you can release manually:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Create release
goreleaser release --clean

# Or create snapshot (without publishing)
goreleaser release --snapshot --clean
```

## Release Artifacts

Each release includes:

### Binary Archives

- `mcp-netbird_VERSION_Linux_x86_64.tar.gz`
- `mcp-netbird_VERSION_Linux_arm64.tar.gz`
- `mcp-netbird_VERSION_Darwin_x86_64.tar.gz`
- `mcp-netbird_VERSION_Darwin_arm64.tar.gz`
- `mcp-netbird_VERSION_Windows_x86_64.zip`

### Debian Packages

- `mcp-netbird_VERSION_linux_x86_64.deb`
- `mcp-netbird_VERSION_linux_arm64.deb`

### Checksums

- `checksums.txt` - SHA256 checksums for all artifacts

### Docker Images

- `ghcr.io/xnet-ngo/mcp-netbird:VERSION`
- `ghcr.io/xnet-ngo/mcp-netbird:latest`

## Installation Instructions

### Linux (Debian/Ubuntu)

```bash
# Download the .deb package
wget https://github.com/XNet-NGO/mcp-netbird/releases/download/v1.0.0/mcp-netbird_1.0.0_linux_x86_64.deb

# Install
sudo dpkg -i mcp-netbird_1.0.0_linux_x86_64.deb

# Verify installation
mcp-netbird --version
```

### Linux (Other Distributions)

```bash
# Download and extract
wget https://github.com/XNet-NGO/mcp-netbird/releases/download/v1.0.0/mcp-netbird_1.0.0_Linux_x86_64.tar.gz
tar -xzf mcp-netbird_1.0.0_Linux_x86_64.tar.gz

# Move to PATH
sudo mv mcp-netbird /usr/local/bin/

# Verify installation
mcp-netbird --version
```

### macOS (Intel)

```bash
# Download and extract
wget https://github.com/XNet-NGO/mcp-netbird/releases/download/v1.0.0/mcp-netbird_1.0.0_Darwin_x86_64.tar.gz
tar -xzf mcp-netbird_1.0.0_Darwin_x86_64.tar.gz

# Move to PATH
sudo mv mcp-netbird /usr/local/bin/

# Verify installation
mcp-netbird --version
```

### macOS (Apple Silicon)

```bash
# Download and extract
wget https://github.com/XNet-NGO/mcp-netbird/releases/download/v1.0.0/mcp-netbird_1.0.0_Darwin_arm64.tar.gz
tar -xzf mcp-netbird_1.0.0_Darwin_arm64.tar.gz

# Move to PATH
sudo mv mcp-netbird /usr/local/bin/

# Verify installation
mcp-netbird --version
```

### Windows

1. Download `mcp-netbird_1.0.0_Windows_x86_64.zip`
2. Extract the ZIP file
3. Add the directory to your PATH
4. Open a new terminal and run `mcp-netbird --version`

### Docker

```bash
# Pull the image
docker pull ghcr.io/xnet-ngo/mcp-netbird:1.0.0

# Or use latest
docker pull ghcr.io/xnet-ngo/mcp-netbird:latest

# Run the container
docker run -d -p 8001:8001 ghcr.io/xnet-ngo/mcp-netbird:latest
```

## Release Checklist

Before releasing:

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version numbers updated
- [ ] No uncommitted changes
- [ ] On main branch
- [ ] All PRs merged

After releasing:

- [ ] Verify GitHub release created
- [ ] Test installation on each platform
- [ ] Verify Docker images work
- [ ] Update documentation if needed
- [ ] Announce release (if applicable)

## Troubleshooting

### Release Failed

1. Check GitHub Actions logs
2. Verify tag format (must be `vX.Y.Z`)
3. Ensure tests pass locally
4. Check GoReleaser configuration

### Missing Artifacts

1. Check `.goreleaser.yaml` configuration
2. Verify platform is included in builds
3. Check for build errors in logs

### Docker Build Failed

1. Verify Dockerfile.sse exists
2. Check Docker build logs
3. Ensure multi-arch build is supported

## Post-Release Tasks

1. **Update Documentation**: Ensure README reflects new version
2. **Announce Release**: Post to relevant channels
3. **Monitor Issues**: Watch for bug reports
4. **Plan Next Release**: Review roadmap and issues

## Hotfix Releases

For critical bug fixes:

1. Create hotfix branch from tag
2. Apply fix and test
3. Create new patch version tag
4. Follow normal release process

```bash
# Create hotfix branch
git checkout -b hotfix/v1.0.1 v1.0.0

# Make fixes
git commit -am "fix: critical bug"

# Tag and push
git tag -a v1.0.1 -m "Hotfix v1.0.1"
git push origin v1.0.1
```

## License and Attribution

All releases include:
- Apache License 2.0
- Copyright notices for XNet Inc. and Joshua S. Doucette
- Attribution to original Grafana Labs project
- NOTICE file with third-party attributions

## Support

For release-related questions:
- Open an issue on GitHub
- Contact XNet Inc.
- Email: joshua@xnet.company

---

**Maintained by XNet Inc.**  
**Lead Developer: Joshua S. Doucette**  
**Licensed under Apache License 2.0**
