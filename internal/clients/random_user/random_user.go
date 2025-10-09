package randomuser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseURL = "https://randomuser.me/api/"
)

type RandomUser struct {
	client *http.Client

	countryCodes []string
	gender       *string
	results      *string
}

func New() *RandomUser {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(client *http.Client) *RandomUser {
	return &RandomUser{
		client:       client,
		countryCodes: make([]string, 0),
	}
}

func (g *RandomUser) WithCountryCode(code string) *RandomUser {
	g.countryCodes = append(g.countryCodes, code)
	return g
}

func (g *RandomUser) WithGender(gender string) *RandomUser {
	g.gender = &gender
	return g
}

func (g *RandomUser) WithResults(results int) *RandomUser {
	r := strconv.Itoa(results)
	g.results = &r
	return g
}

func (g *RandomUser) Generate() (Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Response{}, err
	}

	q := u.Query()
	q.Add("inc", "name,nat")
	q.Add("noinfo", "")
	if len(g.countryCodes) > 0 {
		q.Add("nat", strings.Join(g.countryCodes, ","))
	}
	if g.gender != nil {
		q.Add("gender", *g.gender)
	}
	if g.results != nil {
		q.Add("results", *g.results)
	}
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	// Make GET request
	resp, err := http.Get(u.String())
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		return Response{}, err
	}

	return data, nil
}
