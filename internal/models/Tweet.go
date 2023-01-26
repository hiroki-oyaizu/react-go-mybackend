package models

type Tweet struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Tweet string `json:"tweet"`
}
