# Arbor

Arbor is a CLI tool for analyzing Git repositories. It provides various commands to analyze code metrics, commit history, and generate visualizations.

## Features

- **Lines of Code Analysis**: Analyze lines of code for each language across all commits in a Git repository
- **Chart Generation**: Generate interactive HTML charts showing code evolution over time
- **Multi-language Support**: Automatically detect and track multiple programming languages
- **Git Integration**: Seamlessly work with any Git repository

## Prerequisites

- Go 1.25.0 or later
- The `scc` tool installed on your system

### Installing scc

The `scc` (Sloc, Cloc and Code) tool is required for lines of code counting. You can install it using:

**macOS (using Homebrew):**
```bash
brew install scc
```

**Linux:**
```bash
# Download the latest release from https://github.com/boyter/scc/releases
# or use your package manager
```

**Windows:**
```bash
# Download the latest release from https://github.com/boyter/scc/releases
```

**Note**: The `scc` tool must be available in your system PATH for arbor to function properly. Without it, the analysis will fail.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/andrerfcsantos/arbor.git
cd arbor
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build
```

4. (Optional) Install to your system:
```bash
go install
```

## Usage

### Basic Usage

```bash
# Analyze the current directory (must be a Git repository)
arbor loc

# Analyze a specific repository
arbor loc /path/to/repository

# Get help
arbor --help
arbor loc --help
```

### Commands

#### `loc` - Lines of Code Analysis

Analyzes the lines of code for each language across all commits in a Git repository.

**Usage:**
```bash
arbor loc [repository-path]
```

**Arguments:**
- `repository-path` (optional): Path to the Git repository. If not provided, the current directory is assumed.

**What it does:**
1. Opens the specified Git repository
2. Records the current checkout state (branch/commit)
3. Iterates through all commits chronologically
4. **For each commit:**
   - Checks out the specific commit
   - Counts lines of code for each language using the `scc` tool
   - Shows progress and commit details
5. Generates a chart showing the evolution of lines of code over time
6. **Restores the repository to its original checkout state**

**Progress Tracking:**
The command provides real-time progress updates:
- Shows which commit is being processed
- Displays progress percentage
- Shows commit details (hash, message, author)
- Reports lines of code counts for each language
- Indicates when the original state is being restored

**Output:**
- Creates an HTML file named `loc_analysis.html` with an interactive chart
- The chart shows:
  - X-axis: Commits (chronologically ordered)
  - Y-axis: Lines of code
  - Series: Different programming languages
  - Interactive tooltips with commit details
- **Repository Safety**: The repository is automatically restored to its original state after analysis

## Project Structure

```
arbor/
├── cmd/                    # Command implementations
│   ├── root.go           # Root command and CLI setup
│   └── loc.go            # Lines of code analysis command
├── lib/                   # Shared library code
│   ├── git.go            # Git repository operations
│   ├── scc.go            # Lines of code counting
│   └── chart.go          # Chart generation
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
└── README.md             # This file
```

## Architecture

The project follows a clean architecture pattern:

- **`cmd/` package**: Contains all CLI command implementations using the Cobra library
- **`lib/` package**: Contains shared business logic that is agnostic of CLI commands
- **Dependencies**: Uses well-established Go libraries:
  - `github.com/spf13/cobra` for CLI command structure
  - `github.com/go-git/go-git/v5` for Git operations
  - `github.com/go-echarts/go-echarts/v2` for chart generation

## Development

### Adding New Commands

1. Create a new file in the `cmd/` package (e.g., `cmd/analyze.go`)
2. Define your command using Cobra
3. Add it to the root command in `cmd/root.go`
4. Implement the command logic using functions from the `lib/` package

### Adding New Library Functions

1. Create or modify files in the `lib/` package
2. Keep functions focused and reusable
3. Add appropriate error handling and documentation

### Building and Testing

```bash
# Build the application
go build

# Run tests
go test ./...

# Run with race detection
go test -race ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
```

## Dependencies

- **Cobra**: CLI framework for Go
- **go-git**: Pure Go implementation of Git
- **go-echarts**: Go port of Apache ECharts for chart generation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the LICENSE file for details.

## Future Enhancements

- ✅ **Commit-by-commit analysis**: Each commit is individually checked out and analyzed
- ✅ **Repository state restoration**: Automatically restores the original checkout state
- ✅ **Progress tracking**: Real-time progress indicators and commit details
- Support for analyzing specific file types or directories
- More detailed commit analysis (author statistics, file changes)
- Additional chart types and visualizations
- Integration with CI/CD pipelines
- Support for remote repositories
- Performance optimizations for large repositories
- Parallel processing for faster analysis
- Export results to different formats (CSV, JSON)
- Filter commits by date range or author
