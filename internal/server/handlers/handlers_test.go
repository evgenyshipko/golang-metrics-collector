package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequestHandler(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Проверим что хендлер возвратит 404",
			args: args{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			BadRequestHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}
