package structure

import (
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
