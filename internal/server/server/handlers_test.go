package server

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/setup"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
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

	values, err := setup.GetStartupValues(os.Args[1:])
	if err != nil {
		logger.Instance.Fatalw("Аргументы не прошли валидацию", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := Create(&values)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			server.BadRequestHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}
