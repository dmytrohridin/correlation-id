// Package correlationid provides handler for getting correlationid header from request or generate new and set it to request context.
package correlationid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	// Header name. Can be overridden in Middlware struct.
	DefaultHeaderName = "Correlation-Id"
	// Key used for correlation id value in context.
	ContextKey = "CorrelationId"
)

// Middlware struct with handler settings.
type Middleware struct {
	// Used for correlation id. Default "Correlation-Id".
	HeaderName string
	// Used for indicating should correlation id be included in response headers. Default "true".
	IncludeInResponse bool
	// Used for enforcing client include correlation id header in request. If "true" and header is absent - middlware will return BadRequest status code.
	// Default "false".
	EnforceHeader bool
	// Used for correlation id generation if id is not included in request. By default uses github.com/google/uuid.
	IdGenerator func() string
}

// Initialize default value of middleware. By default uses github.com/google/uuid for correlation id value.
func New() Middleware {
	return Middleware{
		HeaderName:        DefaultHeaderName,
		IncludeInResponse: true,
		EnforceHeader:     false,
		IdGenerator:       defaultGenerator,
	}
}

// Handles http request and proceed correlation id based on Middlware settings.
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

		updCtx := WithCorrelationId(r.Context(), corrId)
		next.ServeHTTP(rw, r.WithContext(updCtx))
	})
}

// Returns correlation id from context if present. Otherwise "".
func FromContext(ctx context.Context) string {
	corrId, ok := ctx.Value(ContextKey).(string)
	if ok {
		return corrId
	}
	return ""
}

// Sets correlation id to context
func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, ContextKey, correlationId)
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
