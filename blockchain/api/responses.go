package api

type HelloResponseGet struct {
	Payload string `json:"payload"`
	Title   string `json:"title"`
}

type HelloResponsePost struct {
	Payload struct {
		Payload string `json:"payload"`
		Title   string `json:"title"`
	} `json:"payload"`
	Title string `json:"title"`
}
