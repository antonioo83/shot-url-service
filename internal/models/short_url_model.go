package models

type ShortURL struct {
	ID            string
	UserCode      int
	CorrelationId string
	Code          string
	OriginalURL   string
	ShortURL      string
}
