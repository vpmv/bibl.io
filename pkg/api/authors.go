package api

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/dto"
)

const (
	ParamAuthorName = `name`
)

var (
	optionAuthors = option.Group(
		option.Query(ParamAuthorName, "Name"),
	)
)

func (api *API) ListAuthors(c fuego.ContextNoBody) (authors []*dto.Author, err error) {
	page := c.QueryParamInt(ParamPage)
	pageSize := c.QueryParamInt(ParamPageSize)

	authors, err = api.store.GetAuthors(page, pageSize)
	return
}

func (api *API) SearchAuthors(c fuego.ContextNoBody) (authors []*dto.Author, err error) {
	name := c.QueryParam(ParamAuthorName)

	api.log.WithFields(logrus.Fields{`name`: name}).Debug(`received authors request`)

	params := &dto.Author{
		Name: name,
	}

	api.log.Debug(`querying...`)
	authors, err = api.store.SearchAuthors(params)
	api.log.Debug(`got query result`)

	go api.queueAuthorSearch(name)

	return authors, err
}

func (api *API) queueAuthorSearch(name string) {
	job := api.openLibrary.SearchAuthors(name)
	if api.isStale(job.Hash()) {
		job.Queue()
	}
	api.store.AddQuery(job.Hash())
}
