package correlationid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type CorrelationIdType string

const (
	DefaultCorrelationIdHeaderName                   = "Correlation-Id"
	CorrelationId                  CorrelationIdType = "CorrelationId"
)

type CorrelationIdMiddleware struct {
	HeaderName        string
	IncludeInResponse bool
	EnforceHeader     bool
	IdGenerator       func() string
}

func New() CorrelationIdMiddleware {
	return CorrelationIdMiddleware{
		HeaderName:        DefaultCorrelationIdHeaderName,
		IncludeInResponse: true,
		EnforceHeader:     false,
		IdGenerator:       uuid.NewString,
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

			corrId = m.generateId()
		}

		if m.IncludeInResponse {
			rw.Header().Set(headerName, corrId)
		}

		updCtx := context.WithValue(r.Context(), CorrelationId, corrId)
		next.ServeHTTP(rw, r.WithContext(updCtx))
	})
}

func FromContext(ctx context.Context) string {
	corrId := ctx.Value(CorrelationId).(string)
	return corrId
}

func (m *CorrelationIdMiddleware) getHeaderName() string {
	if m.HeaderName == "" {
		return DefaultCorrelationIdHeaderName
	}

	return m.HeaderName
}

func (m *CorrelationIdMiddleware) generateId() string {
	if m.IdGenerator != nil {
		return m.IdGenerator()
	}

	return uuid.NewString()
}
