package roster

import (
	"github.com/ilmaruk/go-esms/internal"
	parsername "github.com/ilmaruk/go-esms/internal/clients/parser_name"
	randomuser "github.com/ilmaruk/go-esms/internal/clients/random_user"
	"github.com/rainycape/unidecode"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PersonGenerator interface {
	Generate(int) ([]internal.Person, error)
}

const parserNameApiKey = "99b06ec70c6c0d7c40cd8dcbe5dd46c9"

type parserNameGenerator struct {
	client *parsername.ParserName
}

func newParserNameGenerator(apiKey string) *parserNameGenerator {
	return &parserNameGenerator{
		client: parsername.New(parserNameApiKey),
	}
}

func (g *parserNameGenerator) Generate(qty int) ([]internal.Person, error) {
	persons := make([]internal.Person, 0, qty)

	missing := qty
	for missing > 0 {
		results := missing
		if results > 25 {
			results = 25
		}

		generate := g.client.Generate().
			WithGender("m").
			WithResults(qty)
		resp, err := generate.Do()
		if err != nil {
			return nil, err
		}

		for _, o := range resp.(parsername.GenerateResponse).Data {
			person := internal.Person{
				FirstName: normaliseName(o.Person.FirstName.Name, o.Country.Code),
				LastName:  normaliseName(o.Person.LastName.Name, o.Country.Code),
				Country:   o.Country.CodeAlpha,
			}
			persons = append(persons, person)
		}

		missing -= results
	}

	return persons, nil
}

type randomUserGenerator struct {
	client *randomuser.RandomUser
}

func newRandomUserGenerator() *randomUserGenerator {
	return &randomUserGenerator{
		client: randomuser.New(),
	}
}

func (g *randomUserGenerator) Generate(qty int) ([]internal.Person, error) {
	persons := make([]internal.Person, 0, qty)

	generate := g.client.WithGender("male").
		WithCountryCode("AU").
		WithCountryCode("BR").
		WithCountryCode("CA").
		WithCountryCode("CH").
		WithCountryCode("DE").
		WithCountryCode("DK").
		WithCountryCode("ES").
		WithCountryCode("FI").
		WithCountryCode("FR").
		WithCountryCode("GB").
		WithCountryCode("IE").
		WithCountryCode("MX").
		WithCountryCode("NL").
		WithCountryCode("NO").
		WithCountryCode("NZ").
		WithCountryCode("RS").
		WithCountryCode("TR").
		WithCountryCode("UA").
		WithCountryCode("US").
		WithResults(qty)
	resp, err := generate.Generate()
	if err != nil {
		return nil, err
	}

	for _, o := range resp.Results {
		person := internal.Person{
			FirstName: normaliseName(o.Name.First, o.Nationality),
			LastName:  normaliseName(o.Name.Last, o.Nationality),
			Country:   o.Nationality,
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
