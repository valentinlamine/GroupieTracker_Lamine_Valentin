package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		Price                   float64   `json:"price"`
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
		Description             string    `json:"description,omitempty"`
		CollectionArtistViewURL string    `json:"collectionArtistViewUrl,omitempty"`
	} `json:"results"`
}

type Result struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		Type           string    `json:"type"`
		Title          string    `json:"title"`
		Artist         string    `json:"artist"`
		Album          string    `json:"album"`
		ReleaseDate    time.Time `json:"releaseDate"`
		Explicit       bool      `json:"explicit"`
		PreviewImage   string    `json:"previewImage"`
		PreviewContent string    `json:"previewContent"`
		Price          float64   `json:"price"`
		Description    string    `json:"description"`
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
		request := API_GET(r.FormValue("search"), r.FormValue("type"))
		result := API_Handling(request)
		t.Execute(w, result)
	} else {
		t.Execute(w, nil)
	}

}

func API_GET(search string, media string) Response {
	url := "https://itunes.apple.com/search?term=" + url.QueryEscape(search) + "&media=" + url.QueryEscape(media) //Création de l'url

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête

	res, err := http.DefaultClient.Do(req) //Envoi de la requête

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request Response           //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	fmt.Println("Requette effectuée avec ", request.ResultCount, " résultats de type ", request.Results[0].Kind)
	fmt.Println(len(request.Results))

	return request //Retourne la variable temporaire
}

func API_Handling(Request Response) Result {
	//Création de la structure de retour
	var Result Result
	Result.ResultCount = Request.ResultCount
	Result.Results = make([]struct {
		Type           string    `json:"type"`
		Title          string    `json:"title"`
		Artist         string    `json:"artist"`
		Album          string    `json:"album"`
		ReleaseDate    time.Time `json:"releaseDate"`
		Explicit       bool      `json:"explicit"`
		PreviewImage   string    `json:"previewImage"`
		PreviewContent string    `json:"previewContent"`
		Price          float64   `json:"price"`
		Description    string    `json:"description"`
	}, len(Request.Results))

	//Boucle de traitement des résultats
	for i := 0; i < len(Request.Results); i++ {
		switch Request.Results[i].Kind {
		case "song": //Si le résultat est une chanson
			Result.Results[i].Type = "song"
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate
			Result.Results[i].Explicit = Request.Results[i].TrackExplicitness == "explicit"
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = "Not description available for song"
		case "feature-movie": //Si le résultat est un film
			Result.Results[i].Type = "movie"
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "Not an album"
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate
			Result.Results[i].Explicit = Request.Results[i].TrackExplicitness == "explicit"
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		case "ebook": //Si le résultat est un livre
			Result.Results[i].Type = "ebook"
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "Not an album"
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate
			Result.Results[i].Explicit = false
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].TrackViewURL
			Result.Results[i].Price = Request.Results[i].Price
			Result.Results[i].Description = Request.Results[i].Description
		case "tv-episode": //Si le résultat est un épisode de série
			Result.Results[i].Type = "tv-episode"
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate
			Result.Results[i].Explicit = Request.Results[i].TrackExplicitness == "explicit"
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		default: //Si le résultat n'est pas reconnu
			fmt.Println("Type non reconnu")
		}
	}
	fmt.Println("Requete traitée avec", len(Result.Results), "résultats")
	return Result
}

func PreviewUpscaling(preview string) string {
	//Input : "https://is4-ssl.mzstatic.com/image/thumb/Music124/v4/f3/ee/b3/f3eeb3ff-ca32-273a-15aa-709bdfa64367/mzi.izwiyqez.jpg/100x100bb.jpg"
	//Output : "https://is4-ssl.mzstatic.com/image/thumb/Music124/v4/f3/ee/b3/f3eeb3ff-ca32-273a-15aa-709bdfa64367/mzi.izwiyqez.jpg/1000x1000bb.jpg"
	preview = strings.Replace(preview, "100x100bb", "1000x1000bb", 1)
	return preview
}
