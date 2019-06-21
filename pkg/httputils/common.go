package httputils

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

const (
	// CtxRequestID is the request id context key const set to the context
	CtxRequestID = iota
	// CtxLogger is the logger object set to the context
	CtxLogger
)

var _log *zap.Logger

func log() *zap.Logger {
	if _log == nil {
		_log = zap.L().Named("httputils")
	}
	return _log
}

// RequestID returns the request_id value in context
func RequestID(r *http.Request) string {
	ctx := r.Context()
	if ctx != nil {
		return ctx.Value(CtxRequestID).(string)
	}
	return ""
}

// ReadJSONBody returns the JSON object from the HTTP request
func ReadJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		log().Error("parse json error", zap.String("request_id", RequestID(r)), zap.Error(err))
		WriteJSONWithCode(w, r, errors.New("malformed json text"), http.StatusBadRequest)
	}
	return err
}

// WriteJSON writes an object to the HTTP response
func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	WriteJSONWithCode(w, r, v, -1)
}

// WriteJSONWithCode is same to WriteJSON, but uses the custom HTTP status code
func WriteJSONWithCode(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if code != -1 { // code -1 means don't write status code
		w.WriteHeader(code)
	}

	// parse json
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log().Error("parse json error", zap.String("request_id", RequestID(r)), zap.Error(err))
		WriteJSONWithCode(w, r, errors.New("server error"), http.StatusInternalServerError)
		return
	}

	j = append(j, byte('\n')) // append additional line break to json bytes
	if _code, err := w.Write(j); err != nil {
		log().Error("http body write error", zap.String("request_id", RequestID(r)), zap.Int("code", _code), zap.Error(err))
	}
}
