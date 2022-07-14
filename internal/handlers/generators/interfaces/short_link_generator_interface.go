package interfaces

type ShortLinkGenerator interface {
	GetShortURL(originalURL string, host string, newBaseURL string) (string, string, error)
}
