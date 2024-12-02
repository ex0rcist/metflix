package middleware

import (
	"net"
	"net/http"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
)

// Ensure incoming request is from a trsuted subnet
func FilterUntrustedRequest(next http.Handler, trustedSubnet *net.IPNet) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if trustedSubnet == nil { // skip middleware entirely
			next.ServeHTTP(w, r)
			return
		}

		clientIP := net.ParseIP(r.Header.Get("X-Real-IP"))

		if !trustedSubnet.Contains(clientIP) {
			logging.LogError(entities.UntrustedSubnetError(clientIP))
			http.Error(w, "", http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
