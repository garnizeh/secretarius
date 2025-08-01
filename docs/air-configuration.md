# Air Configuration Files

This project uses [Air](https://github.com/air-verse/air) for live reloading during development. Multiple configuration files are provided for different development scenarios.

## Configuration Files

### `.air.api.toml` - API Server Development
- **Purpose**: Live reloading for the API server
- **Binary**: `./tmp/api`
- **Environment**: Uses `.env.api-dev`
- **Build Command**: Standard build with development version
- **Includes**: All internal modules, API commands, and proto files
- **Usage**: `make watch-api` or `air -c .air.api.toml`

### `.air.worker.toml` - Worker Development
- **Purpose**: Live reloading for the worker server
- **Binary**: `./tmp/worker`
- **Environment**: Uses `.env.worker-dev`
- **Build Command**: Standard build with development version
- **Includes**: Worker-specific modules (ai, config, grpc, logging, worker)
- **Usage**: `make watch-worker` or `air -c .air.worker.toml`

### `.air.debug.toml` - Debug Mode
- **Purpose**: Development with debug flags enabled
- **Binary**: `./tmp/api-debug`
- **Environment**: Uses `.env.dev`
- **Build Command**: Includes race detection and debug symbols
  - `-race`: Race condition detection
  - `-gcflags='-N -l'`: Disable optimizations for debugging
- **Usage**: `make debug-api` or `air -c .air.debug.toml`

## Makefile Commands

```bash
# API Development
make watch-api      # Start API with live reload
make dev-api        # Alias for watch-api
make debug-api      # Start API with debug flags

# Worker Development
make watch-worker   # Start worker with live reload
make dev-worker     # Alias for watch-worker
```

## Features

### Optimized File Watching
- **Includes**: Only relevant directories and file types
- **Excludes**: Test files, vendor, logs, documentation, and build artifacts
- **Extensions**: `.go`, `.proto`, `.sql`, `.yaml`, `.toml`, `.json`

### Build Optimization
- **Faster Rebuilds**: Reduced delay (1000ms -> 500ms for debug)
- **Version Information**: Embedded development version info
- **Graceful Shutdown**: 2-3 second kill delay for proper cleanup
- **Log Files**: Separate build error logs in `tmp/` directory

### Environment Integration
- **Environment Files**: Automatically watches environment files
- **Args**: Environment files passed as arguments to binaries
- **Isolation**: Separate configs for API and Worker services

### Developer Experience
- **Colored Output**: Different colors for each service type
- **Timestamps**: Enabled for better debugging
- **Clean Exit**: Automatic cleanup on exit
- **Clear Screen**: Fresh view on each rebuild

## Usage Examples

### Standard Development
```bash
# Terminal 1: Start infrastructure
make dev-up

# Terminal 2: Start API with live reload
make watch-api

# Terminal 3: Start worker with live reload
make watch-worker
```

### Debug Session
```bash
# Terminal 1: Start infrastructure
make dev-up

# Terminal 2: Start API with debug flags
make debug-api
```

### Manual Air Usage
```bash
# API with custom config
air -c .air.api.toml

# Worker with custom config
air -c .air.worker.toml

# Debug mode
air -c .air.debug.toml
```

## Configuration Details

### Excluded Directories
- `vendor/` - Go modules
- `tmp/` - Temporary files
- `logs/` - Application logs
- `bin/` - Binary outputs
- `deployments/` - Docker configs
- `bruno-collection/` - API testing
- `docs/` - Documentation
- `testdata/`, `test_reports/`, `tests/` - Testing
- `.git/`, `.github/`, `.vscode/` - Version control and IDE

### Excluded Files
- `_test.go` - Test files
- `.test` - Test binaries
- `.sarif` - Security reports

### Build Flags

#### Standard Build
```bash
go build -ldflags='-X main.version=dev' -o ./tmp/api ./cmd/api
```

#### Debug Build
```bash
go build -race -gcflags='-N -l' -ldflags='-X main.version=dev-debug' -o ./tmp/api-debug ./cmd/api
```

## Troubleshooting

### Common Issues

1. **Permission Errors**: Ensure `tmp/` directory is writable
2. **Port Conflicts**: Check if services are already running
3. **Build Failures**: Check `tmp/build-errors.log` for details
4. **Environment Issues**: Verify `.env.*` files exist and are readable

### Performance Tips

1. **Exclude Large Directories**: Add rarely-changed directories to `exclude_dir`
2. **Limit File Types**: Only include necessary extensions in `include_ext`
3. **Adjust Delays**: Increase `delay` if builds are too frequent
4. **Use Specific Includes**: Specify exact directories in `include_dir`

## Version Requirements

- **Air**: v1.49.0 or later
- **Go**: 1.21 or later
- **Make**: For using Makefile commands
