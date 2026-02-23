package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	mbAPI      = "https://musicbrainz.org/ws/2"
	caAPI      = "https://coverartarchive.org"
	tmdbAPI    = "https://api.themoviedb.org/3"
	tmdbImg    = "https://image.tmdb.org/t/p/w500"
	outputJSON = "data/favorites_enriched.json"
	coverDir   = "static/img/covers"
)

var (
	tmdbKey = os.Getenv("TMDB_API_KEY")
	client  = &http.Client{Timeout: 15 * time.Second}
)

type Input struct {
	Albums []struct {
		Title    string
		Artist   string
		Director string
	}
	Movies []struct {
		Title    string
		Director string
	}
	Shows []struct{ Title string }
}

type Album struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Year     string `json:"year"`
	Cover    string `json:"cover"`
}

type Media struct {
	Title    string  `json:"title"`
	Year     string  `json:"year"`
	Poster   string  `json:"poster"`
	Overview string  `json:"overview"`
	Rating   float64 `json:"rating"`
	Director string  `json:"director,omitempty"`
}

type Output struct {
	Albums []Album `json:"albums"`
	Movies []Media `json:"movies"`
	Shows  []Media `json:"shows"`
}

func get(rawURL string, target any) error {
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "FavoritesFetcher/1.0 (personal-site)")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func download(rawURL, path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	os.MkdirAll(filepath.Dir(path), 0755)
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "FavoritesFetcher/1.0 (personal-site)")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func fetchAlbum(title, artist string) Album {
	q := url.QueryEscape(fmt.Sprintf(`release:"%s" AND artist:"%s"`, title, artist))
	var mb struct {
		Releases []struct {
			ID   string `json:"id"`
			Date string `json:"date"`
		} `json:"releases"`
	}
	get(fmt.Sprintf("%s/release/?query=%s&fmt=json&limit=5", mbAPI, q), &mb)
	time.Sleep(1100 * time.Millisecond)

	if len(mb.Releases) == 0 {
		return Album{Title: title, Artist: artist}
	}

	rel := mb.Releases[0]
	year := ""
	if len(rel.Date) >= 4 {
		year = rel.Date[:4]
	}

	var ca struct {
		Images []struct {
			Front      bool              `json:"front"`
			Image      string            `json:"image"`
			Thumbnails map[string]string `json:"thumbnails"`
		} `json:"images"`
	}
	get(fmt.Sprintf("%s/release/%s", caAPI, rel.ID), &ca)
	time.Sleep(500 * time.Millisecond)

	cover := ""
	for _, img := range ca.Images {
		if img.Front {
			if t, ok := img.Thumbnails["500"]; ok {
				cover = t
			} else {
				cover = img.Image
			}
			break
		}
	}

	local := ""
	if cover != "" {
		slug := slugify(artist + "-" + title)
		path := fmt.Sprintf("%s/albums/%s.jpg", coverDir, slug)
		if download(cover, path) == nil {
			local = fmt.Sprintf("/img/covers/albums/%s.jpg", slug)
		}
	}

	return Album{Title: title, Artist: artist, Year: year, Cover: firstNonEmpty(local, cover)}
}

func fetchTMDB(title, mediaType, director string) Media {
	if tmdbKey == "" {
		return Media{Title: title, Director: director}
	}

	var result struct {
		Results []struct {
			Title       string  `json:"title"`
			Name        string  `json:"name"`
			ReleaseDate string  `json:"release_date"`
			FirstAirDate string `json:"first_air_date"`
			PosterPath  string  `json:"poster_path"`
			Overview    string  `json:"overview"`
			VoteAverage float64 `json:"vote_average"`
		} `json:"results"`
	}
	get(fmt.Sprintf("%s/search/%s?api_key=%s&query=%s", tmdbAPI, mediaType, tmdbKey, url.QueryEscape(title)), &result)

	if len(result.Results) == 0 {
		return Media{Title: title, Director: director}
	}

	r := result.Results[0]
	name := r.Title
	date := r.ReleaseDate
	if mediaType == "tv" {
		name = r.Name
		date = r.FirstAirDate
	}
	year := ""
	if len(date) >= 4 {
		year = date[:4]
	}

	poster := ""
	if r.PosterPath != "" {
		imgURL := tmdbImg + r.PosterPath
		slug := slugify(title)
		path := fmt.Sprintf("%s/%ss/%s.jpg", coverDir, mediaType, slug)
		if download(imgURL, path) == nil {
			poster = fmt.Sprintf("/img/covers/%ss/%s.jpg", mediaType, slug)
		} else {
			poster = imgURL
		}
	}

	overview := r.Overview
	if len(overview) > 200 {
		overview = overview[:200]
	}

	return Media{
		Title:    name,
		Year:     year,
		Poster:   poster,
		Overview: overview,
		Rating:   r.VoteAverage,
		Director: director,
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func main() {
	f, err := os.Open("favorites.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var input Input
	yaml.NewDecoder(f).Decode(&input)

	out := Output{}

	for _, a := range input.Albums {
		out.Albums = append(out.Albums, fetchAlbum(a.Title, a.Artist))
	}
	for _, m := range input.Movies {
		out.Movies = append(out.Movies, fetchTMDB(m.Title, "movie", m.Director))
	}
	for _, s := range input.Shows {
		out.Shows = append(out.Shows, fetchTMDB(s.Title, "tv", ""))
	}

	os.MkdirAll(filepath.Dir(outputJSON), 0755)
	outFile, err := os.Create(outputJSON)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(out)
}