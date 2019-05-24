package websocket

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/logger"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
)

func middlewareAllowedAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultAccessControl.CanAccess(r.Context()) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Init(ctx context.Context, config *Config) *Websocket {
	events := repository.Events()

	go func() {
		for {
			if err := eq.feedSessions(ctx, events, store); err != nil {
				logger.Default().Error("session event feed error", zap.Error(err))
			}
		}
	}()
	eq.store(ctx, events)

	return Websocket{}.New(config)
}

func (ws Websocket) ApiServerRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/websocket", func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// @todo make access control injectable
					if !service.DefaultAccessControl.CanAccess(r.Context()) {
						http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
						return
					}

					next.ServeHTTP(w, r)
				})
			})
			r.Get("/", ws.Open)
		})
	})
}
