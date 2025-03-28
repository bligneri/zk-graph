package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bligneri/zk-graph/pkg/graph"
)

func main() {
	jsonFilePath := flag.String("in", "", "Path to the input JSON file")
	outputFileName := flag.String("out", "/tmp/zk-graph/output.html", "Path to the output HTML file")
	serverMode := flag.Bool("server", false, "Start a web server to view output files")
	highlight := flag.String("highlight", "", "Highlight title or filename (comma-separated)")
	templateName := flag.String("template", "", "Name of the template")
	flag.Parse()

	if *serverMode {
		outputDir := filepath.Dir(*outputFileName)
		graph.StartServer(outputDir)
		return
	}

	if *jsonFilePath == "" {
		fmt.Println("Error: json_file is required")
		os.Exit(1)
	}

	var data struct {
		Links []graph.Link `json:"links"`
		Notes []graph.Note `json:"notes"`
	}

	var jsonData []byte
	var err error
	if *jsonFilePath == "-" {
		jsonData, err = io.ReadAll(os.Stdin)
	} else {
		jsonData, err = os.ReadFile(*jsonFilePath)
	}
	if err != nil {
		fmt.Printf("Error reading JSON input: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Printf("Error decoding JSON input: %v\n", err)
		os.Exit(1)
	}

	highlightList := []string{}
	if *highlight != "" {
		highlightList = strings.Split(*highlight, ",") // Split comma-separated highlights into a slice
	}

	err = graph.GenerateForceGraph(data.Notes, data.Links, highlightList, *outputFileName, *templateName)
	if err != nil {
		fmt.Printf("Error generating force graph: %v\n", err)
	}
}
