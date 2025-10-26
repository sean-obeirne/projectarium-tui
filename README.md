# projectarium-tui

A Terminal User Interface (TUI) for the projectarium-v2 backend, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

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
export PROJECTARIUM_API_URL="http://localhost:8080/api/v1"
```

If not set, it defaults to `http://localhost:8080/api/v1`.

## Usage

Run the application:

```bash
./projectarium-tui
```

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
PROJECTARIUM_API_URL="http://localhost:8080/api/v1" ./projectarium-tui
```

## License

See [LICENSE](LICENSE) file.

