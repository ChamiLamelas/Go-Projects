package main

const SHORTEN_ENDPOINT = "/urlshortener/shorten"

const EXPAND_ENDPOINT = "/urlshortener/expand/"

const ANALYTICS_ENDPOINT = "/urlshortener/analytics/"

type ShortenRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type ShortenResponse struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

type ExpandResponse struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

type AnalyticsResponse struct {
	Url        string `json:"url"`
	Alias      string `json:"alias"`
	Expansions int    `json:"expansions"`
}
