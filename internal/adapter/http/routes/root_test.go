package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/internal/adapter/http/routes"
	"github.com/LiquidCats/rater/internal/adapter/http/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler_Handle(t *testing.T) {
	l := zerolog.New(zerolog.NewTestWriter(t))

	router := server.NewRouter(&l)
	handler := routes.NewRootHandler()

	router.GET("/", handler.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"version":"v1"}`, w.Body.String())
}
