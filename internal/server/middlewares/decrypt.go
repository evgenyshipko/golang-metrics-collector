package middlewares

import (
	"bytes"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/decrypt"
	"io"
	"net/http"
)

func DecryptMiddleware(cryptoPrivateKeyPath string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header.Get(`X-Encrypted`) != `true` || cryptoPrivateKeyPath == "" {
				next.ServeHTTP(w, r)
				return
			}

			privateKey, err := decrypt.LoadPrivateKey(cryptoPrivateKeyPath)
			if err != nil {
				logger.Instance.Warnw("Failed to load private key: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r.Body)
			if err != nil {
				logger.Instance.Warnw(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			decryptedMsg, err := decrypt.DecryptWithPrivateKey(buf.Bytes(), privateKey)
			if err != nil {
				logger.Instance.Warnw("Failed to decrypt message: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decryptedMsg))
			r.ContentLength = int64(len(decryptedMsg))

			next.ServeHTTP(w, r)

		})
	}
}
