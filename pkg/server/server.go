package server

import (
	"context"
	"fmt"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	"github.com/vpmv/bibl.io/pkg/api"
	"net/http"
	"os"
	"runtime/debug"
)

// cache middleware for static files
func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600")
		h.ServeHTTP(w, r)
	})
}

func Recover() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stackTrace := debug.Stack()

					logEntry := chiMiddleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rec, stackTrace)
					} else {
						_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rec)
						debug.PrintStack()
					}

					errorText := http.StatusText(http.StatusInternalServerError)

					if os.Getenv(`ENV`) == `development` {
						errorText = fmt.Sprintf("%s\n %+v\n", errorText, rec)
					}

					http.Error(w, errorText, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func New(ctx context.Context, api *api.API, hostAddr string, options ...func(*fuego.Server)) *fuego.Server {
	serverOptions := []func(*fuego.Server){
		fuego.WithAddr(hostAddr),
		//fuego.WithAutoAuth(api.Oauth), // disabled in lieu of bearer tokens
		fuego.WithGlobalResponseTypes(http.StatusForbidden, "Forbidden"),
	}
	options = append(serverOptions, options...)

	app := fuego.NewServer(options...)
	app.OpenApiSpec.Info.Title = "Bibl.io API"
	api.Security = app.Security

	fuego.Use(app, Recover())
	fuego.Use(app, chiMiddleware.Compress(5, "application/json"))
	//fuego.Use(app, cache) // no need for this now

	api.Bootstrap(fuego.Group(app, "/api"), ctx)

	return app
}
