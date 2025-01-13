package main

import (
	"encoding/json"
	"flag"
	log "log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	a2s "github.com/rumblefrog/go-a2s"
)

const (
	defaultListenAddr = ":9112"
	defaultA2SAddr    = "127.0.0.1:27015"
)

// setupLogging enables logging debug mode.
func setupLogging(debug bool) {
	if debug {
		log.SetLogLoggerLevel(log.LevelDebug)
		log.Debug("debug mode on")
	}
}

type a2cInfo struct {
	ServerInfo *a2s.ServerInfo `json:"ServerInfo"`

	// Time when last player was seen
	LastPlayerSeen time.Time `json:"LastPlayerSeen,omitempty"`
}

type JSONExporter struct {
	a2sClient *a2s.Client
	Info      a2cInfo `json:"info"`
}

func (e *JSONExporter) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	queryInfo, err := e.a2sClient.QueryInfo()
	if err != nil {
		log.Error("error getting server info:", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	e.Info.ServerInfo = queryInfo
	if e.Info.ServerInfo.Players > 0 {
		e.Info.LastPlayerSeen = time.Now()
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(e.Info); err != nil {
		log.Error("failed to encode JSON:", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func NewJSONExporter(client *a2s.Client) *JSONExporter {
	return &JSONExporter{a2sClient: client}
}

// Retrieves the value of the environment variable named by the `key`.
// It returns the value if variable present and value not empty.
// Otherwise it returns string value `def`.
func stringFromEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return strings.TrimSpace(v)
	}

	return def
}

// boolFromEnv retrieves the value of the environment variable named by the `key`.
// It returns the boolean value of the variable if present and valid.
// Otherwise, it returns the default value `def`.
func boolFromEnv(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		if err == nil {
			return parsed
		}
	}

	return def
}

// Flags represents a command line parameters.
type Flags struct {
	addr   string // The address to which HTTP server will bind.
	server string // The server address to monitoring.
	debug  bool   // Enables debug mode.
}

func parseFlags() *Flags {
	addrHelpText := `
The address to listen.
Overrides the A2SMON_ADDR environment variable if set.
Default = :9112
	`
	serverHelpText := `
The path to the server status.
Overrides the A2SMON_STATUS environment variable if set.
Default = /status
	`
	debugHelpText := `
Enables debug mode.
Overrides the A2SMON_DEBUG environment variable if set.
Default = false
	`

	flags := &Flags{
		addr:   stringFromEnv("A2SMON_ADDR", defaultListenAddr),
		server: stringFromEnv("A2SMON_SERVER", defaultA2SAddr),
		debug:  boolFromEnv("A2SMON_DEBUG", false),
	}

	flag.StringVar(&flags.addr, "address", flags.addr, strings.TrimSpace(addrHelpText))
	flag.StringVar(&flags.server, "path", flags.server, strings.TrimSpace(serverHelpText))
	flag.BoolVar(&flags.debug, "debug", flags.debug, strings.TrimSpace(debugHelpText))
	flag.Parse()

	return flags
}

func Run() int {
	log.Info("starting A2S Monitoring")

	flags := parseFlags()

	setupLogging(flags.debug)

	log.Debug("create A2S client", "server", flags.server)

	c, err := a2s.NewClient(flags.server)
	if err != nil {
		log.Error("error creating A2S client:", "error", err)

		return 1
	}

	e := NewJSONExporter(c)
	http.Handle("/", e)

	log.Debug("create HTTP server", "address", flags.addr)

	if err := http.ListenAndServe(flags.addr, nil); err != nil {
		log.Error("error starting HTTP server:", "error", err)

		return 1
	}

	return 0
}

func main() {
	os.Exit(Run())
}
