package apiCalls

import (
	"fmt"
	"io"
	"net/http"
)

// GetBodyApiCall takes in an address, does an api call and returns the body and error if any
func GetBodyApiCall(addr string) ([]byte, error) {
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}
