package handlers

import (
	"context"
	"github.com/bingoohuang/gg/pkg/jsoni"
	"log"
	"net/http"
)

func renderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := jsoni.NewEncoder(w).Encode(context.Background(), data); err != nil {
		log.Printf("%v", err)
	}
}
