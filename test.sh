#!/bin/bash
# Integration test for projectarium-tui

set -e

echo "=== projectarium-tui Integration Test ==="
echo

# 1. Build the application
echo "1. Building application..."
go build -v -o projectarium-tui
echo "✓ Build successful"
echo

# 2. Verify the binary exists
if [ ! -f "./projectarium-tui" ]; then
    echo "✗ Binary not found"
    exit 1
fi
echo "✓ Binary exists"
echo

# 3. Build mock server
echo "2. Building mock server..."
cd examples
go build -v -o mock-server mock-server.go
cd ..
echo "✓ Mock server built"
echo

# 4. Start mock server in background
echo "3. Starting mock server..."
./examples/mock-server &
SERVER_PID=$!
echo "✓ Mock server started (PID: $SERVER_PID)"
echo

# Wait for server to be ready
sleep 2

# 5. Test API endpoints
echo "4. Testing API endpoints..."

# Test projects endpoint
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8888/api/projects)
if [ "$HTTP_CODE" != "200" ]; then
    echo "✗ Projects endpoint failed (HTTP $HTTP_CODE)"
    kill $SERVER_PID
    exit 1
fi
echo "✓ Projects endpoint working"

# Test single project endpoint
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8888/api/projects/1)
if [ "$HTTP_CODE" != "200" ]; then
    echo "✗ Single project endpoint failed (HTTP $HTTP_CODE)"
    kill $SERVER_PID
    exit 1
fi
echo "✓ Single project endpoint working"

# Test tasks endpoint
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8888/api/todos)
if [ "$HTTP_CODE" != "200" ]; then
    echo "✗ Todos endpoint failed (HTTP $HTTP_CODE)"
    kill $SERVER_PID
    exit 1
fi
echo "✓ Todos endpoint working"
echo

# 6. Verify data structure
echo "5. Verifying data structure..."
PROJECTS_COUNT=$(curl -s http://localhost:8888/api/projects | jq 'length')
if [ "$PROJECTS_COUNT" != "3" ]; then
    echo "✗ Expected 3 projects, got $PROJECTS_COUNT"
    kill $SERVER_PID
    exit 1
fi
echo "✓ Correct number of projects ($PROJECTS_COUNT)"

TODOS_COUNT=$(curl -s http://localhost:8888/api/todos | jq 'length')
if [ "$TODOS_COUNT" -lt "0" ]; then
    echo "✗ Failed to fetch todos"
    kill $SERVER_PID
    exit 1
fi
echo "✓ Todos endpoint accessible ($TODOS_COUNT todos)"
echo

# 7. Cleanup
echo "6. Cleaning up..."
kill $SERVER_PID
rm -f ./examples/mock-server
echo "✓ Cleanup complete"
echo

echo "=== All tests passed! ==="
echo
echo "To run the application manually:"
echo "  1. Start mock server: cd examples && go run mock-server.go"
echo "  2. In another terminal: ./projectarium-tui"
