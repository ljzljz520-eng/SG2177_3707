package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Offset int
	Limit  int
}

func parsePagination(request *http.Request) Pagination {
	offset, _ := strconv.Atoi(request.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if limit > 500 {
		limit = 500
	}
	return Pagination{Offset: offset, Limit: limit}
}

func writeError(writer http.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, ErrorResponse{Code: code, Message: message})
}

func writeText(writer http.ResponseWriter, status int, text string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(text))
}

func acceptsJSON(request *http.Request) bool {
	accept := request.Header.Get("Accept")
	return accept == "" || strings.Contains(accept, "application/json") || strings.Contains(accept, "*/*")
}

func decodeBody(request *http.Request, target any) error {
	if request.Body == nil {
		return &json.SyntaxError{}
	}
	return json.NewDecoder(request.Body).Decode(target)
}

func methodAllowed(writer http.ResponseWriter, methods ...string) {
	writer.Header().Set("Allow", strings.Join(methods, ", "))
	writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method is not allowed")
}

func queryFlag(request *http.Request, key string) bool {
	value := strings.ToLower(strings.TrimSpace(request.URL.Query().Get(key)))
	return value == "1" || value == "true" || value == "yes"
}

func requestID(request *http.Request) string {
	if value := strings.TrimSpace(request.Header.Get("X-Request-ID")); value != "" {
		return value
	}
	return "local"
}

func setResponseHeaders(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("X-Request-ID", requestID(request))
	writer.Header().Set("Cache-Control", "no-store")
}
