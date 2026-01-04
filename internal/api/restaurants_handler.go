package api

import (
	"encoding/json"
	"errors"
	"htrr-apis/internal/store"
	"htrr-apis/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type RestaurantHandler struct {
	logger *log.Logger
	store  store.RestaurantStore
}

func NewRestaurantHandler(logger *log.Logger, store store.RestaurantStore) *RestaurantHandler {
	return &RestaurantHandler{
		logger: logger,
		store:  store,
	}
}

type registerRestaurantRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func (r *registerRestaurantRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (h *RestaurantHandler) HandleCreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var reqBody registerRestaurantRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		h.logger.Printf("ERROR: decoding HandleCreateRestaurant: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
		return
	}

	err = reqBody.Validate()
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	restaurant := &store.Restaurant{
		Name:     reqBody.Name,
		Address:  reqBody.Address,
		Phone:    reqBody.Phone,
		IsActive: false,
	}

	err = h.store.Create(restaurant)
	if err != nil {
		h.logger.Printf("ERROR: creating failed: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"restaurant": restaurant})
}

func (h *RestaurantHandler) HandleSearchRestaurant(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	// 1. Parse & Validate input
	req := store.SearchRestaurantParams{
		Name:     queries.Get("name"),
		Page:     parseIntOrDefault(queries.Get("page"), 1),
		PageSize: parseIntOrDefault(queries.Get("page_size"), 10),
	}

	list, total, err := h.store.Search(req)
	if err != nil {
		h.logger.Printf("ERROR: search failed: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"restaurants": list,
		"metadata": map[string]any{
			"current_page":  req.Page,
			"page_size":     req.PageSize,
			"total_records": total,
		}})
}

func parseIntOrDefault(value string, def int) int {
	v, err := strconv.Atoi(value)
	if err != nil || v <= 0 {
		return def
	}
	return v
}
