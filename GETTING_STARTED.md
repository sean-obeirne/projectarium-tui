# Getting Started with pj-tui

This guide will help you get started with the pj-tui Terminal User Interface.

## Quick Start

### 1. Build the Application

```bash
make build
# or
go build
```

### 2. Start the Mock Server (for testing)

In one terminal:

```bash
cd examples
go run mock-server.go
```

You should see:
```
Mock projectarium-v2 API server starting on :8888
Endpoints:
  GET /api/projects
  GET /api/projects/:id
  GET /api/todos

Press Ctrl+C to stop
```

### 3. Run the TUI

In another terminal:

```bash
./pj-tui
```

## Using the TUI

### Project List View

When you first launch the application, you'll see a list of projects:

```
Projects
┌─────────────────────────────────┐
│ Website Redesign                │
│ Redesign the company website... │
├─────────────────────────────────┤
│ Mobile App                      │
│ Develop mobile application...   │
├─────────────────────────────────┤
│ API Migration                   │
│ Migrate REST API to GraphQL     │
└─────────────────────────────────┘

↑/k up • ↓/j down • enter select • r refresh • q quit
```

**Navigation:**
- `↑` or `k` - Move up in the list
- `↓` or `j` - Move down in the list
- `Enter` - Select a project to view its kanban board
- `r` - Refresh the project list
- `q` - Quit the application

### Kanban Board View

After selecting a project, you'll see the kanban board:

```
📋 Website Redesign

┌─ To Do (2) ─┐  ┌─ In Progress (1) ─┐  ┌─ Done (1) ─┐
│             │  │                    │  │             │
│ Add contact │  │ Implement         │  │ Design      │
│ form        │  │ responsive layout │  │ homepage    │
│             │  │                    │  │ mockup      │
│ SEO         │  │                    │  │             │
│ optimization│  │                    │  │             │
│             │  │                    │  │             │
└─────────────┘  └────────────────────┘  └─────────────┘

←/h →/l move columns • ↑/k ↓/j move tasks • esc back • r refresh • q quit
```

**Navigation:**
- `←` or `h` - Move to previous column
- `→` or `l` - Move to next column
- `↑` or `k` - Move up in task list
- `↓` or `j` - Move down in task list
- `Esc` - Return to project list
- `r` - Refresh the board
- `q` - Quit the application

## Connecting to Real Backend

To connect to a real projectarium-v2 backend:

1. Set the API URL environment variable:
   ```bash
   export PROJECTARIUM_API_URL="http://your-backend:port/api"
   ```

2. Run the application:
   ```bash
   ./pj-tui
   ```

You can also create a `.env` file (see `.env.example`) and source it before running:

```bash
cp .env.example .env
# Edit .env with your settings
source .env
./pj-tui
```

## Troubleshooting

### "Connection refused" error

If you see:
```
Error: failed to get projects: Get "http://localhost:8888/api/projects": dial tcp: connect: connection refused

Press 'r' to retry, 'q' to quit
```

This means the backend API is not running. Either:
- Start the mock server (see Quick Start section)
- Start your projectarium-v2 backend
- Check that the PROJECTARIUM_API_URL is correct

### No projects shown

If the project list is empty:
- Press `r` to refresh
- Check that your backend has projects
- Verify the API endpoint is correct

## Development

### Running Tests

```bash
./test.sh
```

This will:
1. Build the application
2. Build and start the mock server
3. Test all API endpoints
4. Verify data structure
5. Clean up

### Code Formatting

```bash
go fmt ./...
```

### Code Verification

```bash
go vet ./...
```

## Support

For issues or questions, please refer to the main [README.md](README.md) or open an issue on GitHub.
