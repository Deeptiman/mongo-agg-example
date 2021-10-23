package main

import (
	"time"
)

type Movie struct {
	ID     string `json:"_id"`
	Awards struct {
		Nominations int    `json:"nominations"`
		Text        string `json:"text"`
		Wins        int    `json:"wins"`
	} `json:"awards"`
	Cast      []string `json:"cast"`
	Countries []string `json:"countries"`
	Directors []string `json:"directors"`
	Genres    []string `json:"genres"`
	Imdb      struct {
		ID     int     `json:"id"`
		Rating float64 `json:"rating"`
		Votes  int     `json:"votes"`
	} `json:"imdb"`
	Languages        []string  `json:"languages"`
	Lastupdated      string    `json:"lastupdated"`
	NumMflixComments int       `json:"num_mflix_comments"`
	Plot             string    `json:"plot"`
	Poster           string    `json:"poster"`
	Rated            string    `json:"rated"`
	Released         time.Time `json:"released"`
	Runtime          int       `json:"runtime"`
	Title            string    `json:"title"`
	Type             string    `json:"type"`
	Writers          []string  `json:"writers"`
	Year             int       `json:"year"`
}

type Comment struct {
	ID      string    `json:"_id"`
	Date    time.Time `json:"date"`
	Email   string    `json:"email"`
	MovieID string    `json:"movie_id"`
	Name    string    `json:"name"`
	Text    string    `json:"text"`
}