package transport

import (
	"context"
	"net/http"
	"strings"

	"equipmentlending/internal/reporting"
	"equipmentlending/internal/service"
)

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.service.Store().Healthy(request.Context()); err != nil {
		writeJSON(writer, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleRecords(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		borrower := request.URL.Query().Get("borrower")
		status := request.URL.Query().Get("status")
		records, err := s.service.Query(request.Context(), borrower, status)
		if err != nil {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(writer, http.StatusOK, records)
	case http.MethodPost:
		var input service.CreateRequest
		if err := readJSON(request, &input); err != nil {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		record, err := s.service.Create(request.Context(), input)
		if err != nil {
			writeJSON(writer, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(writer, http.StatusCreated, record)
	default:
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (s *Server) handleRecord(writer http.ResponseWriter, request *http.Request) {
	number := pathNumber(request.URL.Path)
	if number == "" || strings.Contains(number, "/") {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": "record not found"})
		return
	}
	if request.Method != http.MethodDelete {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	record, err := s.service.Delete(request.Context(), number)
	if err != nil {
		writeJSON(writer, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, record)
}

func (s *Server) handleSort(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	mode := request.URL.Query().Get("mode")
	records, err := s.service.Sort(context.Background(), mode)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"records": records, "report": reporting.FormatRecords(records)})
}
