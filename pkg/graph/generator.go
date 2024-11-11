package graph

import (
    "encoding/json"
    "fmt"
    "html/template"
    "os"
)

const (
    ForceGraphTemplateName = "template/force_graph.tmpl"
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

func GenerateForceGraph(idTitleDict map[string][2]string, links []Link, highlight []string, outputFileName string) error {
    nodes := []Node{}
    highlightMap := make(map[string]bool)
    for _, h := range highlight {
        highlightMap[h] = true
    }

    for uid, val := range idTitleDict {
        title := val[0]
        group := 1
        if highlightMap[uid] {
            group = 2
        }
        nodes = append(nodes, Node{ID: title, Group: group})
    }

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

    templateContent, err := os.ReadFile(ForceGraphTemplateName)
    if err != nil {
        return fmt.Errorf("error reading template file: %v", err)
    }

    tmpl, err := template.New("graph").Parse(string(templateContent))
    if err != nil {
        return fmt.Errorf("error parsing template: %v", err)
    }

    outFile, err := os.Create(outputFileName)
    if err != nil {
        return fmt.Errorf("error creating output file: %v", err)
    }
    defer outFile.Close()

    err = tmpl.Execute(outFile, map[string]string{"Data": string(dataJSON)})
    if err != nil {
        return fmt.Errorf("error executing template: %v", err)
    }

    fmt.Printf("Generated graph saved to %s\n", outputFileName)
    return nil
}
