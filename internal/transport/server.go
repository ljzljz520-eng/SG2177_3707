package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"equipmentlending/internal/service"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func NewServer(business *service.Service) *Server {
	server := &Server{service: business, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/records", s.handleRecords)
	s.mux.HandleFunc("/records/", s.handleRecord)
	s.mux.HandleFunc("/sort", s.handleSort)
}

func (s *Server) Handler() http.Handler {
	return loggingMiddleware(s.mux)
}

func (s *Server) ListenAndServe(ctx context.Context, address string) error {
	server := &http.Server{Addr: address, Handler: s.Handler()}
	shutdown := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			server.Shutdown(context.Background())
		case <-shutdown:
		}
	}()
	err := server.ListenAndServe()
	close(shutdown)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Equipment-Service", "equipment-lending")
		next.ServeHTTP(writer, request)
	})
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func readJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func pathNumber(path string) string {
	return strings.TrimPrefix(path, "/records/")
}
