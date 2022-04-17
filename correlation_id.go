package correlationid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	DefaultHeaderName = "Correlation-Id"
	ContextKey        = "CorrelationId"
)

type CorrelationIdMiddleware struct {
	HeaderName        string
	IncludeInResponse bool
	EnforceHeader     bool
	IdGenerator       func(ctx context.Context) string
}

func New() CorrelationIdMiddleware {
	return CorrelationIdMiddleware{
		HeaderName:        DefaultHeaderName,
		IncludeInResponse: true,
		EnforceHeader:     false,
		IdGenerator:       defultGenerator,
	}
}

func (m *CorrelationIdMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		headerName := m.getHeaderName()
		corrId := r.Header.Get(headerName)
		if corrId == "" {
			if m.EnforceHeader {
				http.Error(rw, fmt.Sprintf("%s header is required.", headerName), http.StatusBadRequest)
				return
			}

			corrId = m.generateId(r.Context())
		}

		if m.IncludeInResponse {
			rw.Header().Set(headerName, corrId)
		}

		updCtx := context.WithValue(r.Context(), ContextKey, corrId)
		next.ServeHTTP(rw, r.WithContext(updCtx))
	})
}

func FromContext(ctx context.Context) string {
	corrId := ctx.Value(ContextKey).(string)
	return corrId
}

func defultGenerator(ctx context.Context) string {
	return uuid.NewString()
}

func (m *CorrelationIdMiddleware) getHeaderName() string {
	if m.HeaderName == "" {
		return DefaultHeaderName
	}

	return m.HeaderName
}

func (m *CorrelationIdMiddleware) generateId(ctx context.Context) string {
	if m.IdGenerator != nil {
		return m.IdGenerator(ctx)
	}

	return defultGenerator(ctx)
}
