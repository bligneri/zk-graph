package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ForceGraphTemplateName = "template/force_graph.tmpl"
	DefaultOutputFileName  = "out/output.html"
	DefaultOutDir          = "out"
)

type Link struct {
	SourcePath string `json:"sourcePath"`
	TargetPath string `json:"targetPath"`
}

type Note struct {
	Path     string `json:"path"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
}

type Node struct {
	ID    string `json:"id"`
	Group int    `json:"group"`
}

type GraphData struct {
	Nodes []Node              `json:"nodes"`
	Links []map[string]string `json:"links"`
}

func generateForceGraph(idTitleDict map[string][2]string, links []Link, highlight []string, outputFileName string) error {
	nodes := []Node{}
	highlightMap := make(map[string]bool)
	for _, h := range highlight {
		highlightMap[h] = true
	}

	// Populate nodes with highlighting logic
	for uid, val := range idTitleDict {
		title := val[0]
		group := 1
		if highlightMap[uid] {
			group = 2
		}
		nodes = append(nodes, Node{ID: title, Group: group})
	}

	// Create links with consistent source and target IDs
	linkList := []map[string]string{}
	pathToTitle := make(map[string]string)
	for path, val := range idTitleDict {
		pathToTitle[path] = val[0]
	}

	for _, link := range links {
		sourceID, sourceOK := pathToTitle[link.SourcePath]
		targetID, targetOK := pathToTitle[link.TargetPath]
		if sourceOK && targetOK {
			linkList = append(linkList, map[string]string{
				"source": sourceID,
				"target": targetID,
				"value":  "2",
			})
		} else {
			fmt.Printf("Warning: Missing source or target ID for link from %s to %s\n", link.SourcePath, link.TargetPath)
		}
	}

	data := GraphData{Nodes: nodes, Links: linkList}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling graph data: %v", err)
	}

	// Read template
	templateContent, err := os.ReadFile(ForceGraphTemplateName)
	if err != nil {
		return fmt.Errorf("error reading template file: %v", err)
	}

	tmpl, err := template.New("graph").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Create output file
	outFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Execute the template using GraphData as input
	err = tmpl.Execute(outFile, map[string]string{"Data": string(dataJSON)})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("Generated graph saved to %s\n", outputFileName)
	return nil
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := filepath.Glob(filepath.Join(DefaultOutDir, "*.html"))
		if err != nil {
			http.Error(w, "Error reading output directory", http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			http.Error(w, "No files found in output directory", http.StatusNotFound)
			return
		}

		// If only one file, serve it directly
		if len(files) == 1 {
			http.ServeFile(w, r, files[0])
			return
		}

		// Otherwise, list the files
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<html><body><h1>Select a file to view</h1><ul>")
		for _, file := range files {
			fileName := filepath.Base(file)
			fmt.Fprintf(w, `<li><a href="/view?file=%s">%s</a></li>`, fileName, fileName)
		}
		fmt.Fprintln(w, "</ul></body></html>")
	})

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("file")
		if fileName == "" {
			http.Error(w, "File not specified", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(DefaultOutDir, fileName)
		http.ServeFile(w, r, filePath)
	})

	fmt.Println("Serving files on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	jsonFilePath := flag.String("json_file", "", "Path to the input JSON file")
	highlight := flag.String("highlight", "", "Highlight zettel ID (comma-separated)")
	outputFileName := flag.String("output", DefaultOutputFileName, "Path to the output HTML file")
	serverMode := flag.Bool("server", false, "Start a web server to view output files")
	flag.Parse()

	if *serverMode {
		startServer()
		return
	}

	if *jsonFilePath == "" {
		fmt.Println("Error: json_file is required")
		os.Exit(1)
	}

	var data struct {
		Links []Link `json:"links"`
		Notes []Note `json:"notes"`
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

	err = generateForceGraph(idTitleDict, data.Links, highlightList, *outputFileName)
	if err != nil {
		fmt.Printf("Error generating force graph: %v\n", err)
	}
}
