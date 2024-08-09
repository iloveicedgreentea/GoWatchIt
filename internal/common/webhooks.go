package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/iloveicedgreentea/go-plex/models"
)

func DecodeWebhook(payload []string) (models.PlexWebhookPayload, int, error) {
	var pwhPayload models.PlexWebhookPayload

	err := json.Unmarshal([]byte(payload[0]), &pwhPayload)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// unmarshall error
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request has an invalid value in %q field at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return pwhPayload, http.StatusBadRequest, errors.New(msg)

		default:
			return pwhPayload, http.StatusInternalServerError, err
		}
	}

	return pwhPayload, 0, nil
}


// only used for tests
// func decodeJfWebhook(data []byte) (out models.JellyfinWebhook, err error) {
// 	err = json.Unmarshal(data, &out)
// 	if err != nil {
// 		log.Errorf("Error decoding payload: %v", err)
// 	}

// 	return out, err
// }