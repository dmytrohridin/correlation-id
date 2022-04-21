module github.com/dmytrohridin/correlation-id/examples

go 1.18

replace github.com/dmytrohridin/correlation-id => ../

require (
	github.com/dmytrohridin/correlation-id v1.0.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/gorilla/mux v1.8.0
)

require github.com/google/uuid v1.3.0 // indirect
