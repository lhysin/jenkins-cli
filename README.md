# Jenkins CLI

<p align="center">
  <img src="https://img.shields.io/github/stars/lhysin/jenkins-cli?style=flat-square" alt="GitHub Stars">
  <img src="https://img.shields.io/github/license/lhysin/jenkins-cli?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/v/release/lhysin/jenkins-cli?style=flat-square" alt="Release">
</p>

> A command-line interface for Jenkins REST API. Manage jobs, builds, nodes, and more from your terminal.

[English](README.md) | [한국어](README_ko.md)

## ✨ Features

- 🔗 **Multi-host management** - Login to multiple Jenkins servers and switch between them
- 📋 **Job management** - List, create, update, delete jobs
- 🚀 **Build management** - List builds, view logs, trigger builds
- 👁️ **View management** - List and view jobs in views
- 🖥️ **Node management** - List agents/nodes
- 🔓 **Anonymous access** - Works without authentication for public Jenkins instances

## 📦 Installation

### Quick Install (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/install.sh | sh
```

### Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/uninstall.sh | sh
```

### Download Binary

Download from [GitHub Releases](https://github.com/lhysin/jenkins-cli/releases)

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS | Apple Silicon | `jenkins-cli_darwin_arm64` |
| macOS | Intel | `jenkins-cli_darwin_amd64` |
| Linux | x86_64 | `jenkins-cli_linux_amd64` |
| Linux | ARM64 | `jenkins-cli_linux_arm64` |
| Windows | x86_64 | `jenkins-cli_windows_amd64.exe` |

```bash
# macOS/Linux (install.sh 사용)
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/install.sh | sh

# 또는 직접 다운로드
curl -fsSL https://github.com/lhysin/jenkins-cli/releases/latest/download/jenkins-cli_darwin_arm64 -o /usr/local/bin/jenkins-cli
chmod +x /usr/local/bin/jenkins-cli
```

### Build from Source

```bash
git clone https://github.com/lhysin/jenkins-cli.git
cd jenkins-cli
go build -o jenkins-cli ./cmd/
```

## 🚀 Quick Start

```bash
# Login to Jenkins server
jenkins-cli login local http://localhost:8080 admin
# Token: your-api-token

# Switch between hosts
jenkins-cli use prod

# List hosts
jenkins-cli hosts

# List jobs
jenkins-cli jobs list

# Trigger a build
jenkins-cli builds trigger my-job

# View build logs
jenkins-cli builds logs my-job 42
```

## 📚 Commands

### 🔗 Host Management

| Command | Description |
|---------|-------------|
| `login <name> <url> [user]` | Login to Jenkins (token prompted, anonymous OK) |
| `use <name>` | Switch to a host |
| `logout <name>` | Remove a host |
| `hosts` | List all hosts |

### 📋 Jobs

| Command | Description |
|---------|-------------|
| `jobs list` | List all jobs |
| `jobs info <name>` | Get job details |
| `jobs create <name> [-c config]` | Create a job |
| `jobs delete <name>` | Delete a job |
| `jobs update <name> [-c config]` | Update a job |

### 🚀 Builds

| Command | Description |
|---------|-------------|
| `builds list <job>` | List builds |
| `builds info <job> <n>` | Get build details |
| `builds logs <job> <n>` | Get build logs |
| `builds trigger <job>` | Trigger a build |

### 👁️ Views

| Command | Description |
|---------|-------------|
| `views list` | List all views |
| `views info <name>` | Get view details and jobs |

### 🖥️ Other

| Command | Description |
|---------|-------------|
| `nodes list` | List nodes |
| `status` | Check connection status |

## ⚙️ Configuration

Config stored at `~/.jenkins-cli/config.yaml`:

```yaml
current: local
hosts:
  local:
    url: http://localhost:8080
    user: admin
    token: your-api-token
  prod:
    url: https://jenkins.prod.com
    user: admin
    token: prod-token
```

## 🔑 Generate API Token

1. Login to Jenkins
2. Click your username → Configure
3. Under "API Token", click "Add new Token"
4. Copy the generated token

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run integration tests (requires Docker Jenkins)
./scripts/integration-test.sh
```

## 🤝 Contributing

### Commit Message Format

Version is automatically incremented based on commit messages:

| Type | Description | Version Change |
|------|-------------|----------------|
| `fix:` | Bug fixes | Patch +1 |
| `feat:` | New features | Minor +1 |
| `perf:` | Performance improvements | Patch +1 |
| `BREAKING CHANGE:` | Breaking changes | Major +1 |

Examples:
```bash
git commit -m "fix: resolve login timeout issue"
git commit -m "feat: add pipeline support"
git commit -m "feat!: drop Node.js 16 support"
```

### Workflow

1. Create feature branch: `git checkout -b feature/my-feature`
2. Commit changes with proper prefix
3. Push and create PR to main
4. Merge to main → automatic release

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
