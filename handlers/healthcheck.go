package handlers

import (
	"net/http"
)

func ProcessHealthcheckWebhook() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}

	return http.HandlerFunc(fn)
}