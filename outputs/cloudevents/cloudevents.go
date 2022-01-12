package cloudevents

import (
	"context"
	"errors"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/nsmith5/rekor-sidekick/outputs"
)

const (
	driverName  = `cloudevents`
	eventSource = `github.com/nsmith5/rekor-sidekick`
	eventType   = `rekor-sidekick/event`
)

var (
	ErrConfigMissingURL = errors.New(`cloudevents: driver requires "url" in configuration`)
	ErrURLWrongType     = errors.New(`cloudevents: "url" configuration must be a string`)
)

type driver struct {
	url    string
	client client.Client
}

func (d *driver) Send(e outputs.Event) error {
	event := ce.NewEvent()
	event.SetSource(eventSource)
	event.SetType(eventType)
	event.SetData(ce.ApplicationJSON, e)

	ctx := ce.ContextWithTarget(context.Background(), d.url)

	if result := d.client.Send(ctx, event); !ce.IsACK(result) {
		return result
	}

	return nil
}

func (d *driver) Name() string {
	return driverName
}

func createDriver(conf map[string]interface{}) (outputs.Output, error) {
	opaque, ok := conf[`url`]
	if !ok {
		return nil, ErrConfigMissingURL
	}
	url, ok := opaque.(string)
	if !ok {
		return nil, ErrURLWrongType
	}

	client, err := ce.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	return &driver{url, client}, nil
}

func init() {
	outputs.RegisterDriver(driverName, outputs.CreatorFunc(createDriver))
}
