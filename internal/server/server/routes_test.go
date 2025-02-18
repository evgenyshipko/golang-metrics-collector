package server

import (
	"fmt"
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
			name: "Метод POST /update/ вернет 400  т.к. урл не валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод POST /update/dkdkd/ вернет 400  т.к. имя метрики не передано",
			args: args{
				method:       http.MethodPost,
				url:          "/update/dkdkd/",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/azaza/dhdh/1221/2323 вернет 404  т.к. такого урла нет",
			args: args{
				method:       http.MethodPost,
				url:          "/update/azaza/dhdh/1221/2323",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/gauge/ вернет 404  т.к. не указано имя и значение метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/counter/ вернет 404  т.к. не указано имя и значение метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/counter/",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/gauge/111 вернет 404  т.к. не указано имя метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/111",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/counter/111 вернет 404  т.к. не указано имя метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/counter/111",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод POST /update/counter/name/111 вернет 200 т.к. урл валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/counter/name/111",
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Метод POST /update/counter/name/111.1 вернет 400 т.к. значение метрики counter не может быть типа float",
			args: args{
				method:       http.MethodPost,
				url:          "/update/counter/name/111.1",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод POST /update/counter/111/111 вернет 200 т.к. в качестве имени метрики может быть число",
			args: args{
				method:       http.MethodPost,
				url:          "/update/counter/111/111",
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Метод POST /update/gauge/name/111 вернет 200 т.к. урл валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/name/111",
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Метод POST /update/gauge/name/111.1 вернет 200 т.к. метрика gauge поддерживает значения float",
			args: args{
				method:       http.MethodPost,
				url:          "/update/gauge/name/111.1",
				expectedCode: http.StatusOK,
			},
		},
		{
			name: "Метод POST /update/metric/name/111 вернет 400 т.к. у нас только два типа метрик - counter и gauge",
			args: args{
				method:       http.MethodPost,
				url:          "/update/metric/name/111",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод POST /update/ вернет 404  т.к. не указано имя метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusNotFound,
				json:         "{}",
			},
		},
		{
			name: "Метод POST /update/ вернет 400  т.к. указан в указан неизвестный тип метрики",
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
		// далее тестируем метод POST/ value
		{
			name: "Метод POST /value/ вернет 400  т.к. body пустое",
			args: args{
				method:       http.MethodPost,
				url:          "/value/",
				expectedCode: http.StatusBadRequest,
				json:         "",
			},
		},
		{
			name: "Метод POST /value/ вернет 404  т.к. не указано имя метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/value/",
				expectedCode: http.StatusNotFound,
				json:         "{}",
			},
		},
		{
			name: "Метод POST /value/ вернет 400  т.к. указан в указан неизвестный тип метрики",
			args: args{
				method:       http.MethodPost,
				url:          "/value/",
				expectedCode: http.StatusBadRequest,
				json:         `{"id": "privet", "type": "xxx"}`,
			},
		},
		{
			name: "Метод POST /value/ вернет 404 т.к. значения в базе нет",
			args: args{
				method:       http.MethodPost,
				url:          "/value/",
				expectedCode: http.StatusNotFound,
				json:         `{"id": "privet", "type": "counter"}`,
			},
		},
		// далее тестируем метод GET /value/
		{
			name: "Метод GET /value/ вернет 405 т.к. роут есть, но метод GET не принимает",
			args: args{
				method:       http.MethodGet,
				url:          "/value/",
				expectedCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Метод GET /value/azaz вернет 404 т.к. такого роута нет",
			args: args{
				method:       http.MethodGet,
				url:          "/value/azaz",
				expectedCode: http.StatusNotFound,
			},
		},
		{
			name: "Метод GET /value/azaz/name вернет 400 т.к. передан неизвестный тип метрики",
			args: args{
				method:       http.MethodGet,
				url:          "/value/azaz/name",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод GET /value/counter/name вернет 404 т.к. имени с метрикой name тем в базе",
			args: args{
				method:       http.MethodGet,
				url:          "/value/counter/name",
				expectedCode: http.StatusNotFound,
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

			fmt.Println(responseString)

			if test.args.checkBody {
				assert.Equal(t, test.args.json, responseString)
			}

			defer res.Body.Close()
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
