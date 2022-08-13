package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_handlerF1(t *testing.T) {
	type want struct {
		code     int
		Location string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "positive test #1",
			want: want{
				code:     http.StatusOK,
				Location: "google.com",
			},
			request: "http://localhost:8080/get/id?c0e2fd12-1105-4cbf-b8d8-99881602ad25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlerF1)
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.Location, res.Header.Get("Location"))

		})
	}
}
