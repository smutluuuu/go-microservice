package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}

	// data.BrokerURL = os.Getenv("BROKER_URL")
	// data.BrokerURL = "http://broker-service.info"
	brokerURL := "http://broker-service.info"

	// brokerURL := os.Getenv("BROKER_URL")
	// if brokerURL == "" {
	// 	fmt.Println("Broker URL is not set")
	// 	brokerURL = "http://broker-service.info"
	// } else {
	// 	fmt.Printf("Broker URL is : %s\n", brokerURL)
	// }

	data.BrokerURL = brokerURL

	fmt.Println("Broker URL: ", data.BrokerURL)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
