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

type Middleware struct {
	HeaderName        string
	IncludeInResponse bool
	EnforceHeader     bool
	IdGenerator       func() string
}

func New() Middleware {
	return Middleware{
		HeaderName:        DefaultHeaderName,
		IncludeInResponse: true,
		EnforceHeader:     false,
		IdGenerator:       defaultGenerator,
	}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		headerName := m.getHeaderName()
		corrId := r.Header.Get(headerName)
		if corrId == "" {
			if m.EnforceHeader {
				http.Error(rw, fmt.Sprintf("%s header is required.", headerName), http.StatusBadRequest)
				return
			}

			corrId = m.generateId()
		}

		if m.IncludeInResponse {
			rw.Header().Set(headerName, corrId)
		}

		updCtx := context.WithValue(r.Context(), ContextKey, corrId)
		next.ServeHTTP(rw, r.WithContext(updCtx))
	})
}

func FromContext(ctx context.Context) string {
	corrId, ok := ctx.Value(ContextKey).(string)
	if ok {
		return corrId
	}
	return ""
}

func defaultGenerator() string {
	return uuid.NewString()
}

func (m *Middleware) getHeaderName() string {
	if m.HeaderName == "" {
		return DefaultHeaderName
	}

	return m.HeaderName
}

func (m *Middleware) generateId() string {
	if m.IdGenerator != nil {
		return m.IdGenerator()
	}

	return defaultGenerator()
}
