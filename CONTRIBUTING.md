## Development

### Prerequisites

- Go 1.16 or later
- Make (optional, for using Makefile commands)
- golangci-lint (optional, for linting)

### Common Tasks

```bash
# Build the project
make build

# Run tests
make test

# Generate test coverage report
make coverage

# Run linter
make lint

# Clean up
make clean
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./utils/...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request