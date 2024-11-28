package main

import (
	"context"
	"flag"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/api"
	"github.com/vpmv/bibl.io/pkg/dto"
	"github.com/vpmv/bibl.io/pkg/env"
	"github.com/vpmv/bibl.io/pkg/server"
	"github.com/vpmv/bibl.io/pkg/service/openlibrary"
	"github.com/vpmv/bibl.io/pkg/storage"
	"os"
)

type Config struct {
	Env      string
	LogLevel string
	DB       *storage.Config
}

type SimpleAuthenticator struct {
	tokens map[string]*dto.Authorization
}

func (auth *SimpleAuthenticator) AuthenticateBearer(apiKey string) (*dto.Authorization, bool, error) {
	// todo: remove this block
	if env.IsEnv(`development`) {
		return auth.tokens[`secret`], true, nil
	}

	if app, ok := auth.tokens[apiKey]; ok && app.Token == apiKey {
		return app, true, nil
	}

	return nil, false, nil
}

func main() {
	baseDir := flag.String(`basedir`, `config/`, `Base dir for configurations`)
	flag.Parse()

	env.LoadEnvironment(*baseDir)

	loglevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(`invalid log level: ` + err.Error())
	}

	logger := logrus.New()
	logger.SetLevel(loglevel)

	// Create the authenticator
	auth := &SimpleAuthenticator{
		tokens: map[string]*dto.Authorization{
			"secret": {
				Token:       "secret",
				Description: "Admin",
				Permissions: []dto.Permission{
					{api.PermissionBooksRead, ""},
					{api.PermissionAuthorsRead, ""},
				},
			},
		},
	}

	store, err := storage.NewMysqlClient(&storage.Config{
		Host:     os.Getenv(`DB_HOST`),
		Port:     os.Getenv(`DB_PORT`),
		User:     os.Getenv(`DB_USER`),
		Password: os.Getenv(`DB_PASSWORD`),
		DB:       os.Getenv(`DB_NAME`),
	}, logger)
	if err != nil {
		logger.Fatal(`failed to connect to datastore`, err)
	}

	openlib := openlibrary.NewClient(
		logger,
		env.GetString(`OPENLIBRARY_HOST`, `https://openlibrary.org`),
		env.GetInt(`OPENLIBRARY_RATE_LIMIT`, 20),
	)

	a := api.New(auth, logger, store, openlib)
	srv := server.New(context.Background(), a, os.Getenv(`API_HOST`))

	logger.Println("starting http listener on", aurora.Cyan(os.Getenv(`API_HOST`)))
	logger.Printf("allowed origins %s", aurora.Yellow(os.Getenv(`ALLOWED_ORIGINS)`)))
	logger.Printf("max concurrency per endpoint %d", aurora.Yellow(os.Getenv(`MAX_CONCURRENT`)))

	logger.Fatal(srv.Run())
}
