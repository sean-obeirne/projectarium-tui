# pj-tui

A Terminal User Interface (TUI) for the projectarium-v2 backend, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Preview

### Project List View
```
┌─────────────────────────────────────────────────────────┐
│                       Projects                          │
├─────────────────────────────────────────────────────────┤
│  Website Redesign                                       │
│  Redesign the company website with modern UI            │
│                                                          │
│  Mobile App                                             │
│  Develop mobile application for iOS and Android         │
│                                                          │
│  API Migration                                          │
│  Migrate REST API to GraphQL                            │
└─────────────────────────────────────────────────────────┘
  ↑/k up • ↓/j down • enter select • r refresh • q quit
```

### Kanban Board View
```
📋 Website Redesign

┌─ To Do (2) ──────┐  ┌─ In Progress (1) ┐  ┌─ Done (1) ──────┐
│                  │  │                   │  │                 │
│ ┌──────────────┐ │  │ ┌───────────────┐ │  │ ┌─────────────┐ │
│ │ Add contact  │ │  │ │ Implement     │ │  │ │ Design      │ │
│ │ form         │ │  │ │ responsive    │ │  │ │ homepage    │ │
│ └──────────────┘ │  │ │ layout        │ │  │ │ mockup      │ │
│                  │  │ └───────────────┘ │  │ └─────────────┘ │
│ ┌──────────────┐ │  │                   │  │                 │
│ │ SEO          │ │  │                   │  │                 │
│ │ optimization │ │  │                   │  │                 │
│ └──────────────┘ │  │                   │  │                 │
└──────────────────┘  └───────────────────┘  └─────────────────┘

  ←/h →/l move columns • ↑/k ↓/j move tasks • esc back • r refresh • q quit
```

## Features

- 📋 View and navigate projects
- 🎯 Kanban board view with columns (To Do, In Progress, Done)
- ⌨️ Keyboard-driven navigation
- 🎨 Beautiful terminal UI with lipgloss styling

## Installation

### Build from source

```bash
git clone https://github.com/sean-obeirne/projectarium-tui.git
cd projectarium-tui
go build
```

## Configuration

The application connects to the projectarium-v2 API backend. Configure the API endpoint using the `PROJECTARIUM_API_URL` environment variable:

```bash
export PROJECTARIUM_API_URL="http://localhost:8888/api"
```

If not set, it defaults to `http://localhost:8888/api`.

## Usage

Run the application:

```bash
./pj-tui
```

### Demo

To try the application with sample data, you can use the included mock server:

1. Start the mock server in one terminal:
   ```bash
   cd examples
   go run mock-server.go
   ```

2. In another terminal, run the TUI:
   ```bash
   ./pj-tui
   ```

The mock server provides sample projects and tasks to demonstrate the full functionality.

### Keyboard Controls

#### Project List View
- `↑/k` or `↓/j` - Navigate up/down
- `Enter` - Select project and view kanban board
- `r` - Refresh project list
- `q` - Quit

#### Kanban Board View
- `←/h` or `→/l` - Move between columns
- `↑/k` or `↓/j` - Navigate tasks within a column
- `Esc` - Return to project list
- `r` - Refresh board
- `q` - Quit

## Requirements

- Go 1.21 or later
- projectarium-v2 backend running and accessible

## Development

### Mock Server

A mock server is included in `examples/mock-server.go` for testing and development. It provides sample data with three projects and multiple tasks in different states.

To run the mock server:

```bash
cd examples
go run mock-server.go
```

The server will start on `http://localhost:8888` and provide the following endpoints:
- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get a specific project
- `GET /api/todos` - Get all todos

### Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # API client for projectarium-v2
│   │   ├── client.go      # HTTP client implementation
│   │   └── types.go       # Data structures (Project, Task, etc.)
│   ├── config/            # Configuration management
│   │   └── config.go      # Load config from environment
│   └── tui/               # Terminal UI components
│       ├── model.go       # Main Bubble Tea model
│       └── kanban.go      # Kanban board view
```

### Building

```bash
go build
```

### Running

```bash
PROJECTARIUM_API_URL="http://localhost:8888/api" ./pj-tui
```

## License

See [LICENSE](LICENSE) file.

