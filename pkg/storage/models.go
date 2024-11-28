package storage

import (
	"github.com/vpmv/bibl.io/pkg/dto"
	"gorm.io/gorm"
	"time"
)

type Author struct {
	gorm.Model
	Key         string `gorm:"primaryKey"`
	Name        string
	Description string
	DateOfBirth time.Time
	DateOfDeath time.Time
	Nationality string

	Books []Book `gorm:"many2many:book_authors;"`
}

func (a *Author) DTO(includeBooks bool) *dto.Author {
	death := ``
	if !a.DateOfDeath.IsZero() {
		death = a.DateOfDeath.Format(time.DateOnly)
	}

	var books []*dto.Book
	if includeBooks {
		books = make([]*dto.Book, len(a.Books))
		for i, b := range a.Books {
			books[i] = b.DTO(false)
		}
	}

	return &dto.Author{
		Key:         a.Key,
		Name:        a.Name,
		DateOfBirth: a.DateOfBirth.Format(time.DateOnly),
		DateOfDeath: death,
		Nationality: a.Nationality,
		Books:       books,
	}
}

func makeAuthor(dto *dto.Author) *Author {
	dob, _ := time.Parse(time.DateOnly, dto.DateOfBirth)
	dod, _ := time.Parse(time.DateOnly, dto.DateOfDeath)

	return &Author{
		Key:         dto.Key,
		Name:        dto.Name,
		DateOfBirth: dob,
		DateOfDeath: dod,
		Nationality: dto.Nationality,
		Description: dto.Description,
	}
}

type Book struct {
	gorm.Model
	Key string `gorm:"primaryKey"`
	//ISBN     []string
	Title    string
	Summary  string
	Language string `gorm:"type:varchar(3);"`

	Authors []Author `gorm:"many2many:book_authors;"`
}

func (b Book) DTO(includeAuthors bool) *dto.Book {
	var authors []*dto.Author
	if includeAuthors {
		authors = make([]*dto.Author, len(b.Authors))
		for i, a := range b.Authors {
			authors[i] = a.DTO(false)
		}
	}

	return &dto.Book{
		Key:      b.Key,
		Title:    b.Title,
		Summary:  b.Summary,
		Language: b.Language,
		//ISBN:     b.ISBN,
		Authors: authors,
	}
}

type QueryRegistry struct {
	gorm.Model
	Hash string `gorm:"primaryKey"`
}

func makeBook(dto *dto.Book) *Book {
	authors := make([]Author, len(dto.Authors))
	for i, a := range dto.Authors {
		authors[i] = *makeAuthor(a)
	}
	return &Book{
		Key: dto.Key,
		//ISBN:     dto.ISBN,
		Title:    dto.Title,
		Summary:  dto.Summary,
		Language: dto.Language,
		Authors:  authors,
	}
}
