package api_error

type DataError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type HttpError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Error HttpError `json:"error"`
}
