package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	targetServer = flag.String("server", "", "The remote server (IP:Port) to bridge to.")
	localPort    = flag.String("port", "7777", "The local TCP port to listen on for console connections.")
	serverName   = flag.String("name", "Epic Crossplay World", "The name broadcasted to the console.")
	webPort      = flag.String("webport", "8080", "The port for the Web UI.")
	headless     = flag.Bool("headless", false, "Disable the Web UI and run purely in the terminal.")
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Terraria Proxy Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #0B0E14;
            --card-bg: rgba(20, 24, 34, 0.65);
            --primary: #00FF88;
            --danger: #FF3366;
            --text-main: #FFFFFF;
            --text-muted: #8A92A3;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
        body {
            background: var(--bg-color); color: var(--text-main);
            min-height: 100vh; display: flex; align-items: center; justify-content: center;
            overflow: hidden; position: relative;
        }
        .blob {
            position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.5;
            animation: float 20s infinite alternate; z-index: 1;
        }
        .blob1 { width: 400px; height: 400px; background: rgba(0, 255, 136, 0.15); top: -100px; left: -100px; }
        .blob2 { width: 500px; height: 500px; background: rgba(51, 102, 255, 0.15); bottom: -150px; right: -150px; animation-delay: -10s; }
        @keyframes float { 0% { transform: translate(0, 0) rotate(0deg); } 100% { transform: translate(100px, 50px) rotate(180deg); } }
        
        .container {
            position: relative; z-index: 10; background: var(--card-bg);
            backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);
            border: 1px solid rgba(255, 255, 255, 0.05); border-radius: 24px;
            padding: 40px; width: 90%; max-width: 480px; box-shadow: 0 30px 60px rgba(0,0,0,0.5);
            transform: translateY(20px); opacity: 0; animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }
        @keyframes slideUp { to { transform: translateY(0); opacity: 1; } }

        h1 { font-weight: 800; font-size: 28px; margin-bottom: 8px; letter-spacing: -0.5px; background: linear-gradient(90deg, #fff, #a5b4cb); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        p.subtitle { color: var(--text-muted); font-size: 14px; margin-bottom: 32px; line-height: 1.5; }

        .form-group { margin-bottom: 24px; }
        label { display: block; font-size: 12px; text-transform: uppercase; letter-spacing: 1px; color: var(--text-muted); margin-bottom: 8px; font-weight: 600; }
        input {
            width: 100%; background: rgba(0, 0, 0, 0.2); border: 1px solid rgba(255, 255, 255, 0.1);
            color: white; padding: 16px; border-radius: 12px; font-size: 16px; outline: none; transition: all 0.3s ease;
        }
        input:focus { border-color: var(--primary); box-shadow: 0 0 15px rgba(0, 255, 136, 0.2); }
        input:disabled { opacity: 0.5; cursor: not-allowed; }

        button {
            width: 100%; padding: 18px; border: none; border-radius: 12px; font-size: 16px; font-weight: 700;
            cursor: pointer; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }
        .btn-start { background: var(--primary); color: #000; box-shadow: 0 10px 20px rgba(0, 255, 136, 0.2); }
        .btn-start:hover { transform: translateY(-2px); box-shadow: 0 15px 25px rgba(0, 255, 136, 0.3); }
        .btn-stop { background: rgba(255, 51, 102, 0.1); color: var(--danger); border: 1px solid var(--danger); }
        .btn-stop:hover { background: rgba(255, 51, 102, 0.2); }

        .status-badge {
            display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; border-radius: 20px;
            font-size: 12px; font-weight: 600; background: rgba(255,255,255,0.05); margin-bottom: 24px;
        }
        .dot { width: 8px; height: 8px; border-radius: 50%; }
        .dot.offline { background: var(--text-muted); }
        .dot.online { background: var(--primary); box-shadow: 0 0 10px var(--primary); animation: pulse 2s infinite; }
        @keyframes pulse { 0% { box-shadow: 0 0 0 0 rgba(0, 255, 136, 0.4); } 70% { box-shadow: 0 0 0 10px rgba(0, 255, 136, 0); } 100% { box-shadow: 0 0 0 0 rgba(0, 255, 136, 0); } }

        #toast {
            position: fixed; top: 20px; right: 20px; background: #fff; color: #000; padding: 16px 24px; border-radius: 8px;
            font-weight: 600; transform: translateX(150%); transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1); z-index: 100;
        }
        #toast.show { transform: translateX(0); }
    </style>
</head>
<body>
    <div class="blob blob1"></div><div class="blob blob2"></div>
    <div id="toast">Message</div>
    <div class="container">
        <div class="status-badge" id="statusBadge">
            <div class="dot offline" id="statusDot"></div><span id="statusText">System Offline</span>
        </div>
        <h1>Terraria LAN Proxy</h1>
        <p class="subtitle">Bypass console network restrictions and connect your Xbox/PS/Switch to custom PC servers.</p>
        <div class="form-group">
            <label>Remote Server (IP:Port)</label>
            <input type="text" id="target" placeholder="e.g. 198.51.100.24:7777">
        </div>
        <div class="form-group">
            <label>Console Display Name</label>
            <input type="text" id="name" value="Epic Crossplay Server">
        </div>
        <button id="mainBtn" class="btn-start" onclick="toggleProxy()">START BRIDGE</button>
    </div>
    <script>
        let isRunning = false;
        async function fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                updateUI(data.running, data.target, data.name);
            } catch (e) { console.error(e); }
        }
        function updateUI(running, target, name) {
            isRunning = running;
            document.getElementById('mainBtn').className = running ? 'btn-stop' : 'btn-start';
            document.getElementById('mainBtn').innerText = running ? 'STOP BRIDGE' : 'START BRIDGE';
            document.getElementById('statusDot').className = running ? 'dot online' : 'dot offline';
            document.getElementById('statusText').innerText = running ? 'Routing Traffic' : 'System Offline';
            document.getElementById('statusText').style.color = running ? 'var(--primary)' : '';
            if (target) document.getElementById('target').value = target;
            if (name) document.getElementById('name').value = name;
            document.getElementById('target').disabled = running;
            document.getElementById('name').disabled = running;
        }
        async function toggleProxy() {
            if (isRunning) {
                await fetch('/api/stop', { method: 'POST' });
                showToast("Proxy stopped successfully.");
            } else {
                const target = document.getElementById('target').value;
                const name = document.getElementById('name').value;
                if (!target) return showToast("Server IP cannot be empty!", true);
                const res = await fetch('/api/start', {
                    method: 'POST', headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ target, name })
                });
                if (!res.ok) return showToast("Error: " + await res.text(), true);
                showToast("Proxy bridge established!");
            }
            fetchStatus();
        }
        function showToast(msg, isError=false) {
            const toast = document.getElementById('toast');
            toast.innerText = msg;
            toast.style.background = isError ? "var(--danger)" : "#fff";
            toast.style.color = isError ? "#fff" : "#000";
            toast.classList.add('show');
            setTimeout(() => toast.classList.remove('show'), 3000);
        }
        fetchStatus();
        setInterval(fetchStatus, 2000);
    </script>
</body>
</html>`

type ProxyEngine struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	status bool
	target string
	name   string
}

var engine ProxyEngine

func main() {
	flag.Parse()

	if *headless {
		if *targetServer == "" {
			log.Fatal("Headless mode requires a -server IP to be provided.")
		}
		log.Println("Starting in Headless Mode...")
		engine.Start(*targetServer, *serverName)
		// Block forever
		select {}
	}

	// Auto-start if IP provided even in Web mode
	if *targetServer != "" {
		engine.Start(*targetServer, *serverName)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexHTML))
	})

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		engine.mu.Lock()
		defer engine.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"running": engine.status,
			"target":  engine.target,
			"name":    engine.name,
		})
	})

	http.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Target string `json:"target"`
			Name   string `json:"name"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		err := engine.Start(req.Target, req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/api/stop", func(w http.ResponseWriter, r *http.Request) {
		engine.Stop()
	})

	printWebUIAddresses(*webPort)
	openBrowser("http://127.0.0.1:" + *webPort)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+*webPort, nil))
}

func printWebUIAddresses(port string) {
	fmt.Println("\n=======================================================")
	fmt.Println("  TERRARIA PROXY - WEB UI IS RUNNING")
	fmt.Println("=======================================================")
	fmt.Println("Access the control dashboard from your web browser:")
	fmt.Printf(" ► Local Device: http://127.0.0.1:%s\n", port)

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
					fmt.Printf(" ► Network/Mobile Device: http://%s:%s\n", ip.String(), port)
				}
			}
		}
	}
	fmt.Println("=======================================================\n")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		// Do not attempt to open browser in headless server environments
		if os.Getenv("DISPLAY") != "" {
			err = exec.Command("xdg-open", url).Start()
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		// Silently ignore if browser can't be opened
	}
}

// Engine Methods
func (e *ProxyEngine) Start(target, name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.status {
		return fmt.Errorf("Proxy is already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	e.status = true
	e.target = target
	e.name = name

	go startUDPBroadcast(ctx, name, *localPort)
	go startTCPProxy(ctx, *localPort, target)
	return nil
}

func (e *ProxyEngine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.status && e.cancel != nil {
		e.cancel()
		e.status = false
	}
}

// Core Network Logic with Context Cancellation

func startUDPBroadcast(ctx context.Context, name, tcpPort string) {
	payload := []byte(fmt.Sprintf("Terraria145\x00%s\x00%s", tcpPort, name))
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		log.Printf("UDP Error: %v", err)
		return
	}
	
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			broadcastPacket(conn, payload)
		}
	}
}

func broadcastPacket(conn net.PacketConn, payload []byte) {
	if addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8888"); err == nil {
		conn.WriteTo(payload, addr)
	}
	ifaces, err := net.Interfaces()
	if err != nil { return }
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 { continue }
		addrs, err := iface.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok { continue }
			ip4 := ipnet.IP.To4()
			if ip4 == nil { continue }
			mask := ipnet.Mask
			if len(mask) != 4 { continue }
			bcastIP := net.IPv4(ip4[0]|^mask[0], ip4[1]|^mask[1], ip4[2]|^mask[2], ip4[3]|^mask[3])
			bcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:8888", bcastIP.String()))
			if err == nil {
				conn.WriteTo(payload, bcastAddr)
			}
		}
	}
}

func startTCPProxy(ctx context.Context, localPort, target string) {
	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		log.Printf("TCP Listen Error: %v", err)
		return
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return // Clean exit
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
		go handleConnection(clientConn, target)
	}
}

func handleConnection(clientConn net.Conn, target string) {
	defer clientConn.Close()
	if tcpConn, ok := clientConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	serverConn, err := net.DialTimeout("tcp", target, 10*time.Second)
	if err != nil { return }
	defer serverConn.Close()
	if tcpConn, ok := serverConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(serverConn, clientConn)
		if cw, ok := serverConn.(interface{ CloseWrite() error }); ok { cw.CloseWrite() }
	}()
	go func() {
		defer wg.Done()
		io.Copy(clientConn, serverConn)
		if cw, ok := clientConn.(interface{ CloseWrite() error }); ok { cw.CloseWrite() }
	}()
	wg.Wait()
}
