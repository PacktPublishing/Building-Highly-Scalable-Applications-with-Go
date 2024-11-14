package middleware

import "net/http"

// ResponsePeeker wraps http.ResponseWriter to capture the status code
type ResponsePeeker struct {
	http.ResponseWriter
	Status int
}

func (p *ResponsePeeker) WriteHeader(code int) {
	p.Status = code
	p.ResponseWriter.WriteHeader(code)
}
