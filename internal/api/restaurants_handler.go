package api

import (
	"database/sql"
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

type bulkDeleteRestaurantRequest struct {
	IDs      []string `json:"ids"`
	Strategy string   `json:"strategy"`
}

func (r *bulkDeleteRestaurantRequest) Validate() error {
	if len(r.IDs) == 0 {
		return errors.New("ids array is required and cannot be empty")
	}

	validStrategies := map[string]bool{
		"atomic":      true,
		"partial":     true,
		"best_effort": true,
	}

	if r.Strategy == "" {
		r.Strategy = "atomic"
	}

	if !validStrategies[r.Strategy] {
		return errors.New("invalid strategy: must be 'atomic', 'partial', or 'best_effort'")
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

func (h *RestaurantHandler) HandleGetRestaurantById(w http.ResponseWriter, r *http.Request) {
	paramsId, err := utils.GetIdUrlParams(r)
	if err != nil {
		h.logger.Printf("ERROR: GetIdUrlParams %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id"})
		return
	}

	restaurant, err := h.store.GetRestaurantById(paramsId)
	if err != nil {
		h.logger.Printf("ERROR: GetRestaurantById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}

	if restaurant == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"restaurant": restaurant})
}

func (h *RestaurantHandler) HandleUpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	rId, err := utils.GetIdUrlParams(r)
	if err != nil {
		h.logger.Printf("ERROR: GetIdViaUrl %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	existingRestaurant, err := h.store.GetRestaurantById(rId)
	if err != nil {
		h.logger.Printf("ERROR: GetRestaurantById %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}

	if existingRestaurant == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "cannot find id"})
		return
	}

	type updateRestaurantRequest struct {
		Name     *string `json:"name"`
		Address  *string `json:"address"`
		IsActive *bool   `json:"is_active"`
		Phone    *string `json:"phone"`
	}

	var rqBody updateRestaurantRequest
	err = json.NewDecoder(r.Body).Decode(&rqBody)
	if err != nil {
		h.logger.Printf("ERROR: decode json failed %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if rqBody.Name != nil {
		existingRestaurant.Name = *rqBody.Name
	}
	if rqBody.IsActive != nil {
		existingRestaurant.IsActive = *rqBody.IsActive
	}
	if rqBody.Address != nil {
		existingRestaurant.Address = *rqBody.Address
	}
	if rqBody.Phone != nil {
		existingRestaurant.Phone = *rqBody.Phone
	}

	err = h.store.Update(existingRestaurant)
	if err != nil {
		h.logger.Printf("ERROR: updateRestaurant %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"restaurant": existingRestaurant})
}

func (h *RestaurantHandler) HandleDeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdUrlParams(r)
	if err != nil {
		h.logger.Printf("ERROR: GetIdUrlParams, %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.store.Delete(id)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "not found id"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "deleted~"})

}

func (h *RestaurantHandler) HandleBulkDeleteRestaurants(w http.ResponseWriter, r *http.Request) {
	var req bulkDeleteRestaurantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: Failed to decode request body, %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Printf("ERROR: Validation failed, %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	switch req.Strategy {
	case "atomic":
		h.handleAtomicDelete(w, req.IDs)
	case "partial":
		h.handlePartialDelete(w, req.IDs)
	case "best_effort":
		h.handleBestEffortDelete(w, req.IDs)
	}
}

func (h *RestaurantHandler) handleAtomicDelete(w http.ResponseWriter, ids []string) {
	_, err := h.store.BulkDeleteAtomic(ids)
	if err != nil {
		if err.Error() == "Invalid id format" {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "one or more ids not found"})
			return
		}
		h.logger.Printf("ERROR: BulkDeleteAtomic failed, %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "deleted successfully"})
}

func (h *RestaurantHandler) handlePartialDelete(w http.ResponseWriter, ids []string) {
	result, err := h.store.BulkDeletePartial(ids)
	if err != nil {
		h.logger.Printf("ERROR: BulkDeletePartial failed, %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	status := http.StatusOK
	if result.FailedCount > 0 && result.DeletedCount == 0 {
		status = http.StatusNotFound
	} else if result.FailedCount > 0 {
		status = http.StatusPartialContent
	}

	utils.WriteJSON(w, status, utils.Envelope{
		"deleted_count": result.DeletedCount,
		"failed_count":  result.FailedCount,
		"deleted_ids":   result.DeletedIDs,
		"failed_ids":    result.FailedIDs,
		"message":       "bulk delete completed with details",
	})
}

func (h *RestaurantHandler) handleBestEffortDelete(w http.ResponseWriter, ids []string) {
	count, err := h.store.BulkDeleteBestEffort(ids)
	if err != nil {
		h.logger.Printf("ERROR: BulkDeleteBestEffort failed, %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"deleted_count": count,
		"message":       "deleted successfully",
	})
}
