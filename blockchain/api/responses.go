package api

type Response struct {
	Payload map[string]interface{} `json:"payload"`
	Title   string                 `json:"title"`
}
