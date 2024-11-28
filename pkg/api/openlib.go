package api

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/dto"
	"github.com/vpmv/bibl.io/pkg/service/openlibrary"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func (api API) isStale(hash string) bool {
	t := api.store.GetJobQueryTime(hash)
	now := time.Now()

	stale := t == nil || now.Sub(*t).Seconds() > staleJobTimeout
	api.log.WithFields(logrus.Fields{`stale`: stale, `hash`: hash}).Debug(`checked job timings`)
	return stale
}

func (api API) OpenLibJobResolver(job openlibrary.Job) {
	api.log.WithFields(logrus.Fields{`hash`: job.Hash()}).Debugf(`retrieving job result %s`, reflect.TypeOf(job))
	res := job.Value()

	switch job.Type() {
	case openlibrary.JobTypeBook:
		api.StoreBook(res.(*openlibrary.Book))
	case openlibrary.JobTypeBookSearch:
		for _, book := range res.([]*openlibrary.Book) {
			api.StoreBook(book)
		}
	case openlibrary.JobTypeAuthor:
		api.StoreAuthor(res.(*openlibrary.Author))
	case openlibrary.JobTypeAuthorSearch:
		for _, author := range res.([]*openlibrary.Author) {
			api.StoreAuthor(author)
		}
	case openlibrary.JobTypeWork:
		work := res.(*openlibrary.Work)
		book, err := api.store.GetBook(work.Key, false)
		if err != nil {
			api.log.Errorf(`error getting book for work %s: %s`, work.Key, err.Error())
			return
		}
		convertWork(book, work)
		if err := api.store.SaveBook(book); err != nil {
			api.log.Errorf("failed to save book: %v", err)
		}
	default:
		api.log.Error(`fu`)
	}

}

func (api API) StoreBook(book *openlibrary.Book) {
	api.log.Debugf(`processing book %s (%s)`, book.Key, book.Title)

	obj := convertBook(book)
	if err := api.store.SaveBook(obj); err != nil {
		api.log.Errorf("failed to save book: %v", err)
	}
	api.log.Debugf(`book %s saved`, book.Key)

	for _, authorKey := range book.Authors {
		if author, err := api.store.GetAuthor(authorKey, false); author == nil || author.Name == `` || errors.Is(err, gorm.ErrRecordNotFound) {
			api.log.Debugf(`retrieving author %s`, authorKey)
			job := api.openLibrary.GetAuthor(authorKey)
			if api.isStale(job.Hash()) {
				job.Queue()
			}
			api.store.AddQuery(job.Hash())
		}
	}
	api.openLibrary.GetWork(book.WorkKey)
}

func (api API) StoreAuthor(author *openlibrary.Author) {
	obj := convertAuthor(author)

	api.log.Debugf(`Saving author %s`, author.Key)
	if err := api.store.SaveAuthor(obj); err != nil {
		api.log.Errorf("failed to save author: %s", err.Error())
	}
	return
}

func convertBook(b *openlibrary.Book) *dto.Book {
	authors := make([]*dto.Author, len(b.Authors))
	for i, author := range b.Authors {
		authors[i] = &dto.Author{
			Key: author,
		}
	}

	return &dto.Book{
		Key:   b.Key,
		ISBN:  b.ISBNList,
		Title: b.Title,
		//Language: b.Languages[0], //fixme
		Authors: authors,
	}
}

// fixme
func convertWork(dto *dto.Book, work *openlibrary.Work) {
	dto.Summary = work.Description
}

func convertAuthor(a *openlibrary.Author) *dto.Author {
	return &dto.Author{
		Key:         a.Key,
		Name:        a.Name,
		DateOfBirth: a.BirthDate.Format(time.DateOnly),
		DateOfDeath: a.DeathDate.Format(time.DateOnly),
		Nationality: "",
	}
}
