package main

import "time"

type ApiResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		League   string `json:"league"`
		Date     string `json:"date"`
		Timezone string `json:"timezone"`
		Season   string `json:"season"`
	} `json:"parameters"`
	Errors  []any `json:"errors"`
	Results int   `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []struct {
		Fixture struct {
			ID        int       `json:"id"`
			Referee   string    `json:"referee"`
			Timezone  string    `json:"timezone"`
			Date      time.Time `json:"date"`
			Timestamp int       `json:"timestamp"`
			Periods   struct {
				First  int `json:"first"`
				Second int `json:"second"`
			} `json:"periods"`
			Venue struct {
				ID   any    `json:"id"`
				Name string `json:"name"`
				City string `json:"city"`
			} `json:"venue"`
			Status struct {
				Long    string `json:"long"`
				Short   string `json:"short"`
				Elapsed int    `json:"elapsed"`
				Extra   int    `json:"extra"`
			} `json:"status"`
		} `json:"fixture"`
		League struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Country   string `json:"country"`
			Logo      string `json:"logo"`
			Flag      any    `json:"flag"`
			Season    int    `json:"season"`
			Round     string `json:"round"`
			Standings bool   `json:"standings"`
		} `json:"league"`
		Teams struct {
			Home struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
			} `json:"home"`
			Away struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
			} `json:"away"`
		} `json:"teams"`
		Goals struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"goals"`
		Score struct {
			Halftime struct {
				Home int `json:"home"`
				Away int `json:"away"`
			} `json:"halftime"`
			Fulltime struct {
				Home int `json:"home"`
				Away int `json:"away"`
			} `json:"fulltime"`
			Extratime struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"extratime"`
			Penalty struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"penalty"`
		} `json:"score"`
	} `json:"response"`
}
