package middlewares

import "net/http"

type Claim struct {
	Username string `json:"username"`
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: verify token...

		next.ServeHTTP(w, r)
	})
}

func verifyToken(w http.ResponseWriter, r *http.Request) (string, *Claim, error) {
	// TODO: refactor
	w.Header().Add("Vary", "Authorization")
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return "", nil, nil
	}
	return "", nil, nil
}
