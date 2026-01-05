package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	indentJson, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	indentJson = append(indentJson, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(indentJson)
	return nil
}

func GetOffset(pageNo, pageSize *int) (limit, offset int) {
	defaultPage := 1
	defaultSize := 10

	p := defaultPage
	if pageNo != nil && *pageNo > 0 {
		p = *pageNo
	}

	s := defaultSize
	if pageSize != nil && *pageSize > 0 {
		s = *pageSize
	}
	limit = s
	offset = (p - 1) * s

	return limit, offset
}

func GetURLQuery(key string, r *http.Request) string {
	return r.URL.Query().Get(key)
}

func GetIdUrlParams(r *http.Request) (string, error) {
	paramId := chi.URLParam(r, "id")
	if paramId == "" {
		return "", errors.New("invalid id params")
	}
	return paramId, nil
}
