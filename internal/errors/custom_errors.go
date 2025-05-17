package errors

import "errors"

var (
	ErrNoUnsentMessages       = errors.New("no unsent messages found")
	ErrWebhookFailed          = errors.New("webhook returned unexpected status code")
	ErrInvalidWebhookResponse = errors.New("failed to decode webhook response")
	ErrMarshalPayload         = errors.New("failed to marshal message payload")
	ErrCreateHTTPRequest      = errors.New("failed to create HTTP request")
	ErrHTTPCall               = errors.New("HTTP request to webhook failed")
	ErrFetchMessages          = errors.New("could not fetch unsent messages from repository")
	ErrMarkMessageSent        = errors.New("could not mark message as sent in repository")
)
