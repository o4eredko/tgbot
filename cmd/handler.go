package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"tg-bot/internal/sheets"
)

type errResponse struct {
	Error error `json:"error"`
}

type Handler struct {
	service *sheets.Service
}

func NewHandler(service *sheets.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/cities", h.getCities).Methods(http.MethodGet)
	r.HandleFunc("/cities/{city}/organizations", h.getOrganizations).Methods(http.MethodGet)
	r.HandleFunc("/organizations/{organization}", h.getOrganization).Methods(http.MethodGet)
}

func (h *Handler) getCities(w http.ResponseWriter, r *http.Request) {
	const tableCities = "Організації"

	resp, err := h.service.GetHeaders(r.Context(), tableCities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, err)

		return
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, fmt.Errorf("marshal response: %w", err))

		return
	}

	log.Println("get cities")
	w.Write(rawResp)
}

func (h *Handler) getOrganizations(w http.ResponseWriter, r *http.Request) {
	const tableOrganizations = "Організації"

	city := mux.Vars(r)["city"]

	resp, err := h.service.GetColumn(r.Context(), tableOrganizations, city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, err)

		return
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, fmt.Errorf("marshal response: %w", err))

		return
	}

	log.Println("get organizations")
	w.Write(rawResp)
}

func (h *Handler) getOrganization(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["organization"]

	resp, err := h.service.GetTableRows(r.Context(), organization)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, err)

		return
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.writeErr(w, fmt.Errorf("marshal response: %w", err))

		return
	}

	log.Println("get organization: ", organization)
	w.Write(rawResp)
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	resp := errResponse{Error: err}
	rawResp, _ := json.Marshal(resp)

	w.Write(rawResp)
}
