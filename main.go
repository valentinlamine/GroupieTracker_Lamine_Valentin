package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Server started on port 80 : http://localhost")
	//Chargement des fichiers CSS
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	//Chargement des ASSETS
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/css/", http.StripPrefix("/assets/", assets))
	//Gestion des templates
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func API_TEST() {
	url := "https://randomuser.me/api/"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	err := ioutil.WriteFile("test.json", body, 0644)
	if err != nil {
		panic(err)
	}

}
