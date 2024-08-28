# SliverMoveMoveClient

SliverMoveMoveClient is a Go-based client tool designed to interact with the [Sliver C2](https://github.com/BishopFox/sliver) server. This tool implements several common lateral movement modules, allowing various operations on target hosts, such as steal passwords and credential searching.


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

- `timmy_mac.cfg` Sliver client connection configration file
- `main.go` - The main entry point of the program, containing the implementation of all modules.
- `go.mod` & `go.sum` - Go dependency management files.
- `modified_file.conf` - temp file for modify PAM files.
- `logger` - log ssh su sudo password and send in telegram (YOU NEED TO BUILD BY YOUR SELF)
Reference: https://github.com/qoo7972365/pam_logger

### 4. Run the Client

Run the following command to start the client and interact with the Sliver server:

```bash
./SliverMoveMoveClient --sliver-config /path/to/your/sliver.cfg --pam-logger /path/to/your/logger --command-logger /path/to/your/logger
```

### 5. Interaction

After starting the client, you can choose the session to operate on and select the module to execute. The tool supports selecting sessions and modules using the keyboard.

## Module List
- **Pam Logger**: Logs PAM authentication events to a specified location. 
- **SSH known host Search in all user**: Searches for SSH hosts in all users. (Not Available Under Working now)
- **Credential Search (Files)**: Searches for credentials in the file system. (Not Available Under Working now)
- **Credential Search (Memory)**: Searches for credentials in memory. (Not Available Under Working now)
