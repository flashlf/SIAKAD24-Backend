package middleware

import "net/http"

// StripSlashes trims one trailing slash from the request path before it
// reaches the wrapped handler, mirroring chi.Middleware.StripSlashes without
// issuing an HTTP redirect. The root path "/" is left untouched.
func StripSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			r.URL.Path = path[:len(path)-1]
		}
		next.ServeHTTP(w, r)
	})
}
