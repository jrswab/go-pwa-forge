# go-pwa-forge

A Go template combining Chi and templ for building modern hypermedia-driven web applications with minimal JavaScript and maximum efficiency.

## Features

- **Chi** - High performance, minimalist web framework
- **Templ** - Type-safe HTML templates for Go
- **TailwindCSS** - Utility-first CSS framework
- **OpenBSD Ready** - Configured for OpenBSD deployment
- **GitHub Actions** - Automated PR and Production workflows

## Quick Start

```bash
# Clone the template
git clone https://github.com/jrswab/go-htmx-forge your-project-name

# Navigate to project
cd your-project-name

# Install dependencies
go mod tidy

# Start development server
make run
```

To run Tailwind use Tailwind's [Standalone CLI](https://tailwindcss.com/blog/standalone-cli). This forge does not use the npm/npx version when running locally.

Visit `http://localhost:3000` to see your application running.

## Project Structure

```
.
├── .github/workflows/            # GitHub Actions: pr-workflow.yml, prod-workflow.yml
├── .agents/ and .mem/            # Agent skills and project memory
├── config/
│   ├── config.go                 # Load, CORS origin parsing and validation
│   └── config_test.go
├── handlers/
│   ├── decoder.go                # POST body decoding
│   ├── home.go                   # GET / handler
│   ├── render.go                 # Templ render helper
│   └── render_test.go
├── internal/httpapp/             # Chi router, middleware, and routes (app.go)
├── media/                        # Untracked; init.sh creates it, served at /media/
├── scripts/
│   ├── appd                      # OpenBSD rc.d service script
│   └── init_obsd.sh              # Installs appd on OpenBSD
├── static/                       # Embedded and served at /static/
│   ├── htmx.min.js               # htmx v2.0.4
│   ├── input.css                 # Tailwind input
│   ├── manifest.json             # PWA metadata
│   ├── offline.html              # Served when the client is offline
│   ├── output.css                # Committed Tailwind output
│   └── sw.js                     # PWA service worker
├── views/
│   ├── home/                     # home_content.templ + generated _templ.go
│   └── layout/                   # base.templ + generated _templ.go
├── main.go                       # Process lifecycle: flags, logging, config, listener
├── main_test.go
├── init.sh                       # One-shot bootstrap for a new project
├── Makefile
├── tailwind.config.js
├── go.mod
└── go.sum
```

`main.go` creates `logs/` (with `site.log`) at runtime, so it is not in the tree. `init.sh` also creates `logs/` and `media/`.

## Development

### Prerequisites

- Go 1.23 or higher (the `go` directive in `go.mod`)
- Make

### Available Commands

- `make gen` - Regenerates `_templ.go` files from `.templ` sources. Runs templ at the version pinned in `go.mod` (`go run ... generate`); no local templ install is needed
- `make lint` - Installs and runs golangci-lint v2.1.6 over `./...`
- `make test` - Runs `go test -race -v` with a coverage profile (`coverage.out`)
- `make coverage` - Prints per-function coverage from `coverage.out` (`go tool cover -func`). Run `make test` first; it writes the profile
- `make run` - Runs `make gen`, then `go run main.go -l` (serves on `:3000`)
- `make prep` - Runs `gen`, `lint`, `test` (writes `coverage.out`), and `coverage`
- `make check` - Runs `lint`

## Deployment

The template includes two GitHub Actions workflows:

### PR Workflow
- Triggers on pushes to non-master branches
- Runs linting and tests
- Deploys to `/var/www/builds/<branch-name>/`

### Production Workflow
- Triggers on pushes to master
- Runs linting and tests
- Deploys to `/var/www/`
- Restarts the OpenBSD service

### Required Secrets

- `SSH_HOST` - Deployment server hostname
- `SSH_USER` - SSH username
- `SSH_PRIVATE_KEY` - SSH private key

### Environment Variables

- `CORS_ALLOWED_ORIGINS` - Optional. Unset or blank means no CORS. Comma-separated exact origins (`scheme://host[:port]`, `http` or `https`); a malformed list stops startup.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the BSD 3-Clause License.

## Support

If you find this template useful, consider:
- ⭐ Starring the repository
- 🐛 Opening issues for feature requests or bugs
- 🔀 Submitting pull requests
- ☕ Supporting the project on [Liberapay](https://liberapay.com/jrswab)

Created by [jrswab](https://github.com/jrswab) - Building tools to enhance developer workflows.
