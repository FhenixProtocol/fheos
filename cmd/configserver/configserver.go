package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

var configDir string // Global variable to hold the config directory

func main() {
	// Define and parse the command-line flag
	flag.StringVar(&configDir, "configDir", "./config", "Directory containing the config files")
	flag.Parse()

	// Set up the HTTP handler and start the server
	http.HandleFunc("/", serveConfigFile)
	log.Printf("Starting server on :8080 with config directory %s...", configDir)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveConfigFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[1:] // Extract file name from URL

	if fileName == "" {
		// for debugging
		listConfigFiles(w, r)
		return
	}

	// Clean the file name to prevent directory traversal
	fileName = filepath.Clean(fileName)
	filePath := filepath.Join(configDir, fileName)

	// Ensure the file path is absolute
	absConfigDir, err := filepath.Abs(configDir)
	if err != nil {
		http.Error(w, "Server error.", 500)
		return
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil || !filepath.HasPrefix(absFilePath, absConfigDir) {
		http.Error(w, "Invalid file path.", 400)
		return
	}

	data, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	w.Write(data)
}

func listConfigFiles(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		http.Error(w, "Could not list configuration files.", 500)
		return
	}

	fmt.Fprintf(w, "List of Configuration Files:\n")
	for _, file := range files {
		if !file.IsDir() {
			filePath := "/" + file.Name()
			fmt.Fprintf(w, "<a href='%s'>%s</a><br>", filePath, file.Name())
		}
	}
}
