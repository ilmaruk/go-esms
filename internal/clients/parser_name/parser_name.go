package parsername

import "net/http"

const (
	baseURL = "https://api.parser.name/"
)

type ParserName struct {
	apiKey string
	client *http.Client
}

func New(apiKey string) *ParserName {
	return NewWithClient(apiKey, http.DefaultClient)
}

func NewWithClient(apiKey string, client *http.Client) *ParserName {
	return &ParserName{
		apiKey: apiKey,
		client: client,
	}
}

func (pn ParserName) Generate() *Generate {
	return &Generate{
		apiKey: pn.apiKey,
		client: pn.client,
	}
}
