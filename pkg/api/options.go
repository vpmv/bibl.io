package api

import (
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

const (
	ParamPage     = "page"
	ParamPageSize = "pageSize"
)

var (
	optionPage     = option.QueryInt(ParamPage, "Page number", param.Default(1))
	optionPageSize = option.QueryInt(ParamPageSize, "Page size", param.Default(100))

	optionPagination = option.Group(
		optionPage,
		optionPageSize,
	)
)
