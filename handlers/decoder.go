package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// DecodePostRequest is used to parse post body data and apply it to the struct
// passed into v.
func DecodePostRequest[T any](r *http.Request, v *T) (*T, error) {
	// Ensure this is a POST request
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Check if the body is empty
	if len(body) == 0 {
		return nil, errors.New("request body is empty")
	}

	// Parse the JSON body into reqStruct struct
	if err := json.Unmarshal(body, v); err != nil {
		return nil, err
	}

	return v, nil
}
