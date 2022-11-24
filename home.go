package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

func IncludeHTML(path string) template.HTML {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("includeHTML - error reading file: %v", err)
		return ""
	}

	return template.HTML(string(b))
}

func GenerateDiagram(s string) string {
	println(s)
	graph, _ := d2compiler.Compile("", strings.NewReader(s), &d2compiler.CompileOptions{UTF16: true})
	ruler, _ := textmeasure.NewRuler()
	graph.SetDimensions(nil, ruler)
	d2dagrelayout.Layout(context.Background(), graph)
	diagram, _ := d2exporter.Export(context.Background(), graph, d2themescatalog.NeutralDefault.ID)
	out, _ := d2svg.Render(diagram)
	ioutil.WriteFile(filepath.Join("out.svg"), out, 0600)
	return "out.svg"
}

func GetDiagram(w http.ResponseWriter, r *http.Request) {
	GenerateDiagram("a->b")

	tmpl := template.New("sample")
	tmpl.Funcs(template.FuncMap{
		"IncludeHTML": IncludeHTML,
	})

	tmpl, err := tmpl.Parse(`
<!DOCTYPE>
<html>
<body>
    {{ IncludeHTML "out.svg" }}
</body>
</html>
    `)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Error executing template: %v", err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/get", GetDiagram).Methods("GET")

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
