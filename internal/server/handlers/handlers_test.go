package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostMetric(t *testing.T) {
	type args struct {
		method       string
		url          string
		expectedCode int
	}
	type TestStruct struct {
		name string
		args args
	}

	tests := []TestStruct{
		{
			name: "Метод GET / не доступен",
			args: args{
				method:       http.MethodGet,
				url:          "/",
				expectedCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Метод GET /update/ не доступен",
			args: args{
				method:       http.MethodGet,
				url:          "/update/",
				expectedCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Метод POST /update/ вернет 400  т.к. урл не валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод POST /update/dkdkd/ вернет 400  т.к. урл не валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/dkdkd/",
				expectedCode: http.StatusBadRequest,
			},
		},
		{
			name: "Метод POST /update/azaza/dhdh/1221/2323 вернет 400  т.к. урл не валиден",
			args: args{
				method:       http.MethodPost,
				url:          "/update/azaza/dhdh/1221/2323",
				expectedCode: http.StatusBadRequest,
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(test.args.method, test.args.url, nil)
			w := httptest.NewRecorder()
			PostMetric(w, request)

			res := w.Result()
			assert.Equal(t, test.args.expectedCode, res.StatusCode)
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
