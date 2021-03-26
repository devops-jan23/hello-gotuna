package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/templating"
	"github.com/gorilla/mux"
)

func Logger(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())
					// log error and stack trace to console
					logger.Printf("PANIC RECOVERED: %v", err)
					logger.Println(stacktrace)

					//fmt.Println(err, stacktrace)

					w.WriteHeader(http.StatusInternalServerError)

					tmpl := templating.GetNativeTemplatingEngine(i18n.NewTranslator(i18n.En)) // TODO
					tmpl.
						Set("error", err).
						Set("stacktrace", stacktrace).
						Render(w, r, "app.html", "error.html")
					return
				}
			}()

			next.ServeHTTP(w, r)

			logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func whoops(err interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusInternalServerError)

		tmpl := templating.GetNativeTemplatingEngine(i18n.NewTranslator(i18n.En)) // TODO
		tmpl.
			Set("error", err).
			Set("stacktrace", string(debug.Stack())).
			Render(w, r, "app.html", "error.html")

	})
}
