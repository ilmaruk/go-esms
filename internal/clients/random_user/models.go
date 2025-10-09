package randomuser

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type result struct {
	Name        Name   `json:"name"`
	Nationality string `json:"nat"`
}

type Response struct {
	Results []result `json:"results"`
}
