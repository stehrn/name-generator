package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/docker/docker/pkg/namesgenerator"
)

// templateData data to render into template
type templateData struct {
	Message       string
}

type handlerContext struct {
	template   *template.Template
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, ctx.templateData(""))
	case "POST":
		name := namesgenerator.GetRandomName(0)
		ctx.template.Execute(w, ctx.templateData(name))
	}
	return nil
}

func (ctx handlerContext) templateData(message string) *templateData {
	return &templateData{message}
}

func main() {
	log.Print("Starting client")

	templatePath := os.Getenv("TEMPLATE_PATH")
	clientTemplate := filepath.Join(templatePath, "./index.tmpl")
	log.Printf("Loading template from: %s", clientTemplate)

	ctx := &handlerContext{
		template:   template.Must(template.ParseFiles(clientTemplate))}

	handler := errorHandler(ctx.handle)
	http.HandleFunc("/", handler)

	port := port()
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// port get port form env, default to 8080
func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	return port
}

// errorHandler handler for errors to reduce boilerplate code
func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

