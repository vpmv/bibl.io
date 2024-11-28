package api

import (
	"fmt"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/dto"
)

const (
	ParamBookTitle = `title`
	ParamBookISBN  = `isbn`
)

var (
	optionBooks = option.Group(
		option.Query(ParamBookTitle, "Title"),
		option.QueryInt(ParamBookISBN, "ISBN"),
	)
)

func (api *API) ListBooks(c fuego.ContextNoBody) (books []*dto.Book, err error) {
	page := c.QueryParamInt(ParamPage)
	pageSize := c.QueryParamInt(ParamPageSize)

	books, err = api.store.GetBooks(page, pageSize)
	return
}

func (api *API) SearchBooks(c fuego.ContextNoBody) (books []*dto.Book, err error) {
	title := c.QueryParam(ParamBookTitle)
	isbn := api.QueryParamUint(c, ParamBookISBN)
	if isbn < 13 && isbn < 10 {
		isbn = 0
	}
	api.log.WithFields(logrus.Fields{`title`: title, `isbn`: isbn}).Debug(`received request`)

	params := &dto.Book{
		Title: title,
	}
	if isbn > 0 {
		params.ISBN = []string{fmt.Sprint(isbn)} // fixme
	}

	api.log.Debug(`querying...`)
	books, err = api.store.SearchBooks(params)
	api.log.Debug(`got query result`)

	go api.queueBookSearch(title, isbn)

	return books, err
}

func (api API) queueBookSearch(title string, isbn uint64) {
	if isbn > 0 {
		job := api.openLibrary.GetBookByISBN(isbn)
		if api.isStale(job.Hash()) {
			job.Queue()
		}
	}

	job := api.openLibrary.SearchBooks(title, ``)
	if api.isStale(job.Hash()) {
		job.Queue()
	}

	api.store.AddQuery(job.Hash())
}
