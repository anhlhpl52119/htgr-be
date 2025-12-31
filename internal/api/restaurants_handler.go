package api

import (
	"encoding/json"
	"go-todo-apis/internal/store"
	"go-todo-apis/internal/utils"
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

func (h *RestaurantHandler) HandleCreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var restaurant store.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		h.logger.Printf("ERROR: decoding HandleCreateRestaurant: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
		return
	}

	if restaurant.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "name is required"})
		return
	}

	restaurant.IsActive = false

	err = h.store.Create(&restaurant)
	if err != nil {
		h.logger.Printf("ERROR: creating failed: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"restaurant": restaurant})
}

func (h *RestaurantHandler) HandleSearchRestaurant(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	page, err := strconv.Atoi(queries.Get("page"))
	if err != nil || page == 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(queries.Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := store.SearchRestaurantParams{
		Page:     page,
		PageSize: pageSize,
		Name:     queries.Get("name"),
	}

	list, total, err := h.store.Search(search)
	if err != nil {
		h.logger.Printf("ERROR: search failed: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"restaurants": list,
		"page":        page,
		"pageSize":    pageSize,
		"total":       total})
}
