package parsername

type Name struct {
	Name      string `json:"name"`
	NameASCII string `json:"name_ascii"`
}

type Data struct {
	Person struct {
		FirstName Name `json:"firstname"`
		LastName  Name `json:"lastname"`
	} `json:"name"`
	Country struct {
		Code      string `json:"country_code"`
		CodeAlpha string `json:"country_code_alpha"`
	} `json:"country"`
}

type GenerateResponse struct {
	Data []Data `json:"data"`
}
