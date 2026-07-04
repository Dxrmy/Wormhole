# Wormhole

A local proxy that tricks your Xbox, PlayStation, or Switch into thinking a remote PC server is running on your home network, allowing you to easily join custom Terraria servers.

It runs locally on Windows, macOS, Linux, or mobile network emulators.

## Installation

### Windows
1. Download `proxy-windows-amd64.exe` from this repository.
2. Double-click to run it.
3. Your web browser will automatically open a dashboard.
4. Enter the IP of the server you want to join and click **Start Bridge**.

### macOS & Linux
Download the binary for your specific system and run it via terminal:
```bash
chmod +x proxy-mac-arm64
./proxy-mac-arm64
```
It will print out a local IP address (like `http://127.0.0.1:8080`). Open that in your browser to access the dashboard.

## How to Connect on Console

Once the proxy is running and connected to a remote server:
1. Open Terraria on your Xbox, PlayStation, or Switch.
2. Select **Multiplayer** -> **Local** (or LAN).
3. The server will appear in the list (default name: "Epic Crossplay World").
4. Click Join!

## Command Line Options

For advanced users or headless servers, Wormhole can be completely controlled via terminal flags:

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-server` | *(empty)* | The remote server IP and port to bridge to (e.g. `198.51.100.24:7777`). |
| `-name` | `"Epic Crossplay World"` | The server name that broadcasts to your console. |
| `-port` | `7777` | The local TCP port the proxy listens on. |
| `-webport` | `8080` | The port used to host the Web UI dashboard. |
| `-headless` | `false` | Disables the Web UI entirely. (Requires `-server` to be set). |

**Example (Headless Mode):**
```bash
./proxy-linux-amd64 -headless -server 198.51.100.24:7777 -name "My Custom Server"
```

## What this does

1. Spoofs a UDP LAN broadcast so it shows up in your console's "Local Network" multiplayer tab.
2. Forwards the TCP connection from your console out to the remote server over the internet.
3. Provides a clean web interface to manage the connection without typing terminal commands.
