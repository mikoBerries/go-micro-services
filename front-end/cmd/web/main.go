package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var myTemplate *template.Template

func init() {
	//parsing all template
	myTemplate = template.Must(template.ParseGlob("./templates/*"))
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renders(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port :80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}

}

// renders render all go html
func renders(w http.ResponseWriter, t string) {
	// partials := []string{
	// 	"./templates/base.layout.gohtml",
	// 	"./templates/header.partial.gohtml",
	// 	"./templates/footer.partial.gohtml",
	// }

	// var templateSlice []string
	// templateSlice = append(templateSlice, fmt.Sprintf("./templates/%s", t))

	// // for _, x := range partials {
	// templateSlice = append(templateSlice, partials...)
	// // }

	// tmpl, err := template.ParseFiles(templateSlice...)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	if err := myTemplate.ExecuteTemplate(w, t, nil); err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
