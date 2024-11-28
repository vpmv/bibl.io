package openlibrary

import (
	"strconv"
	"strings"
	"time"
)

type KeyRelation struct {
	Key string `json:"key"`
}

type JobResponse[T any] interface {
	Normalize() T
}

// BookSearchResponse holds result set of Books
// where Type T is *Book, Exported X is []*Book and Response R is JobResponse[T: *Book]
type BookSearchResponse[T *Book, X []T, R JobResponse[T]] struct {
	Docs []R `json:"docs"`
}

func (sr BookSearchResponse[T, X, R]) Normalize() []*Book {
	var res []*Book
	for _, d := range sr.Docs {
		res = append(res, d.Normalize())
	}
	return res
}

type AuthorSearchResponse[T *Author, X []T, R JobResponse[T]] struct {
	Docs []R `json:"docs"`
}

func (sr AuthorSearchResponse[T, X, R]) Normalize() []*Author {
	var res []*Author
	for _, d := range sr.Docs {
		res = append(res, d.Normalize())
	}
	return res
}

type AuthorResponse struct {
	Key            string   `json:"key"`
	BirthDate      string   `json:"birth_date"`
	DeathDate      string   `json:"death_date"`
	Name           string   `json:"name"`
	Biography      string   `json:"bio"`
	TopSubjects    []string `json:"top_subjects"`
	TopWork        string   `json:"top_work"`
	Type           string   `json:"type"`
	WorkCount      int      `json:"work_count"`
	RatingsAverage float64  `json:"ratings_average"`
	RatingsCount   int      `json:"ratings_count"`
}

type Author struct {
	Key            string    `json:"key"`
	BirthDate      time.Time `json:"birth_date"`
	DeathDate      time.Time `json:"death_date"`
	Name           string    `json:"name"`
	TopSubjects    []string  `json:"top_subjects"`
	TopWork        string    `json:"top_work"`
	WorkCount      int       `json:"work_count"`
	RatingsAverage float64   `json:"ratings_average"`
	RatingsCount   int       `json:"ratings_count"`
}

func (ar AuthorResponse) Normalize() *Author {
	a := &Author{
		Key:            strings.Replace(ar.Key, `/authors/`, ``, -1),
		BirthDate:      findDate(ar.BirthDate),
		DeathDate:      findDate(ar.DeathDate),
		Name:           ar.Name,
		TopSubjects:    ar.TopSubjects,
		TopWork:        ar.TopWork,
		WorkCount:      ar.WorkCount,
		RatingsAverage: ar.RatingsAverage,
		RatingsCount:   ar.RatingsCount,
	}

	return a
}

type BookResponse struct {
	Key           string   `json:"key"`
	Title         string   `json:"title"`
	AuthorKeys    []string `json:"author_key"`
	PublishDates  []string `json:"publish_date"`
	Publishers    []string `json:"publishers"`
	Covers        []int    `json:"covers"`
	Contributors  []string `json:"contributors"`
	Languages     []string `json:"languages"`
	FirstSentence []string `json:"first_sentence"`
	NumberOfPages int      `json:"number_of_pages"`
	//Classifications struct{} `json:"classifications"`
	OCA_ID         string   `json:"ocaid"`
	ISBNs          []string `json:"isbn"`
	LatestRevision int      `json:"latest_revision"`
	Revision       int      `json:"revision"`
}

type Book struct {
	Key            string    `json:"key"`
	WorkKey        string    `json:"work_key"`
	Title          string    `json:"title"`
	Authors        []string  `json:"authors"`
	PublishDate    time.Time `json:"publish_date"`
	Publishers     []string  `json:"publishers"`
	Contributors   []string  `json:"contributors"`
	Languages      []string  `json:"language"`
	FirstSentence  string    `json:"first_sentence"`
	NumberOfPages  int       `json:"number_of_pages"`
	OCA_ID         string    `json:"ocaid"`
	ISBNList       []string  `json:"isbn_list"`
	LatestRevision int       `json:"latest_revision"`
	Revision       int       `json:"revision"`
}

func (br BookResponse) Normalize() *Book {
	var pubDate time.Time
	for _, d := range br.PublishDates {
		pubDate = findDate(d)
		if !pubDate.IsZero() {
			break
		}
	}

	return &Book{
		Key:          strings.Replace(strings.Replace(br.Key, `/works/`, ``, -1), `/books/`, ``, -1), // fixme: search response returns /works/XYZ
		WorkKey:      strings.Replace(br.Key, `/works/`, ``, -1),
		Title:        br.Title,
		Authors:      br.AuthorKeys,
		PublishDate:  pubDate,
		Publishers:   br.Publishers,
		Contributors: br.Contributors,
		Languages:    br.Languages,
		//FirstSentence:  br.FirstSentence[0], // fixme: api doesn't specify language
		NumberOfPages:  br.NumberOfPages,
		ISBNList:       br.ISBNs,
		OCA_ID:         br.OCA_ID,
		LatestRevision: br.LatestRevision,
		Revision:       br.Revision,
	}
}

type ISBNResponse struct {
	Authors       []KeyRelation `json:"authors"`
	BookKey       string        `json:"key"`
	Work          []KeyRelation `json:"works"`
	ISBN_10       string        `json:"isbn_10"`
	ISBN_13       string        `json:"isbn_13"`
	FirstSentence struct {
		Value string `json:"value"`
	} `json:"first_sentence"`
}

type ISBN struct {
	Key       string `json:"key"`
	AuthorKey string `json:"author_key"`
	WorkKey   string `json:"work_key"`
	ISBN_10   uint64 `json:"isbn_10"`
	ISBN_13   uint64 `json:"isbn_13"`
}

func (ir ISBNResponse) Normalize() *ISBN {
	isbn10, _ := strconv.ParseUint(ir.ISBN_10, 10, 64)
	isbn13, _ := strconv.ParseUint(ir.ISBN_13, 10, 64)

	return &ISBN{
		Key:       strings.Replace(ir.BookKey, `/books/`, ``, -1),
		AuthorKey: strings.Replace(ir.Authors[0].Key, `/authors/`, ``, -1),
		WorkKey:   strings.Replace(ir.Work[0].Key, `/works/`, ``, -1),
		ISBN_10:   isbn10,
		ISBN_13:   isbn13,
	}
}

type WorkResponse struct {
	Key           string   `json:"key"`
	Description   string   `json:"description"`
	Subjects      []string `json:"subjects"`
	SubjectPlaces []string `json:"subject_places"`
	Covers        []int    `json:"covers"`
}

type Work struct {
	WorkResponse
}

func (wr WorkResponse) Normalize() *Work {
	return &Work{wr}
}

// findDate attempts to resolve one of the many date stamps used by openlibrary
func findDate(i string) time.Time {
	formats := []string{
		time.DateOnly,
		`2 January 2006`,
		`2 January, 2006`,
		`January 2 2006`,
		`January 2, 2006`,
		`2006 January 2`,
		`2006, January 2`,
		`2006 January`,
		`2006`,
	}

	for _, format := range formats {
		t, e := time.Parse(format, i)
		if e == nil {
			return t
		}
	}

	return time.Time{}
}
