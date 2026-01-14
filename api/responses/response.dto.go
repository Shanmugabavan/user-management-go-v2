package responses

type Response struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}
