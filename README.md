# report-tfvars

`report-tfvars` is a Go project designed to parse Terraform variable definitions from `.tf` files and output them in a structured format. The tool can handle complex default values and cleans up unnecessary whitespace and quotes. Useful to add a `make vars` target to existing Terraform directories using the `-dir` flag in this project can quickly show a help output of the potential variables offered by the project.

## Features

- Parses Terraform variable definitions
- Handles optional descriptions and default values
- Cleans up unnecessary whitespace and quotes
- Supports multiple operating systems and architectures

## Installation

To build the project, you need to have Go and Make installed on your system. If you do not have Go installed, you can use the [install-go](https://github.com/andreiwashere/install-go) package and ensure you have at least `1.22.0` on your system before building this project.

### Building with Make

This project uses a `Makefile` to facilitate building for different operating systems and architectures.

#### Building Binaries

To build binaries for all target OS/Arch combinations, run:

```bash
make all
```

This will generate binaries in the `bin` directory for the following combinations:

- darwin/amd64
- darwin/arm64
- linux/amd64
- linux/arm64
- windows/amd64

#### Cleaning Up

To remove all generated binaries, run:

```bash
make clean
```

#### Help

To display the help message, run:

```bash
make help
```

## Usage

After building the project, you can run the generated binary with the following options:

### Command Line Options

- `--file <file>`: Path to the Terraform file that contains variables.
- `--dir <directory>`: Path to the Terraform directory to scan for variables.

### Example

To parse variables from a specific file:

```bash
./bin/report-tfvars-linux-amd64 --file path/to/vars.tf
```

To scan a directory for `.tf` files and parse variables:

```bash
./bin/report-tfvars-linux-amd64 --dir path/to/terraform
```

## Development

### Prerequisites

- Go (https://golang.org/doc/install)
- Make (https://www.gnu.org/software/make/)

### Project Structure

- `main.go`: The main entry point for the application.
- `Makefile`: The Makefile for building the project.

### Running Locally

To run the project locally without building binaries, you can use the Go `run` command:

```bash
go run main.go --file path/to/vars.tf
```

or

```bash
go run main.go --dir path/to/terraform
```

## Contributing

Feel free to fork the project and then submit a pull request if you wish to improve this script.

## License

This script is licensed with the Apache 2.0 License.

