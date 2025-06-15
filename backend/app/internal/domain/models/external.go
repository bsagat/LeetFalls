package models

type Character struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image"`
	Quote    string `json:"quote"`
}
