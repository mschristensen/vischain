package api

// APIResponse describes the format of a response
// received from the API
type APIResponse struct {
	Payload map[string]interface{} `json:"payload"`
	Title   string                 `json:"title"`
}

func (r *APIResponse) FromMap(m map[string]interface{}) {
	result := APIResponse{
		Payload: m["payload"].(map[string]interface{}),
		Title:   m["title"].(string),
	}
	*r = result
}

// OKResponse describes the shape of a HTTP 200
// response we send the the API
type OKResponse struct {
	Code int8
}
