package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func middlewareCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (cfg *config) middlewareLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.logger.Info("received request",
			slog.String("ip", r.RemoteAddr),
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()))

		next.ServeHTTP(w, r)
	})
}

func (cfg *config) middlewareRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				// make go http server close this connection
				w.Header().Set("Connection", "close")
				cfg.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (cfg *config) middlewareRequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cfg.isAuthenticated(r) {
			cfg.sessionManager.Put(r.Context(), "flash", "You need to be logged in")
			cfg.redirect(w, r, "/users/login")
			return
		}

		// do not cache pages that requre auth
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (cfg *config) middlewareAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := cfg.sessionManager.GetString(r.Context(), "authenticatedUserID")
		if idStr == "" {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := uuid.Parse(idStr)
		if err != nil {
			cfg.sessionManager.Remove(r.Context(), "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}

		user, err := cfg.dbQueries.GetUserByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				cfg.sessionManager.Remove(r.Context(), "authenticatedUserID")
				next.ServeHTTP(w, r)
				return
			}
			cfg.serverError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func middlewarePreventCSRF(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()
	return cop.Handler(next)
}
