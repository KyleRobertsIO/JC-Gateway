package azure

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponseError struct {
	Error ErrorDetails `json:"error"`
}
