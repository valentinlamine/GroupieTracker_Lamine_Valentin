package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
		Type           string  `json:"type"`
		Id             int     `json:"id"`
		Title          string  `json:"title"`
		Artist         string  `json:"artist"`
		Album          string  `json:"album"`
		ReleaseDate    string  `json:"releaseDate"`
		PreviewImage   string  `json:"previewImage"`
		PreviewContent string  `json:"previewContent"`
		Duration       string  `json:"duration"`
		Price          float64 `json:"price"`
		Description    string  `json:"description"`
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
	start := time.Now() //Début du timer
	fmt.Print("\n")
	t, _ := template.ParseFiles("templates/recherche.html")
	t2, _ := template.ParseFiles("templates/resultat.html")
	reg := regexp.MustCompile(`\/search\?id=(?P<id>\d+)`)
	if r.URL.RequestURI() == "/search" {
		if r.Method == "POST" {
			request := RequestByName(r.FormValue("search"), r.FormValue("type"))
			result := RequestHandler(request)
			t.Execute(w, result)
		} else {
			t.Execute(w, nil)
		}
	} else if reg.MatchString(r.URL.RequestURI()) { ///search?id=648817663
		id := reg.FindStringSubmatch(r.URL.RequestURI())[1]
		request := RequestById(id)
		result := RequestHandler(request)
		t2.Execute(w, result)
	} else {
		t.Execute(w, nil)
	}
	fmt.Println("SearchHandler  | Request took", time.Since(start))
}

func RequestByName(search string, media string) Response {
	url := "https://itunes.apple.com/search?country=FR&term=" + url.QueryEscape(search) + "&media=" + url.QueryEscape(media) //Création de l'url

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête

	res, err := http.DefaultClient.Do(req) //Envoi de la requête

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request Response           //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	for i := 0; i < len(request.Results); i++ {
		if request.Results[i].Kind == "song" {
			if request.Results[i].IsStreamable == false {
				request.Results = append(request.Results[:i], request.Results[i+1:]...)
				i--
			}
		}
	}

	if request.ResultCount != 0 {
		fmt.Println("RequestByName  | Sucessful request with", request.ResultCount, "results of type", request.Results[0].Kind)
	}
	return request //Retourne la variable temporaire
}

func RequestById(id string) Response {
	url := "https://itunes.apple.com/lookup?country=FR&id=" + id //Création de l'url

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête

	res, err := http.DefaultClient.Do(req) //Envoi de la requête

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request Response           //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	if request.ResultCount != 0 {
		fmt.Println("RequestById    | Sucessful request with", request.ResultCount, "results of type", request.Results[0].Kind)
	}
	return request //Retourne la variable temporaire
}

func RequestHandler(Request Response) Result {
	//Création de la structure de retour
	var Result Result
	Result.ResultCount = Request.ResultCount
	Result.Results = make([]struct {
		Type           string  `json:"type"`
		Id             int     `json:"id"`
		Title          string  `json:"title"`
		Artist         string  `json:"artist"`
		Album          string  `json:"album"`
		ReleaseDate    string  `json:"releaseDate"`
		PreviewImage   string  `json:"previewImage"`
		PreviewContent string  `json:"previewContent"`
		Duration       string  `json:"duration"`
		Price          float64 `json:"price"`
		Description    string  `json:"description"`
	}, len(Request.Results))

	//Boucle de traitement des résultats
	for i := 0; i < len(Request.Results); i++ {
		switch Request.Results[i].Kind {
		case "song": //Si le résultat est une chanson
			Result.Results[i].Type = "song"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = FormatDate(Request.Results[i].ReleaseDate)
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, "song")
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = "Not description available for song"
		case "feature-movie": //Si le résultat est un film
			Result.Results[i].Type = "movie"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "Not an album"
			Result.Results[i].ReleaseDate = FormatDate(Request.Results[i].ReleaseDate)
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, "movie")
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		case "ebook": //Si le résultat est un livre
			Result.Results[i].Type = "ebook"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "Not an album"
			Result.Results[i].ReleaseDate = FormatDate(Request.Results[i].ReleaseDate)
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].TrackViewURL
			Result.Results[i].Duration = "Not a song"
			Result.Results[i].Price = Request.Results[i].Price
			Result.Results[i].Description = Request.Results[i].Description
		case "tv-episode": //Si le résultat est un épisode de série
			Result.Results[i].Type = "tv-episode"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = FormatDate(Request.Results[i].ReleaseDate)
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, "movie")
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		case "music-video": //Si le résultat est une vidéo musicale
			Result.Results[i].Type = "music-video"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = FormatDate(Request.Results[i].ReleaseDate)
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, "song")
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		default: //Si le résultat n'est pas reconnu
			fmt.Println("Type non reconnu : ", Request.Results[i].Kind)
		}
	}
	fmt.Println("RequestHandler | Succesfully handled request")
	return Result
}

func PreviewUpscaling(preview string) string {
	//Input : "https://is4-ssl.mzstatic.com/image/thumb/Music124/v4/f3/ee/b3/f3eeb3ff-ca32-273a-15aa-709bdfa64367/mzi.izwiyqez.jpg/100x100bb.jpg"
	//Output : "https://is4-ssl.mzstatic.com/image/thumb/Music124/v4/f3/ee/b3/f3eeb3ff-ca32-273a-15aa-709bdfa64367/mzi.izwiyqez.jpg/1000x1000bb.jpg"
	preview = strings.Replace(preview, "100x100bb", "1000x1000bb", 1)
	return preview
}

func FormatDate(date time.Time) string {
	//Input : 2021-01-01T00:00:00Z
	//Output : 01/01/2021
	return date.Format("02/01/2006")
}

func IsExplicit(title, explicit string) string {
	if explicit == "explicit" {
		return title + " ⓔ"
	} else {
		return title
	}
}

func FormatDuration(duration int, format string) string {
	if format == "song" {
		minutes := duration / 60000
		seconds := (duration % 60000) / 1000
		if seconds < 10 {
			return strconv.Itoa(minutes) + ":0" + strconv.Itoa(seconds)
		}
		return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
	} else if format == "movie" {
		hours := duration / 3600000
		minutes := (duration % 3600000) / 60000
		if minutes < 10 {
			return strconv.Itoa(hours) + "h0" + strconv.Itoa(minutes) + "m"
		}
		return strconv.Itoa(hours) + "h" + strconv.Itoa(minutes) + "m"
	} else {
		return "This never happens because i'm perfect"
	}
}
