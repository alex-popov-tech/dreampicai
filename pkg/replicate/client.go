package replicate

import (
	"net/http"

	r "github.com/replicate/replicate-go"
)

var Client *r.Client

type (
	Webhook              = r.Webhook
	PredictionInput      = r.PredictionInput
	WebhookEventType     = r.WebhookEventType
	WebhookSigningSecret = r.WebhookSigningSecret
)

func InitClient(token string) (*r.Client, error) {
	c, err := r.NewClient(r.WithToken(token))
	// init global replicate client for later usage across the app
	// if err != nil app won't start, so its safe to use global here
	Client = c
	return c, err
}

func ValidateWebhookRequest(req *http.Request, secret string) (bool, error) {
	return r.ValidateWebhookRequest(req, r.WebhookSigningSecret{Key: secret})
}
