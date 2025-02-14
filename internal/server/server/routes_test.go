package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostMetric(t *testing.T) {
	type args struct {
		method       string
		url          string
		expectedCode int
		json         string
		checkBody    bool
	}
	type TestStruct struct {
		name string
		args args
	}

	tests := []TestStruct{
		{
			name: "Метод POST /update/ вернет 400  т.к. body пустое",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/",
				expectedCode: http.StatusBadRequest,
				json:         "",
			},
		},
		{
			name: "Метод POST /update/ вернет 404  т.к. не указано имя метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/",
				expectedCode: http.StatusNotFound,
				json:         "{}",
			},
		},
		{
			name: "Метод POST /update/ вернет 400  т.к. указан в MType указан неизвестный тип метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "xxx"}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 400  т.к. не указано значение метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "gauge"}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 400  т.к. указывать оба значения (Value и Delta) нельзя",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "counter", "value": 1.1, "delta": 1}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 400 т.к. для метрики counter надо указывать значение Delta",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "counter", "value": 1.1}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 400 т.к. для метрики gauge надо указывать значение Value",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "gauge", "delta": 1}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 400 т.к. значение метрики counter не может быть типа float",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "counter", "delta": 1.1}`,
			},
		},
		{
			name: "Метод POST /update/ вернет 200 т.к. данные для counter валидны и в ответ приходят данные о метриках",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusOK,
				json:         `{"id":"privet","type":"counter","delta":1}`,
				checkBody:    true,
			},
		},
		{
			name: "Метод POST /update/ вернет 200 т.к. данные для gauge валидны и в ответ приходят данные о метриках",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusOK,
				json:         `{"id":"privet","type":"gauge","value":1.1}`,
				checkBody:    true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			server := Setup()

			request := httptest.NewRequest(test.args.method, test.args.url, strings.NewReader(test.args.json))
			w := httptest.NewRecorder()
			server.router.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, test.args.expectedCode, res.StatusCode)

			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			responseString := string(body)

			if test.args.checkBody {
				assert.Equal(t, test.args.json, responseString)
			}

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
