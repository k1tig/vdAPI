package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var VDHOLDER string

type Middleware func(http.Handler) http.Handler
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			r.URL.Path = strings.TrimSuffix(path, "/")
			http.Redirect(w, r, r.URL.Path, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key") // Get the API key from the header.

		if apiKey == "" || apiKey != VDHOLDER {

			http.Error(w, "Unauthorized ", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) // Proceed to the next handler if the API key is valid.
	})
}

func GetKeys() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	VDHOLDER = os.Getenv("VD_USERKEY")

}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
