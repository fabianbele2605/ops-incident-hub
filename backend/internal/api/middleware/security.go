package middleware

import "net/http"

// SecurityHeadersMiddleware agrega headers de seguridad HTTP
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevenir MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Prevenir clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// Habilitar XSS protection en browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Forzar HTTPS (solo en producción)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		next.ServeHTTP(w, r)
	})
}
