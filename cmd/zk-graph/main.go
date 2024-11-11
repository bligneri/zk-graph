package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bligneri/zk-graph/pkg/graph"
)

func main() {
	jsonFilePath := flag.String("json_file", "", "Path to the input JSON file")
	highlight := flag.String("highlight", "", "Highlight zettel ID (comma-separated)")
	outputFileName := flag.String("output", "out/output.html", "Path to the output HTML file")
	serverMode := flag.Bool("server", false, "Start a web server to view output files")
	flag.Parse()

	if *serverMode {
		graph.StartServer()
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
		jsonData, err = os.ReadFile(*jsonFilePath)
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

	idTitleDict := make(map[string][2]string)
	for _, note := range data.Notes {
		title := note.Title
		if title == "" {
			title = note.Filename
		}
		idTitleDict[note.Path] = [2]string{title, note.Path}
	}

	highlightList := []string{}
	if *highlight != "" {
		highlightList = append(highlightList, *highlight) // Split by commas if necessary
	}

	err = graph.GenerateForceGraph(idTitleDict, data.Links, highlightList, *outputFileName)
	if err != nil {
		fmt.Printf("Error generating force graph: %v\n", err)
	}
}
