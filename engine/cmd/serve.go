package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	servePort   int
	serveStatic bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the dev server",
	RunE:  runServe,
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func runServe(cmd *cobra.Command, args []string) error {
	contentDir := "../content"
	distDir := "../dist"
	if serveStatic {
		// Serve pre-built files
		distDir = "../dist"
	} else {
		// Build first
		buildCmd2 := &cobra.Command{}
		buildOutput = distDir
		if err := runBuild(buildCmd2, nil); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
	}

	// Find LAN address
	addr := fmt.Sprintf("0.0.0.0:%d", servePort)

	fmt.Printf("🌐 lyt dev server\n")
	fmt.Printf("   Local:   http://localhost:%d\n", servePort)
	fmt.Printf("   LAN:     http://%s\n", lanAddr())
	fmt.Printf("   Static:  %s\n\n", distDir)

	// File watcher for rebuild-on-change
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		debounce := time.NewTimer(500 * time.Millisecond)
		debounce.Stop()
		for {
			select {
			case <-watcher.Events:
				debounce.Reset(500 * time.Millisecond)
			case err := <-watcher.Errors:
				fmt.Printf("   ⚠️  watcher: %v\n", err)
			case <-debounce.C:
				buildCmd2 := &cobra.Command{}
				buildOutput = distDir
				if err := runBuild(buildCmd2, nil); err == nil {
					notifyClients("reload")
				}
			}
		}
	}()

	watchDir(watcher, contentDir)
	watchDir(watcher, "../templates")

	// HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, distDir)
	})

	return http.ListenAndServe(addr, mux)
}

func serveFile(w http.ResponseWriter, r *http.Request, dir string) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	filePath := filepath.Join(dir, filepath.FromSlash(strings.TrimPrefix(path, "/")))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Try index.html fallback (SPA-style)
		filePath = filepath.Join(dir, "index.html")
	}

	http.ServeFile(w, r, filePath)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	clientsMu.Lock()
	clients = append(clients, conn)
	clientsMu.Unlock()
	defer func() {
		clientsMu.Lock()
		for i, c := range clients {
			if c == conn {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		clientsMu.Unlock()
	}()
	for {
		conn.ReadMessage()
	}
}

var clients []*websocket.Conn
var clientsMu sync.Mutex

func notifyClients(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for _, c := range clients {
		c.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func watchDir(w *fsnotify.Watcher, dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			w.Add(path)
		}
		return nil
	})
}

func lanAddr() string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			if ip := toIP(a); ip != nil && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return "127.0.0.1"
}

func toIP(a net.Addr) net.IP {
	switch v := a.(type) {
	case *net.IPNet:
		return v.IP
	}
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 5173, "Port to serve on")
	serveCmd.Flags().BoolVar(&serveStatic, "static", false, "Serve pre-built dist without watching")
}
