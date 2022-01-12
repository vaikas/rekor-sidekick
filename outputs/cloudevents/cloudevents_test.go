package cloudevents

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nsmith5/rekor-sidekick/outputs"
)

func TestCreateDriver(t *testing.T) {
	tests := map[string]struct {
		Conf  map[string]interface{}
		Error error
	}{
		"valid configuration": {
			Conf: map[string]interface{}{
				`url`: `https://localhost:8080`,
			},
			Error: nil,
		},
		"wrong type for url": {
			Conf: map[string]interface{}{
				`url`: 22,
			},
			Error: ErrURLWrongType,
		},
		"missing url in config": {
			Conf:  map[string]interface{}{},
			Error: ErrConfigMissingURL,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := createDriver(data.Conf)
			if err != data.Error {
				t.Errorf("Expected err %q, but recieved %q", data.Error, err)
			}
		})
	}
}

func TestSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Ce-Source") != eventSource {
			t.Errorf("expected event source to be %s but got %s", eventSource, r.Header.Get("Ce-Source"))
		}
		if r.Header.Get("Ce-Type") != eventType {
			t.Errorf("expected event type to be %s but got %s", eventSource, r.Header.Get("Ce-Type"))
		}
	}))

	conf := map[string]interface{}{
		`url`: ts.URL,
	}

	driver, err := createDriver(conf)
	if err != nil {
		t.Fatal(err)
	}

	err = driver.Send(outputs.Event{})
	if err != nil {
		t.Error(err)
	}
}
