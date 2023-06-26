package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/jacobtie/rating-party/server/internal/platform/werrors"
)

type errResponse struct {
	RequestID string `json:"requestId,omitempty"`
	Error     string `json:"error"`
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, werrors.ErrNotFound) {
		respondError(ctx, w, err, http.StatusNotFound)
		return
	}
	if errors.Is(err, werrors.ErrBadRequest) {
		respondError(ctx, w, err, http.StatusBadRequest)
		return
	}
	if errors.Is(err, werrors.ErrUnauthorized) {
		respondError(ctx, w, err, http.StatusUnauthorized)
		return
	}
	if errors.Is(err, werrors.ErrForbidden) {
		respondError(ctx, w, err, http.StatusForbidden)
		return
	}
	respondError(ctx, w, err, http.StatusInternalServerError)
}

func respondError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	response := errResponse{
		Error: err.Error(),
	}
	v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
	if ok {
		v.Error = true
		if v.RequestID != "" {
			response.RequestID = v.RequestID
		}
	}
	Respond(ctx, w, response, code)
}

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) {
	v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
	if ok {
		v.StatusCode = code
	}
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		respondError(ctx, w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}
