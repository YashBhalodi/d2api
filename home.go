package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"oss.terrastruct.com/d2/d2compiler"
	"oss.terrastruct.com/d2/d2exporter"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2renderers/textmeasure"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GenerateDiagram(s string) string {
	graph, err := d2compiler.Compile("", strings.NewReader(s), &d2compiler.CompileOptions{UTF16: true})
	if err != nil {
		fmt.Println(err)
		return "Error"
	}
	ruler, _ := textmeasure.NewRuler()
	graph.SetDimensions(nil, ruler)
	d2dagrelayout.Layout(context.Background(), graph)
	diagram, _ := d2exporter.Export(context.Background(), graph, d2themescatalog.NeutralDefault.ID)
	out, _ := d2svg.Render(diagram)
	ioutil.WriteFile(filepath.Join("out.svg"), out, 0600)
	return "out.svg"
}

type requestBody struct {
	DiagramCode string `json:"diagramCode"`
}

func GetDiagram(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	json.NewDecoder(r.Body).Decode(&body)
	outputFilePath := GenerateDiagram(body.DiagramCode)
	if outputFilePath == "Error" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorBody := make(map[string]string)
		errorBody["errorMessage"] = "Compiler error"
		errorBody["errorKey"] = "COMPILER_ERROR"
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	fileBytes, _ := ioutil.ReadFile(outputFilePath)
	enableCors(&w)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
	os.Remove(outputFilePath)
}

func main() {
	router := mux.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.HandleFunc("/getSvg", GetDiagram).Methods("POST")

	fmt.Println("Server at " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
