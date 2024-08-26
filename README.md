# SliverMoveMoveClient

SliverMoveMoveClient is a Go-based client tool designed to interact with the [Sliver C2](https://github.com/BishopFox/sliver) server. This tool implements several common lateral movement modules, allowing various operations on target hosts, such as log management and credential searching.

## Features

- **Pam Logger Module**: Logs passwords from su, sudo, and ssh authentication and sends them to a specified location.
- **Credential Search**: Searches for sensitive credentials in the file system and memory.
- **File Management**: Supports uploading, downloading, modifying, and managing file permissions.

## Usage

### 1. Clone the Repository

```bash
git clone https://github.com/qoo7972365/SliverMoveMoveCleint.git
cd SliverMoveMoveClient
```

### 2. Build the Project

Ensure you have Go installed on your system. Run the following command to build the project:

```bash
go build -o SliverMoveMoveClient main.go
```

### 3. Configuration File

Before running the client, make sure you have a valid Sliver client configuration file in the project directory, such as `timmy_mac_35.236.161.97.cfg`.

### 4. Run the Client

Run the following command to start the client and interact with the Sliver server:

```bash
./SliverMoveMoveClient --config /path/to/your/config/file
```

### 5. Interaction

After starting the client, you can choose the session to operate on and select the module to execute. The tool supports selecting sessions and modules using the keyboard.

## Module List

- **Pam Logger**: Logs PAM authentication events to a specified location.
- **Credential Search (Files)**: Searches for credentials in the file system.
- **Credential Search (Memory)**: Searches for credentials in memory.

## Project Structure

- `main.go` - The main entry point of the program, containing the implementation of all modules.
- `go.mod` & `go.sum` - Go dependency management files.
- `modified_file.conf` - Example configuration file showing how to modify PAM files.

## Contribution Guidelines

If you find any issues or have suggestions for improvements, feel free to submit an Issue or Pull Request. For more detailed contribution guidelines, refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file (if available).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
