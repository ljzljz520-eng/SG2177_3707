package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"equipmentlending/internal/persistence"
	"equipmentlending/internal/service"
)

func TestHTTPHandlers(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	business := service.New(store)
	if err := business.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	server := NewServer(business)
	body, _ := json.Marshal(service.CreateRequest{EquipmentNumber: "H-1", Name: "Hose", Borrower: "Zoe", BorrowDate: "2026-06-01"})
	request := httptest.NewRequest(http.MethodPost, "/records", bytes.NewReader(body))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create endpoint returned %d", response.Code)
	}
	get := httptest.NewRequest(http.MethodGet, "/records?borrower=zoe", nil)
	getResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(getResponse, get)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("query endpoint returned %d", getResponse.Code)
	}
}
