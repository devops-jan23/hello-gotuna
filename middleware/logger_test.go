package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestLogging(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	middleware := middleware.Logger(logger)
	handler := middleware(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Contains(t, wlog.String(), "GET")
	assert.Contains(t, wlog.String(), "/sample")
}

func TestRecoveringFromPanic(t *testing.T) {

	needle := "assignment to entry in nil map"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this code will panic with: assignment to entry in nil map
	})

	// basic whoops html template
	//	whoopsTmpl := templating.GetNativeTemplatingEngine(i18n.NewTranslator(nil)).
	//		Mount(
	//			doubles.NewFileSystemStub(
	//				map[string][]byte{
	//					"app.html":   []byte(`{{define "app"}}{{block "sub" .}}{{end}}{{end}}`),
	//					"error.html": []byte(`{{define "sub"}}{{.Data.error}}<hr>{{.Data.stacktrace}}{{end}}`),
	//				}))

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	middleware := middleware.Logger(logger)
	handler := middleware(badHandler)

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), needle)
	assert.Contains(t, wlog.String(), needle)
}
