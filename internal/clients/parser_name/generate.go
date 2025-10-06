package parsername

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Generate struct {
	apiKey string
	client *http.Client

	countryCode *string
	gender      *string
	results     *string
}

func (g *Generate) WithCountryCode(code string) *Generate {
	g.countryCode = &code
	return g
}

func (g *Generate) WithGender(gender string) *Generate {
	g.gender = &gender
	return g
}

func (g *Generate) WithResults(results int) *Generate {
	r := strconv.Itoa(results)
	g.results = &r
	return g
}

func (g *Generate) Do() (any, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("api_key", g.apiKey)
	q.Add("endpoint", "generate")
	if g.countryCode != nil {
		q.Add("country_code", *g.countryCode)
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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data GenerateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
