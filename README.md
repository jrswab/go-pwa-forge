# go-htmx-forge

A blazing-fast Go template combining Echo, HTMX, and templ for building modern hypermedia-driven web applications with minimal JavaScript and maximum efficiency.

## Features

- **Go Echo Framework** - High performance, minimalist web framework
- **Templ Integration** - Type-safe HTML templates for Go
- **TailwindCSS** - Utility-first CSS framework
- **OpenBSD Ready** - Configured for OpenBSD deployment
- **Redis Support** - Prepared for Redis integration
- **GitHub Actions** - Automated PR and Production workflows

## Quick Start

```bash
# Clone the template
git clone https://github.com/jrswab/go-htmx-forge your-project-name

# Navigate to project
cd your-project-name

# Install dependencies
go mod tidy
npm install tailwindcss
go install github.com/a-h/templ/cmd/templ@latest

# Start development server
make run
```

Visit `http://localhost:3000` to see your application running.

## Project Structure

```
.
├── .github/
│   └── workflows/        # GitHub Actions workflows
│       ├── pr-workflow.yml
│       └── prod-workflow.yml
├── cmd/
│   └── main.go          # Application entrypoint
├── logs/                # Application logs
├── static/              # Static assets
│   ├── input.css        # TailwindCSS input
│   └── output.css       # Generated CSS
└── media/               # Media files
```

## Development

### Prerequisites

- Go 1.20 or higher
- Node.js (for TailwindCSS)
- Make

### Available Commands

- `make run` - Start development server with templ generation and TailwindCSS compilation
- `make pipeline` - Build assets for deployment

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
- `REDIS_HOST` - Redis server hostname (production only)
- `REDIS_PASS` - Redis password (production only)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## Support

If you find this template useful, consider:
- ⭐ Starring the repository
- 🐛 Opening issues for feature requests or bugs
- 🔀 Submitting pull requests
- ☕ Supporting the project on [Liberapay](https://liberapay.com/jrswab)

Created by [jrswab](https://github.com/jrswab) - Building command-line tools and server templates that enhance developer workflows.
