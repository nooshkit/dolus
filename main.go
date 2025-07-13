package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type MethodConfig struct {
	Response int                    `yaml:"response"`
	Content  map[string]interface{} `yaml:"content"`
}

type ServerConfig struct {
	Port       int    `yaml:"port"`
	Listenaddr string `yaml:"listenaddr"`
}

type PathConfig map[string]MethodConfig // method -> MethodConfig

type Config struct {
	Paths  map[string]PathConfig `yaml:"paths"`
	Server ServerConfig          `yaml:"server"`
}

var routeTable = make(map[string]map[string]MethodConfig) // path -> method -> config
var serverConf string
var verbose bool

func loadConfig(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Normalize and store the config into routeTable
	for path, methods := range config.Paths {
		routeTable[path] = make(map[string]MethodConfig)
		for method, cfg := range methods {
			routeTable[path][strings.ToUpper(method)] = cfg
		}
	}

	// generate the serverConf string
	serverConf = fmt.Sprintf("%s:%d", config.Server.Listenaddr, config.Server.Port)
}

func logRequest(r *http.Request) {
	fmt.Println("---- Incoming Request ----")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Path:   %s\n", r.URL.Path)
	fmt.Println("Query Params:")
	for key, values := range r.URL.Query() {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Println("Headers:")
	for name, values := range r.Header {
		fmt.Printf("  %s: %v\n", name, values)
	}
	if r.Body != nil && r.ContentLength > 0 {
		body, _ := io.ReadAll(r.Body)
		fmt.Printf("Body: %s\n", string(body))
		r.Body = io.NopCloser(strings.NewReader(string(body)))
	}
	fmt.Println("--------------------------")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if verbose {
		logRequest(r)
	}

	path := r.URL.Path
	method := r.Method

	methods, exists := routeTable[path]
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cfg, methodExists := methods[method]
	if !methodExists {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(cfg.Response)
	json.NewEncoder(w).Encode(cfg.Content)
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	v := flag.Bool("v", false, "Turn on verbose request logging")

	flag.Parse()
	if *v {
		fmt.Println("---- Running In Verbose Mode ----")
		verbose = *v
	}

	loadConfig(*configPath)

	// Dynamically register a wildcard handler for all routes
	http.HandleFunc("/", handler)

	fmt.Printf("Dolus running at http://%s", serverConf)
	log.Fatal(http.ListenAndServe(serverConf, nil))
}
