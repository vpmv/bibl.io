package openlibrary

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type JobType int

const (
	JobTypeBook JobType = iota + 1
	JobTypeAuthor
	JobTypeBookSearch
	JobTypeAuthorSearch
	JobTypeISBN
	JobTypeWork
)

// JobResolver is responsible for parsing the Job result as desired
type JobResolver func(j Job)

// An implementable Job interface to overcome the shortcomings of generics
type Job interface {
	Hash() string
	Type() JobType
	Queue()
	StatusCode() int
	// Value returns the normalized job response
	Value() any
	uri() string
	queryParams() map[string]string
	setResult(code int, body io.ReadCloser)
	decode() (any, error)
}

type GenericJob[R JobResponse[T], T any] struct {
	hash       string
	jobType    JobType
	uriPath    string
	params     map[string]string
	res        io.ReadCloser
	statusCode int
	rawValue   R
	value      T
}

func (j *GenericJob[R, T]) Type() JobType {
	return j.jobType
}

func (j *GenericJob[R, T]) Queue() {
	jobQueue.queue(j)
}

func (j *GenericJob[R, T]) StatusCode() int {
	return j.statusCode
}

func (j *GenericJob[R, T]) Value() any {
	return j.value
}

func (j *GenericJob[R, T]) uri() string {
	return j.uriPath
}

func (j *GenericJob[R, T]) queryParams() map[string]string {
	return j.params
}

func (j *GenericJob[R, T]) setResult(code int, body io.ReadCloser) {
	j.statusCode = code
	j.res = body
}

func (j *GenericJob[R, T]) Hash() string {
	return j.hash
}

func (j *GenericJob[R, T]) decode() (any, error) {
	var t T

	err := json.NewDecoder(j.res).Decode(&j.rawValue)
	if err != nil {
		return t, fmt.Errorf(`error decoding %s response: %s`, reflect.TypeOf(j), err.Error())
	}
	j.value = j.rawValue.Normalize()

	return j.value, nil

}

// // NewJob creates a new Job
// // FIXME: does not support solr queries
func NewJob(jobType JobType, params map[string]string) (j Job) {
	md5Sum := md5.Sum([]byte(fmt.Sprintf("%s_%v", jobType, params)))
	hash := fmt.Sprintf("%x", md5Sum)

	switch jobType {
	case JobTypeBook:
		j = &GenericJob[BookResponse, *Book]{
			hash:     hash,
			jobType:  jobType,
			params:   nil,
			uriPath:  `/books/` + params[`key`],
			rawValue: BookResponse{},
		}

	case JobTypeBookSearch:
		j = &GenericJob[BookSearchResponse[*Book, []*Book, BookResponse], []*Book]{
			hash:    hash,
			jobType: jobType,
			params:  params,
			uriPath: `/search.json`,
			rawValue: BookSearchResponse[*Book, []*Book, BookResponse]{
				Docs: []BookResponse{},
			},
		}
	case JobTypeAuthor:
		j = &GenericJob[AuthorResponse, *Author]{
			hash:    hash,
			jobType: jobType,
			params:  nil,
			uriPath: `/authors/` + params[`key`],
		}
	case JobTypeAuthorSearch:
		j = &GenericJob[AuthorSearchResponse[*Author, []*Author, AuthorResponse], []*Author]{
			hash:    hash,
			jobType: jobType,
			params:  params,
			uriPath: `/authors/search.json`,
		}
	case JobTypeWork:
		j = &GenericJob[WorkResponse, *Work]{
			hash:    hash,
			jobType: jobType,
			params:  nil,
			uriPath: `/works/` + params[`key`],
		}

	case JobTypeISBN:
		j = &GenericJob[ISBNResponse, *ISBN]{
			hash:    hash,
			jobType: jobType,
			params:  nil,
			uriPath: `/isbn/` + params[`isbn`],
		}

	}

	return j
}

// Albeit now needing a type for every job, we can still easily use channels
//
//type abstractJob struct {
//	hash       string
//	jobType    JobType
//	params     map[string]string
//	res        io.ReadCloser
//	statusCode int
//	value      any
//}
//
//func (j abstractJob) Type() JobType {
//	return j.jobType
//}
//
//func (j *abstractJob) Hash() string {
//	return j.hash
//}
//
//func (j abstractJob) StatusCode() int {
//	return j.statusCode
//}
//
//func (j *abstractJob) setResult(code int, body io.ReadCloser) {
//	j.statusCode = code
//	j.res = body
//}
//
////func (j *abstractJob) resolveResponse(job Job) (any, error) {
////	responseObj := job.responseObject()
////	err := json.NewDecoder(j.res).Decode(responseObj)
////	if err != nil {
////		return nil, fmt.Errorf(`error decoding %s response: %s`, reflect.TypeOf(job), err.Error())
////	}
////	j.value = responseObj.Normalize()
////
////	return j.value, nil
////}
//
//func (j abstractJob) Value() any {
//	return j.value
//}
//
//type ISBNJob struct {
//	*abstractJob
//}
//
//func (j ISBNJob) uri() string {
//	return `/isbn/` + j.params[`isbn`]
//}
//
//func (j ISBNJob) queryParams() map[string]string {
//	return nil
//}
//
//func (j ISBNJob) decode() (any, error) {
//	return j.resolveResponse(j)
//}
//
//func (j ISBNJob) responseObject() JobResponse {
//	return new(ISBNResponse)
//}
//
//func (j ISBNJob) queue() {
//	jobQueue <- j
//}
//
//type BookJob struct {
//	*abstractJob
//	search bool
//}
//
//func (j BookJob) uri() string {
//	if j.search {
//		return `/search.json`
//	}
//	return `/books/` + j.params[`key`]
//}
//
//func (j BookJob) queryParams() map[string]string {
//	if !j.search {
//		return nil
//	}
//	return j.params
//}
//
//func (j BookJob) responseObject() JobResponse {
//	if j.search {
//		return new(SearchResponse[*BookResponse])
//	}
//	return new(BookResponse)
//}
//
//func (j BookJob) decode() (any, error) {
//	return j.resolveResponse(j)
//}
//
//type BookSearchJob struct {
//	BookJob
//}
//
//func (j BookSearchJob) decode() (any, error) {
//	return j.resolveResponse(j)
//}
//
//type WorkJob struct {
//	*abstractJob
//}
//
//func (j WorkJob) uri() string {
//	return `/works/` + j.params[`key`]
//}
//
//func (j WorkJob) queryParams() map[string]string {
//	return nil
//}
//
//func (j WorkJob) decode() (any, error) {
//	return j.resolveResponse(j)
//}
//
//func (j WorkJob) responseObject() JobResponse {
//	return new(WorkResponse)
//}
//
//type AuthorJob struct {
//	*abstractJob
//	search bool
//}
//
//type AuthorSearchJob struct {
//	AuthorJob
//	result *SearchResponse[*AuthorResponse]
//}
//
//func (j AuthorJob) uri() string {
//	if j.search {
//		return `/search/authors.json`
//	}
//	return `/authors/` + j.params[`key`]
//}
//
//func (j AuthorJob) queryParams() map[string]string {
//	if !j.search {
//		return nil
//	}
//	return j.params
//}
//
//func (j AuthorJob) decode() (any, error) {
//	return j.resolveResponse(j)
//}
//
//func (j AuthorJob) responseObject() JobResponse {
//	if j.search {
//		return new(SearchResponse[*AuthorResponse])
//	}
//	return new(AuthorResponse)
//}
