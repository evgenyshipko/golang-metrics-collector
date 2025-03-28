package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"github.com/evgenyshipko/golang-metrics-collector/internal/server/middlewares/utils"
	"net/http"
)

func ValidateSHA256(hashKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			hashFromHeader := r.Header.Get(`HashSHA256`)

			var hashPointer *string

			if hashFromHeader != "" && hashKey != "" && r.Body != nil && r.Body != http.NoBody {

				requestBody, err := utils.GetBodyAndRestore(r)
				if err != nil {
					logger.Instance.Warnw("GetBodyAndRestore", "err", err)
					http.Error(w, "GetBodyAndRestore ошибка", http.StatusBadRequest)
					return
				}

				h := hmac.New(sha256.New, []byte(hashKey))
				h.Write([]byte(requestBody))
				hash := hex.EncodeToString(h.Sum(nil))

				if hashFromHeader != hash {
					logger.Instance.Warnw("Хеши не совпадают", "hashFromHeader", hashFromHeader, "hash", hash)
					http.Error(w, "отказано в доступе", http.StatusForbidden)
					return
				}

				hashPointer = &hash
			}

			next.ServeHTTP(w, r)

			if hashPointer != nil {
				w.Header().Set("HashSHA256", *hashPointer)
				logger.Instance.Infow("ValidateSHA256", "Requiest Headers", w.Header())
			}
		})
	}
}
