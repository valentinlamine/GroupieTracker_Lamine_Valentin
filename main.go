package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		WrapperType             string    `json:"wrapperType"`
		Kind                    string    `json:"kind"`
		ArtistID                int       `json:"artistId"`
		CollectionID            int       `json:"collectionId"`
		TrackID                 int       `json:"trackId"`
		ArtistName              string    `json:"artistName"`
		CollectionName          string    `json:"collectionName"`
		TrackName               string    `json:"trackName"`
		CollectionCensoredName  string    `json:"collectionCensoredName"`
		TrackCensoredName       string    `json:"trackCensoredName"`
		ArtistViewURL           string    `json:"artistViewUrl"`
		CollectionViewURL       string    `json:"collectionViewUrl"`
		TrackViewURL            string    `json:"trackViewUrl"`
		PreviewURL              string    `json:"previewUrl"`
		ArtworkURL30            string    `json:"artworkUrl30"`
		ArtworkURL60            string    `json:"artworkUrl60"`
		ArtworkURL100           string    `json:"artworkUrl100"`
		CollectionPrice         float64   `json:"collectionPrice"`
		TrackPrice              float64   `json:"trackPrice"`
		ReleaseDate             time.Time `json:"releaseDate"`
		CollectionExplicitness  string    `json:"collectionExplicitness"`
		TrackExplicitness       string    `json:"trackExplicitness"`
		DiscCount               int       `json:"discCount"`
		DiscNumber              int       `json:"discNumber"`
		TrackCount              int       `json:"trackCount"`
		TrackNumber             int       `json:"trackNumber"`
		TrackTimeMillis         int       `json:"trackTimeMillis"`
		Country                 string    `json:"country"`
		Currency                string    `json:"currency"`
		PrimaryGenreName        string    `json:"primaryGenreName"`
		IsStreamable            bool      `json:"isStreamable,omitempty"`
		CollectionArtistName    string    `json:"collectionArtistName,omitempty"`
		CollectionArtistID      int       `json:"collectionArtistId,omitempty"`
		ContentAdvisoryRating   string    `json:"contentAdvisoryRating,omitempty"`
		ShortDescription        string    `json:"shortDescription,omitempty"`
		LongDescription         string    `json:"longDescription,omitempty"`
		CollectionArtistViewURL string    `json:"collectionArtistViewUrl,omitempty"`
	} `json:"results"`
}

func main() {
	fmt.Println("Server started on port 80 : http://localhost")
	//Chargement des fichiers CSS
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	//Chargement des ASSETS
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	//Gestion des templates
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/search", SearchHandler)
	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/recherche.html")
	if r.Method == "POST" {
		request := API_GET(r.FormValue("search"))
		t.Execute(w, request.Results)
	} else {
		t.Execute(w, nil)
	}

}

func API_GET(search string) Response {
	url := "https://itunes.apple.com/search?term=" + url.QueryEscape(search)

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête

	res, err := http.DefaultClient.Do(req) //Envoi de la requête

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request Response           //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	fmt.Println("Requette effectuée avec ", request.ResultCount, " résultats") //Affichage du nombre de résultats
	fmt.Println(len(request.Results))

	return request //Retourne la variable temporaire
}

/*
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
*/
