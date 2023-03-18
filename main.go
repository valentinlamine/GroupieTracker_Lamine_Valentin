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
	"structure/structure"
	"time"
)

func main() {
	fmt.Println("Server started on port 80 : http://localhost")
	//Chargement des fichiers CSS
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	//Chargement des ASSETS
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	//Chargement des fichiers JS
	js := http.FileServer(http.Dir("scripts"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", js))
	//Gestion des templates
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/search", SearchHandler)
	http.ListenAndServe(":80", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //Si l'url n'est pas valide
		ErrorHandler(w, r, 404) //On affiche une erreur 404
		return
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" { //Si l'url n'est pas valide
		ErrorHandler(w, r, 404) //On affiche une erreur 404
		return
	}
	start := time.Now() //Début du timer
	fmt.Print("\n")
	//Définition des templates
	t, _ := template.ParseFiles("templates/recherche.html")
	t2, _ := template.ParseFiles("templates/resultat.html")
	//Regex pour récupérer l'id dans l'url
	reg := regexp.MustCompile(`\/search\?id=(?P<id>\d+)`)
	if r.URL.RequestURI() == "/search" { //Si l'url est /search
		if r.Method == "POST" { //Si on a effectué une recherche
			request := RequestByName(r.FormValue("search"), r.FormValue("type"))
			result := RequestHandler(request)
			t.Execute(w, result) //On affiche les résultats
		} else { //Si on a pas effectué de recherche
			t.Execute(w, nil) //On affiche la page de recherche
		}
	} else if reg.MatchString(r.URL.RequestURI()) { //Si on accède à un résultat
		id := reg.FindStringSubmatch(r.URL.RequestURI())[1] //On récupère l'id
		request := RequestById(id)                          //On récupère les informations du résultat
		result := RequestHandler(request)                   //On traite les informations
		t2.Execute(w, result)                               //On affiche le résultat
	} else { //Si l'url n'est pas valide
		t.Execute(w, nil) //On affiche la page de recherche
	}
	fmt.Println("SearchHandler  | Request took", time.Since(start))
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == 404 { //Si l'erreur est 404
		t, _ := template.ParseFiles("templates/404.html")
		t.Execute(w, nil)
	}
}

func RequestByName(search string, media string) structure.Response { //Fonction de requête par nom
	url := "https://itunes.apple.com/search?country=FR&term=" + url.QueryEscape(search) + "&media=" + url.QueryEscape(media) //Création de l'url

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête
	res, err := http.DefaultClient.Do(req)     //Envoi de la requête
	if err != nil {                            //Si il y a une erreur
		panic(err) //On arrête le programme
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request structure.Response //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	for i := 0; i < len(request.Results); i++ { //Boucle de traitement des résultats
		if request.Results[i].Kind == "song" { //Si le résultat est une chanson
			if request.Results[i].IsStreamable == false { //Si la chanson n'est pas streamable
				request.Results = append(request.Results[:i], request.Results[i+1:]...) //On supprime le résultat
				i--
			}
		}
	}

	if request.ResultCount != 0 { //Si il y a des résultats
		fmt.Println("RequestByName  | Sucessful request with", request.ResultCount, "results of type", request.Results[0].Kind)
	}
	return request //Retourne la variable temporaire
}

func RequestById(id string) structure.Response { //Fonction de requête par id
	url := "https://itunes.apple.com/lookup?country=FR&id=" + id //Création de l'url

	req, _ := http.NewRequest("GET", url, nil) //Création de la requête
	res, err := http.DefaultClient.Do(req)     //Envoi de la requête
	if err != nil {                            //Si il y a une erreur
		panic(err) //On arrête le programme
	}

	defer res.Body.Close()          //Fermeture du body
	body, _ := io.ReadAll(res.Body) //Lecture du body

	var request structure.Response //Création de la variable temporaire
	json.Unmarshal(body, &request) //Désérialisation du JSON

	if request.ResultCount != 0 { //Si il y a des résultats
		fmt.Println("RequestById    | Sucessful request with", request.ResultCount, "results of type", request.Results[0].Kind)
	}
	return request //Retourne la variable temporaire
}

func RequestHandler(Request structure.Response) structure.Result {
	//Création de la structure de retour
	var Result structure.Result
	Result.ResultCount = Request.ResultCount //On copie le nombre de résultats
	Result.Results = make([]struct {         //On crée la liste des résultats
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
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate.Format("02/01/2006")
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, true)
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = "There is no description for this song"
		case "feature-movie": //Si le résultat est un film
			Result.Results[i].Type = "movie"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "A movie is not an album"
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate.Format("02/01/2006")
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, false)
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		case "ebook": //Si le résultat est un livre
			Result.Results[i].Type = "ebook"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = Request.Results[i].TrackName
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = "An ebook is not an album"
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate.Format("02/01/2006")
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].TrackViewURL
			Result.Results[i].Duration = "No data"
			Result.Results[i].Price = Request.Results[i].Price
			Result.Results[i].Description = Request.Results[i].Description
		case "tv-episode": //Si le résultat est un épisode de série
			Result.Results[i].Type = "tv-episode"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate.Format("02/01/2006")
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, false)
			Result.Results[i].Price = Request.Results[i].TrackPrice
			Result.Results[i].Description = Request.Results[i].LongDescription
		case "music-video": //Si le résultat est une vidéo musicale
			Result.Results[i].Type = "music-video"
			Result.Results[i].Id = Request.Results[i].TrackID
			Result.Results[i].Title = IsExplicit(Request.Results[i].TrackName, Request.Results[i].TrackExplicitness)
			Result.Results[i].Artist = Request.Results[i].ArtistName
			Result.Results[i].Album = Request.Results[i].CollectionName
			Result.Results[i].ReleaseDate = Request.Results[i].ReleaseDate.Format("02/01/2006")
			Result.Results[i].PreviewImage = PreviewUpscaling(Request.Results[i].ArtworkURL100)
			Result.Results[i].PreviewContent = Request.Results[i].PreviewURL
			Result.Results[i].Duration = FormatDuration(Request.Results[i].TrackTimeMillis, true)
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

func IsExplicit(title, explicit string) string {
	if explicit == "explicit" {
		return title + " ⓔ"
	} else {
		return title
	}
}

func FormatDuration(duration int, isSong bool) string {
	if isSong {
		minutes := duration / 60000
		seconds := (duration % 60000) / 1000
		if seconds < 10 {
			return strconv.Itoa(minutes) + ":0" + strconv.Itoa(seconds)
		}
		return strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)
	} else {
		hours := duration / 3600000
		minutes := (duration % 3600000) / 60000
		if minutes < 10 {
			return strconv.Itoa(hours) + "h0" + strconv.Itoa(minutes) + "m"
		}
		return strconv.Itoa(hours) + "h" + strconv.Itoa(minutes) + "m"
	}
}
