package storage

import (
	"fmt"
	"github.com/vpmv/bibl.io/pkg/dto"
	"time"
)

type StorageClient interface {
	Migrate() error
	GetAPIToken(id string, token string) (*dto.Authorization, error)

	GetAuthor(key string, includeBooks bool) (*dto.Author, error)
	GetAuthors(page, size int) ([]*dto.Author, error)
	SearchAuthors(params *dto.Author) ([]*dto.Author, error)
	SaveAuthor(author *dto.Author) error

	GetBook(key string, includeAuthors bool) (*dto.Book, error)
	GetBooks(page, size int) ([]*dto.Book, error)
	SearchBooks(params *dto.Book) ([]*dto.Book, error)
	SaveBook(book *dto.Book) error

	AddQuery(hash string)
	GetJobQueryTime(hash string) *time.Time
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func (cfg Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB,
	)
}
