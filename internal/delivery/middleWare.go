package delivery

import (
	"context"
	"net/http"

	"TransactionTest/internal/logger"
	"github.com/google/uuid"
)

type LoggerMiddleWare struct {

}

func NewMiddlerWare() *LoggerMiddleWare {
	return &LoggerMiddleWare{}
}

func (lmw *LoggerMiddleWare) loggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), logger.RequestID, uuid.New().String())
		r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}