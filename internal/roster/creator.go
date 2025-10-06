package roster

import "fmt"

func CreateRoster(qty int) error {
	persons, err := generatePersons(qty)
	if err != nil {
		return err
	}

	for _, p := range persons {
		fmt.Println(p, p.Short())
	}

	return nil
}
