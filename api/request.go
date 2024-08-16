package httptransport

import (
	"encoding/json"
	"strings"
)

type CreateUser struct {
	name    string `json:"name"`
	address string `json:"address"`
}

const contentType string = "application/json"

func DecodeRequestCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	// Lấy giá trị của tham số "name" từ query string
	if !strings.Contains(r.Header.Get("Content-type", contentType)){
		return nil, errors.New("Unsupported content type")
	}
	req := CreateUser{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("Malformed entity")
	}
	// Tạo một đối tượng Request từ giá trị "name" trích xuất

	return req, nil
}