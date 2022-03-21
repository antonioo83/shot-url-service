package models

type ShortURL struct {
	ID          string
	UserCode    int
	Code        string
	OriginalURL string
	ShortURL    string
}
