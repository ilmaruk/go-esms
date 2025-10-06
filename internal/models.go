package internal

import "strings"

type Person struct {
	FirstName string
	LastName  string
	Country   string
}

func (p Person) String() string {
	return strings.Join([]string{p.FirstName, p.LastName}, " ")
}

func (p Person) Short() string {
	fn := ""
	parts := strings.Split(p.FirstName, " ")
	for _, p := range parts {
		fn += strings.ToUpper(p[0:1])
	}
	return strings.Join([]string{fn, p.LastName}, " ")
}
