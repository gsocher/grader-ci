package route

import (
	"encoding/json"
	"net/http"
)

type errorMsg struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func writeError(rw http.ResponseWriter, statusCode int, e error) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	m := &errorMsg{
		Status: statusCode,
		Msg:    e.Error(),
	}

	b, _ := json.Marshal(m)
	rw.Write(b)

	return nil
}

func writeOk(rw http.ResponseWriter, b []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
