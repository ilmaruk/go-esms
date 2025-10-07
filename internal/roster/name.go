package roster

import (
	"fmt"

	"github.com/ilmaruk/go-esms/internal"
	parsername "github.com/ilmaruk/go-esms/internal/clients/parser_name"
	"github.com/rainycape/unidecode"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const parserNameApiKey = "99b06ec70c6c0d7c40cd8dcbe5dd46c9"

type parserNameName struct {
	Name      string `json:"name"`
	NameASCII string `json:"name_ascii"`
}

type parserNameObject struct {
	Person struct {
		FirstName parserNameName `json:"firstname"`
		LastName  parserNameName `json:"lastname"`
	} `json:"name"`
	Country struct {
		Code      string `json:"country_code"`
		CodeAlpha string `json:"country_code_alpha"`
	} `json:"country"`
}

type parserNameGenerate struct {
	Data []parserNameObject `json:"data"`
}

func generatePersons(qty int) ([]internal.Person, error) {
	if qty > 25 {
		return nil, fmt.Errorf("invalid quantity; max is 25")
	}

	parserNameClient := parsername.New(parserNameApiKey)
	generate := parserNameClient.Generate().
		WithGender("m").
		WithResults(qty)
	resp, err := generate.Do()
	if err != nil {
		return nil, err
	}

	persons := make([]internal.Person, 0, qty)
	for _, o := range resp.(parsername.GenerateResponse).Data {
		person := internal.Person{
			FirstName: normaliseName(o.Person.FirstName.Name, o.Country.Code),
			LastName:  normaliseName(o.Person.LastName.Name, o.Country.Code),
			Country:   o.Country.CodeAlpha,
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func normaliseName(name, country string) string {
	t := getLanguageTag(country)
	name = capitalise(name, t)
	return unidecode.Unidecode(name)
}

func capitalise(s string, t language.Tag) string {
	caser := cases.Title(t)

	return caser.String(s)
}

func getLanguageTag(code string) language.Tag {
	// TODO: user a map instead
	switch code {
	case "RU":
		return language.Russian
	case "UA":
		return language.Ukrainian
	case "BG":
		return language.Bulgarian
	case "CY":
		fallthrough
	case "GR":
		return language.Greek
	default:
		return language.English
	}
}
