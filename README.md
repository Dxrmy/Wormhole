# Wormhole

A tool that lets console players join custom Terraria servers by spoofing a local LAN game and routing the connection over the internet.

Works flawlessly on Windows, macOS, Linux, and mobile emulators without needing any prior networking knowledge.

## Usage

### Windows
1. Download `proxy-windows-amd64.exe`
2. Double-click to run it.
3. Your default web browser will automatically open to the Control Dashboard.
4. Enter your custom server IP and click **START BRIDGE**.

### macOS & Linux
Open your **Terminal** and run:
```bash
chmod +x proxy-linux-amd64 # or proxy-mac-arm64
./proxy-linux-amd64
```
Access the Dashboard via the local IP printed in your terminal.

## What this does

1. Spoofs a UDP LAN broadcast that shows up in the Terraria Console "Local Network" tab.
2. Intercepts the Xbox/PlayStation/Switch connection and tunnels it to your chosen remote PC server.
3. Hosts a beautiful local Web-UI dashboard for zero-friction configuration.
4. Operates entirely at Layer 4 (TCP), meaning zero latency impact and total immunity to crossplay version mismatches.
