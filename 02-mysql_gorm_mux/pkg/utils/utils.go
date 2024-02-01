package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// ParseBody 解析请求主体
func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err = json.Unmarshal(body, x); err != nil {
			return
		}
	}
}
