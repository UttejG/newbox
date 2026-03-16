# newbox 📦

> Set up a brand-new machine in minutes with an interactive TUI

newbox is a cross-platform CLI that helps you install your favorite tools on a fresh machine using native package managers — `brew` on macOS, `winget` on Windows, and `apt`/`dnf`/`pacman` on Linux.

## Install

**macOS / Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/UttejG/newbox/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/UttejG/newbox/main/scripts/install.ps1 | iex
```

**Go install:**
```bash
go install github.com/uttejg/newbox/cmd/newbox@latest
```

## Usage

```bash
newbox              # Launch TUI — pick a profile and install
newbox --dry-run    # Preview what would be installed (no changes)
newbox --dry-run --summary  # Print text summary of install plan
newbox --dry-run --json     # Machine-readable JSON output
newbox list         # List all available tools for your OS
newbox list --category messaging  # Filter by category
newbox --version    # Print version
```

## TUI Flow

```
┌─────────────────────────────────────┐
│   📦 newbox                         │
│   Set up your machine               │
│                                     │
│   Detected: macOS (ARM64)           │
│   Package manager: brew             │
│   Press Enter to continue           │
└─────────────────────────────────────┘
           │
           ▼
  Select a Profile:
  ▸ 🏗️  Developer   Full dev environment
    🎨  Creative    Media and design tools
    ⚡  Minimal     Just the essentials
    🚀  Full        Everything available
    🔧  Custom      Pick your own

           │
           ▼
  Select Categories:     [Tab: next]
  [x] 💬 Messaging
  [x] 🌐 Web Browsers
  [x] 🔧 CLI Essentials
  [ ] 🎬 Media Players
  ...

           │
           ▼
  Installing... (3/13)
  ✅ Signal
  ✅ Telegram
  ⏳ WhatsApp (installing)
  ⬚  Discord
  ░░░░████████ 23%
```

## Features

- **Interactive TUI** — navigate with arrow keys, toggle with space
- **5 profiles** — Developer, Creative, Minimal, Full, Custom
- **23 categories** — from messaging to cloud CLIs to AI tools
- **100+ tools** — with per-OS package mappings
- **Dry-run mode** — preview everything before installing
- **Resume support** — picks up where it left off if interrupted
- **dotfiles integration** — tools from your existing dotfiles pre-selected (★)
- **Cross-platform** — macOS (brew/mas), Windows (winget), Linux (apt/dnf/pacman/flatpak)

## Categories

| Category | Tools |
|----------|-------|
| 📧 Email Clients | Spark, Thunderbird, Mailspring |
| 💬 Messaging | Signal, Telegram, WhatsApp, Discord, Slack |
| 🌐 Web Browsers | Chrome, Firefox, Brave, Arc, Edge |
| 🖥️ Terminal Emulators | iTerm2, Ghostty, Alacritty, Warp |
| 📝 Text Editors | VS Code, Cursor, Zed, Neovim, Vim |
| 🏗️ IDEs | JetBrains Toolbox, Xcode, Android Studio |
| 🔧 CLI Essentials | git, ripgrep, fzf, bat, eza, jq, curl |
| 🐚 Shell Setup | Fish, Starship, zsh plugins |
| 🐳 Containers & VMs | Docker Desktop, Podman, UTM |
| ☁️ Cloud CLIs | AWS CLI, Azure CLI, kubectl, Terraform |
| 🗄️ Databases | PostgreSQL, MySQL, Redis, SQLite |
| 🗄️ DB Clients | TablePlus, DBeaver |
| 🐍 Languages | Node.js, Python, Go, Rust, Java, .NET |
| 📦 Version Managers | mise, pyenv, rbenv |
| 🎬 Media Players | VLC, Spotify, IINA |
| 🎨 Creative Tools | Figma, GIMP, OBS Studio |
| 📋 Productivity | Obsidian, Notion, Raycast, Rectangle |
| 🔐 Security | 1Password, Bitwarden, GnuPG |
| ☁️ Cloud Storage | Dropbox, Google Drive |
| 🤖 AI Tools | Ollama, ChatGPT, Claude |
| 💰 Finance | GnuCash |
| 📡 Networking | Wireshark, Postman, ngrok |
| 🖼️ Screenshot & Record | OBS, Flameshot, Monosnap |

## Dry-Run Mode

Preview your install plan without making any changes:

```bash
newbox --dry-run
```

Or get a machine-readable JSON plan for CI:

```bash
newbox --dry-run --json
```

```json
{
  "platform": {"os": "macos", "arch": "arm64", "package_manager": "brew"},
  "steps": [
    {"tool": "Signal", "command": "brew install --cask signal", "status": "would_install"},
    {"tool": "git",    "command": "brew install git",           "status": "already_installed"}
  ],
  "summary": {"would_install": 12, "already_installed": 5}
}
```

## Architecture

newbox uses **Hexagonal Architecture** (Ports & Adapters) — the idiomatic Go equivalent of C#'s Clean Architecture:

```
┌─────────────────────────────────────────┐
│              ADAPTERS (outer)            │
│  Driving:  TUI, CLI flags               │
│  Driven:   brew, winget, apt, FileStore │
├─────────────────────────────────────────┤
│              PORTS (interfaces)          │
│  PackageManager, CommandRunner           │
│  CatalogService, StateStore             │
├─────────────────────────────────────────┤
│              CORE (inner)                │
│  Domain:    Tool, Category, Platform    │
│  Services:  InstallService, CatalogSvc  │
└─────────────────────────────────────────┘
```

**Dependency rule**: Core never imports adapter packages. All I/O goes through port interfaces.

## Development

```bash
git clone https://github.com/UttejG/newbox
cd newbox
make test          # run all unit tests
make test-race     # with race detector
make build         # build binary
make coverage      # coverage report
go run ./cmd/newbox --dry-run  # test locally
```

## Contributing

1. Fork the repo
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `make test` passes
5. Open a PR

## License

MIT — see [LICENSE](LICENSE)
