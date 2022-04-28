# Correlation Id
Correlations Id is used in distributed applications to trace requests across multiple services. This package provides a lightweight correlation id middlware. Request headers are checked for a correlation id. If found or generated, this correlation id is attached to the request context which can be used to access the current correlation id where it is required for logging etc. Based on middleware settings correlation id can be returned with response headers.

## Install and update

`go get -u github.com/dmytrohridin/correlation-id`

## Get Correlation Id value
To get correlation id value use `FromContext` function.
```go
id := correlationid.FromContext(req.Context())
``` 

## HeaderName setting
By default `Correlation-Id` key used as header name. Behavior can be overridden.
```go
middleware := correlationid.New()
middleware.HeaderName = "Request-Id"
```

## IdGenerator setting
By default [google/uuid package](https://github.com/google/uuid) is used for correlation id when `Correlation-Id` header is not included in request headers. Behavior can be overridden by applying a custom function.
```go
middleware := correlationid.New()
middleware.IdGenerator = func() string {
    return "any_id"
}
```

## IncludeInResponse setting
By default `Correlation-Id` is sent in response header. Behavior can be overridden.
```go
middleware := correlationid.New()
middleware.IncludeInResponse = false
```

## EnforceHeader setting
By default `Correlation-Id` is not required in request headers. Behavior can be overridden.
```go
middleware := correlationid.New()
middleware.EnforceHeader = true
```
If header is enforced and client does not send it - `BadRequest` will be returned.
