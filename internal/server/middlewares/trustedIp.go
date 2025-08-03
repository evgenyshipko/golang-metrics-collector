package middlewares

import (
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
	"net"
	"net/http"
)

func TrustedIpMiddleware(trustedSubnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if trustedSubnet == nil {
				next.ServeHTTP(w, r)
				return
			}

			logger.Instance.Info("МИДЛВАРА ТРАСТЕД РАБОТАЕТ")

			ipStr := r.Header.Get("X-Real-IP")
			if ipStr == "" {
				http.Error(w, "X-Real-IP header is required", http.StatusForbidden)
				return
			}

			ip := net.ParseIP(ipStr)
			if ip == nil {
				http.Error(w, "Invalid IP address in X-Real-IP header", http.StatusForbidden)
				return
			}

			// Проверяем принадлежность IP к доверенной подсети
			if !trustedSubnet.Contains(ip) {
				http.Error(w, "Access denied: IP not in trusted subnet", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
