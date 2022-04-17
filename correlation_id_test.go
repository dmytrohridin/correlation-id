package correlationid

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDefaultFlow(t *testing.T) {
	correlationid := New()
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", DefaultHeaderName)
	}

	id, err := uuid.Parse(header)
	if err != nil {
		t.Errorf("%s can not be parsed with default uuid provider", DefaultHeaderName)
	}

	if id == uuid.Nil {
		t.Errorf("%s is empty. Value: %s", DefaultHeaderName, id.String())
	}
}

func TestFlowWhenHeaderNameSetupedEmpty(t *testing.T) {
	correlationid := New()
	correlationid.HeaderName = ""
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", DefaultHeaderName)
	}

	id, err := uuid.Parse(header)
	if err != nil {
		t.Errorf("%s can not be parsed with default uuid provider", DefaultHeaderName)
	}

	if id == uuid.Nil {
		t.Errorf("%s is empty. Value: %s", DefaultHeaderName, id.String())
	}
}

func TestFlowWithCustomHeaderName(t *testing.T) {
	correlationid := New()
	correlationid.HeaderName = "X-Correlation-Id"
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(correlationid.HeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", correlationid.HeaderName)
	}

	id, err := uuid.Parse(header)
	if err != nil {
		t.Errorf("%s can not be parsed with default uuid provider", correlationid.HeaderName)
	}

	if id == uuid.Nil {
		t.Errorf("%s is empty. Value: %s", correlationid.HeaderName, id.String())
	}
}

func TestDefaultWithCustomIdGenerator(t *testing.T) {
	testId := "custom_id"
	correlationid := New()
	correlationid.IdGenerator = func() string {
		return testId
	}
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", DefaultHeaderName)
	}

	if header != testId {
		t.Errorf("Unexpected %s value. Expected: %s. Actual %s", DefaultHeaderName, testId, header)
	}
}

func TestFlowWhenNoIncludeInResponse(t *testing.T) {
	correlationid := New()
	correlationid.IncludeInResponse = false
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header != "" {
		t.Errorf("%s header value should be empty. Actual: %s", DefaultHeaderName, header)
	}
}

func TestFlowEnforceHeaderReturnBadRequest(t *testing.T) {
	correlationid := New()
	correlationid.EnforceHeader = true
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 400 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header != "" {
		t.Errorf("%s header should be empty", DefaultHeaderName)
	}
}

func TestFlowEnforceHeaderReturnOk(t *testing.T) {
	testId := "custom_id"
	correlationid := New()
	correlationid.EnforceHeader = true
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	req.Header.Add(DefaultHeaderName, testId)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", DefaultHeaderName)
	}

	if header != testId {
		t.Errorf("Unexpected %s value. Expected: %s. Actual %s", DefaultHeaderName, testId, header)
	}
}

func TestFlowWhenHeaderProvided(t *testing.T) {
	testId := "custom_id"
	correlationid := New()
	ts := httptest.NewServer(correlationid.Handle(getTestHandler()))
	defer ts.Close()

	u := fmt.Sprintf("%s/test", ts.URL)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", u, http.NoBody)
	req.Header.Add(DefaultHeaderName, testId)
	must(err, t)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected response StatusCode. Expected: 200. Actual: %d", resp.StatusCode)
	}
	must(err, t)
	defer resp.Body.Close()

	header := resp.Header.Get(DefaultHeaderName)
	if header == "" {
		t.Errorf("%s header value is empty", DefaultHeaderName)
	}

	if header != testId {
		t.Errorf("Unexpected %s value. Expected: %s. Actual %s", DefaultHeaderName, testId, header)
	}
}

func TestFromContextReturnValue(t *testing.T) {
	testVal := "test_value"
	ctx := context.WithValue(context.Background(), ContextKey, testVal)
	res := FromContext(ctx)
	if res != testVal {
		t.Errorf("Unexpected result. Expected: %s. Actual: %s", testVal, res)
	}
}

func TestFromContextReturnEmpty(t *testing.T) {
	res := FromContext(context.Background())
	if res != "" {
		t.Errorf("Unexpected result. Should be empty. Actual: %s", res)
	}
}

func getTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("test"))
	}
	return fn
}

func must(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
