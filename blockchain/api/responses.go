package api

// Response describes the format of a response
// received from the API
type Response struct {
	Payload map[string]interface{} `json:"payload"`
	Title   string                 `json:"title"`
}

// FromMap accepts an empty interface map describing an Response
// (parsed from a JSON) and stores its parsed contents in the Response.
func (r *Response) FromMap(m map[string]interface{}) error {
	result := Response{
		Payload: m["payload"].(map[string]interface{}),
		Title:   m["title"].(string),
	}
	*r = result
	return nil
}

// PeerResponse describes the shape of a HTTP response
// received from a peer which we return to the API
type PeerResponse map[string]interface{}
