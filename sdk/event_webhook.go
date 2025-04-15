package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// EventWebhook is a Sendgrid event webhook settings.
type EventWebhook struct { //nolint:maligned
	Id                string `json:"id,omitempty"`
	Enabled           bool   `json:"enabled"`
	URL               string `json:"url,omitempty"`
	GroupResubscribe  bool   `json:"group_resubscribe"` //nolint:tagliatelle
	Delivered         bool   `json:"delivered"`
	GroupUnsubscribe  bool   `json:"group_unsubscribe"` //nolint:tagliatelle
	SpamReport        bool   `json:"spam_report"`       //nolint:tagliatelle
	Bounce            bool   `json:"bounce"`
	Deferred          bool   `json:"deferred"`
	Unsubscribe       bool   `json:"unsubscribe"`
	Processed         bool   `json:"processed"`
	Open              bool   `json:"open"`
	Click             bool   `json:"click"`
	Dropped           bool   `json:"dropped"`
	OAuthClientID     string `json:"oauth_client_id,omitempty"`     //nolint:tagliatelle
	OAuthClientSecret string `json:"oauth_client_secret,omitempty"` //nolint:tagliatelle
	OAuthTokenURL     string `json:"oauth_token_url,omitempty"`     //nolint:tagliatelle
}

type EventWebhookSigning struct {
	Enabled   bool   `json:"enabled"`
	PublicKey string `json:"public_key"` //nolint:tagliatelle
}

func parseEventWebhook(respBody string) (*EventWebhook, RequestError) {
	var body EventWebhook
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseEventWebhookSigning(respBody string) (*EventWebhookSigning, RequestError) {
	var body EventWebhookSigning
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateEventWebhook creates an EventWebhook and returns it.
func (c *Client) CreateEventWebhook(webhook *EventWebhook) (*EventWebhook, RequestError) {
	if webhook.URL == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrURLRequired,
		}
	}

	respBody, statusCode, err := c.Post("POST", "/user/webhooks/event/settings", webhook)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed posting event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedPatchingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhook(respBody)
}

// PatchEventWebhook creates an EventWebhook and returns it.
func (c *Client) PatchEventWebhook(webhook *EventWebhook) (*EventWebhook, RequestError) {
	if webhook.URL == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrURLRequired,
		}
	}

	if webhook.Id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrIdRequired,
		}
	}

	respBody, statusCode, err := c.Post("PATCH", fmt.Sprintf("/user/webhooks/event/settings/%s", webhook.Id), webhook)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed patching event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedPatchingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhook(respBody)
}

// ReadEventWebhook retrieves an EventWebhook and returns it.
func (c *Client) ReadEventWebhook(id string) (*EventWebhook, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrIdRequired,
		}
	}

	respBody, _, err := c.Get("GET", fmt.Sprintf("/user/webhooks/event/settings/%s", id))
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseEventWebhook(respBody)
}

func (c *Client) ConfigureEventWebhookSigning(id string, enabled bool) (*EventWebhookSigning, RequestError) {
	respBody, statusCode, err := c.Post("PATCH", fmt.Sprintf("/user/webhooks/event/settings/signed/%s", id), EventWebhookSigning{
		Enabled: enabled,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed patching event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedPatchingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhookSigning(respBody)
}

func (c *Client) ReadEventWebhookSigning(id string) (*EventWebhookSigning, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrIdRequired,
		}
	}

	respBody, _, err := c.Get("GET", fmt.Sprintf("/user/webhooks/event/settings/signed/%s", id))
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseEventWebhookSigning(respBody)
}
