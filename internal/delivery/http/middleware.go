package http

import (
	"context"
	"fmt"
	httpBase "net/http"
	"time"

	"TransactionTest/internal/domain"
	"TransactionTest/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDMiddleware добавляет уникальный RequestID в контекст запроса
func RequestIDMiddleware(next httpBase.Handler) httpBase.Handler {
	return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), logger.RequestID, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware логирует информацию о запросах
func LoggingMiddleware(log logger.Logger) func(httpBase.Handler) httpBase.Handler {
	return func(next httpBase.Handler) httpBase.Handler {
		return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
			start := time.Now()

			// Создаем response writer wrapper для захвата статуса
			wrapped := &responseWriter{ResponseWriter: w, statusCode: httpBase.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			log.Info(r.Context(), "HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

// RecoveryMiddleware обрабатывает паники и возвращает 500 ошибку
func RecoveryMiddleware(log logger.Logger) func(httpBase.Handler) httpBase.Handler {
	return func(next httpBase.Handler) httpBase.Handler {
		return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(r.Context(), "Panic recovered", zap.Any("error", err))

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(httpBase.StatusInternalServerError)

					jsonResponse := fmt.Sprintf(
						`{"error":"%s","code":"%s","message":"%s"}`,
						"Internal Server Error",
						domain.CodeInternal,
						"An unexpected error occurred",
					)
					w.Write([]byte(jsonResponse))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wrapper для захвата статуса ответа
type responseWriter struct {
	httpBase.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
