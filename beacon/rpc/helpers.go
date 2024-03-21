package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func HandleError(w http.ResponseWriter, message string, code int) {
	errJson := &DefaultJsonError{
		Message: message,
		Code:    code,
	}
	WriteError(w, errJson)
}

// WriteJson writes the response message in JSON format.
func WriteJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", JsonMediaType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logrus.WithError(err).Error("Could not write response message")
	}
}

// WriteError writes the error by manipulating headers and the body of the final response.
func WriteError(w http.ResponseWriter, errJson HasStatusCode) {
	j, err := json.Marshal(errJson)
	if err != nil {
		logrus.WithError(err).Error("Could not marshal error message")
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(j)))
	w.Header().Set("Content-Type", JsonMediaType)
	w.WriteHeader(errJson.StatusCode())
	if _, err := io.Copy(w, io.NopCloser(bytes.NewReader(j))); err != nil {
		logrus.WithError(err).Error("Could not write error message")
	}
}
