package api

import (
	"context"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	chiCORS "github.com/go-chi/cors"
	"github.com/go-fuego/fuego"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/env"
	"github.com/vpmv/bibl.io/pkg/service/openlibrary"
	"github.com/vpmv/bibl.io/pkg/storage"
	"net/http"
)

const (
	PermissionBooksRead   = `books.read`
	PermissionAuthorsRead = `authors.read`
)

var staleJobTimeout float64 = 86400 // 1 day

type API struct {
	auth        Authenticator
	log         *logrus.Logger
	store       storage.StorageClient
	openLibrary *openlibrary.Client
	Security    fuego.Security
}

func New(auth Authenticator, logger *logrus.Logger, store storage.StorageClient, openLibClient *openlibrary.Client) *API {
	return &API{
		auth:        auth,
		store:       store,
		log:         logger,
		openLibrary: openLibClient,
	}

}

// Bootstrap of the API
func (api *API) Bootstrap(fs *fuego.Server, ctx context.Context) {
	api.log.Debug(`Migrating DB...`)
	if err := api.store.Migrate(); err != nil {
		api.log.Fatal(`Error migrating models`)
	}

	api.log.Debug(`Setting up router...`)
	cors := chiCORS.New(chiCORS.Options{
		//AllowedOrigins:   strings.Split(os.Getenv(`ALLOWED_ORIGINS`), `;`),
		AllowedMethods:   []string{http.MethodGet},
		AllowedHeaders:   []string{"Cookie", "Accept", "Authorization", "Content-Type", "Origin"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	})
	fuego.Use(fs, cors.Handler)

	fuego.Use(fs, chiMiddleware.Throttle(env.GetInt(`THROTTLE_THRESHOLD`, 100)))

	// Health check
	fuego.Get(fs, "/", api.HealthCheck)

	v1 := fuego.Group(fs, `/v1`)
	v1.Header(`Authorization`, `Auth token`)
	fuego.Use(v1, api.bearerAuthorization)

	// /v1/books/...
	booksGroup := fuego.Group(v1, `/books`)
	fuego.Use(booksGroup, api.hasPermission(PermissionBooksRead))

	fuego.Get(booksGroup, `/`, api.ListBooks, optionPagination)
	fuego.Get(booksGroup, `/search`, api.SearchBooks, optionBooks)

	api.log.Debug(`Starting OpenLibrary worker...`)
	staleJobTimeout = env.GetFloat(`OPENLIBRARY_STALE_TIMEOUT`, staleJobTimeout)
	go api.openLibrary.StartWorker(ctx, api.OpenLibJobResolver)

	// /v1/authors/....
	authorsGroup := fuego.Group(v1, `/authors`)
	fuego.Use(authorsGroup, api.hasPermission(PermissionAuthorsRead))
	fuego.Get(authorsGroup, `/`, api.ListAuthors, optionPagination)
	fuego.Get(authorsGroup, `/search`, api.SearchAuthors, optionAuthors)
	api.log.Debug(`API bootstrapped!`)
}

// HealthCheck is used to do server health checks outside normal api routing
func (api *API) HealthCheck(c fuego.ContextNoBody) (string, error) {
	return `tabula rasa`, nil
}

func (api *API) QueryParamUint(c fuego.ContextNoBody, param string) uint64 {
	i := c.QueryParamInt(param)
	return uint64(i)
}
