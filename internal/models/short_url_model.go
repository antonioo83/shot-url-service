package models

type ShortURL struct {
	ID            string
	UserCode      int
	CorrelationID string
	Code          string
	OriginalURL   string
	ShortURL      string
}
