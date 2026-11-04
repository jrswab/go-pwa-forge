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
go install github.com/a-h/templ/cmd/templ@latest

# Start development server
make run
```

To run Tailwind use Tailwind's [Standalone CLI](https://tailwindcss.com/blog/standalone-cli). This forge does not use the npm/npx version when running locally.

Visit `http://localhost:3000` to see your application running.

## Project Structure

```
.
├── .github/
│   └── workflows/                # GitHub Actions workflows
│       ├── pr-workflow.yml
│       └── prod-workflow.yml
├── cmd/
│   └── main.go                   # Application entrypoint
├── config/
│   └── config.go                 # Application configuration
├── handlers/
│   └── handlers.go               # Handlers for serving pages
├── logs/                         # Application logs
│   └── site.log                  # The actual log file
├── media/                        # Media files
├── static/                       # Static assets
│   ├── input.css                 # TailwindCSS input
│   ├── manifest.json             # PWA Metadata for your site
│   ├── offline.html              # What to serve when the client is offline
│   ├── output.css                # Generated CSS
│   ├── htmx.min.js               # Currently v2.0.4
│   └── sw.js                     # Specific Javascript for the PWA Service Worker
├── views/             
│   ├── home/        
│   │   ├── home_content.templ    # The content for the home page
│   │   └── home_content_templ.go # Generated Go file from Templ
│   └── layout/        
│       ├── base.templ            # Main structure of you website (this won't change per page)
│       └── base_templ.go         # Generated Go file from Templ
```

## Development

### Prerequisites

- Go 1.20 or higher
- Make

### Available Commands

- `make run` - Starts a local development server with templ generation and TailwindCSS compilation
- `make pipeline` - Build assets for deployment (this is used by Github)

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
