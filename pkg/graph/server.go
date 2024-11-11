package graph

import (
    "fmt"
    "net/http"
    "path/filepath"
)

const DefaultOutDir = "out"

func StartServer() {
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

        if len(files) == 1 {
            http.ServeFile(w, r, files[0])
            return
        }

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
