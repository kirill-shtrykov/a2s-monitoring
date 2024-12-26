package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	a2s "github.com/rumblefrog/go-a2s"
)

type a2cInfo struct {
	ServerInfo *a2s.ServerInfo

	// Time when last player was seen
	LastPlayerSeen time.Time `json:"LastPlayerSeen,omitempty"`
}

type jsonExporter struct {
	a2sClient *a2s.Client
	info      a2cInfo
}

func (e *jsonExporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryInfo, err := e.a2sClient.QueryInfo()

	if err != nil {
		log.Printf("error getting server info: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 Internal Server Error")
		return
	}

	e.info.ServerInfo = queryInfo
	if e.info.ServerInfo.Players > 0 {
		e.info.LastPlayerSeen = time.Now()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e.info)
}

func NewJsonExporter(client *a2s.Client) *jsonExporter {
	return &jsonExporter{a2sClient: client}
}

// Retrieves the value of the environment variable named by the `key`
// It returns the value if variable present and value not empty
// Otherwise it returns string value `def`
func stringFromEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return strings.TrimSpace(v)
	}
	return def
}

func Run() int {
	addrHelpText := `
The address to listen.
Overrides the A2SMON_ADDR environment variable if set.
Default = :9112
	`
	statHelpText := `
The path to the server status.
Overrides the A2SMON_STATUS environment variable if set.
Default = /status
	`
	serverHelpText := `
The server address to monitoring.
Overrides the A2SMON_SERVER environment variable if set.
Default = :27015
	`

	addr := flag.String("address", stringFromEnv("A2SMON_ADDR", ":9112"), strings.TrimSpace(addrHelpText))
	stat := flag.String("status", stringFromEnv("A2SMON_STATUS", "/status"), strings.TrimSpace(statHelpText))
	server := flag.String("server", stringFromEnv("A2SMON_SERVER", ":27015"), strings.TrimSpace(serverHelpText))

	c, err := a2s.NewClient(*server)
	if err != nil {
		log.Printf("error creating A2S client: %v", err)
		return 1
	}

	e := NewJsonExporter(c)
	http.Handle(*stat, e)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>A2S Server Monitoring</title></head>
             <body>
             <h1>A2S Server Monitoring</h1>
             <p><a href='` + *stat + `'>Status</a></p>
             </body>
             </html>`))
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Printf("error starting HTTP server: %v", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(Run())
}
