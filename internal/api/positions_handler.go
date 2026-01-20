package api

import (
	"encoding/json"
	"htrr-apis/internal/store"
	"htrr-apis/internal/utils"
	"log"
	"net/http"
)

type PositionHandler struct {
	logger *log.Logger
	store  store.PostgresPosition
}

func NewPositionHandler(logger *log.Logger, positionStore store.PostgresPosition) *PositionHandler {
	return &PositionHandler{
		logger: logger,
		store:  positionStore,
	}
}

type createPositionRequest struct {
	Title string
}

func (h *PositionHandler) HandleCreatePosition(w http.ResponseWriter, r *http.Request) {
	var body createPositionRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.Printf("ERROR: decode body: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if body.Title == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "title is required"})
		return
	}

	pos := &store.Position{
		Title: body.Title,
	}

	err = h.store.Create(pos)
	if err != nil {
		h.logger.Printf("Error: Create position failed :%v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"position": pos})
}

func (h *PositionHandler) HandleUpdatePosition(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdUrlParams(r)

	if err != nil {
		h.logger.Printf("ERROR: parse id via params: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "id is not valid"})
		return
	}

	var body createPositionRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.Printf("ERROR: decode body failed: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	pos, err := h.store.GetById(id)
	if err != nil {
		h.logger.Printf("ERROR: get position by id failed %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if pos == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "position does not exits"})
		return
	}

	err = h.store.Update(*pos)
	if err != nil {
		h.logger.Printf("ERROR: update position: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"position": pos})
}

func (h *PositionHandler) HandleGetPositionById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdUrlParams(r)
	if err != nil {
		h.logger.Printf("ERROR: parse id via params: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	pos, err := h.store.GetById(id)
	if err != nil {
		h.logger.Printf("ERROR: cannot get position by id :%v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if pos == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "position does not exist!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"position": pos})
}
