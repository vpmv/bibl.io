package storage

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vpmv/bibl.io/pkg/dto"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type MysqlClient struct {
	db  *gorm.DB
	log *logrus.Logger
}

func (m MysqlClient) Migrate() error {
	return m.db.Migrator().AutoMigrate(&Author{}, &Book{}, &QueryRegistry{})
}

// fixme
func (m MysqlClient) GetAPIToken(extID string, token string) (*dto.Authorization, error) {
	apiToken := &APIToken{Token: token, ExternalID: extID}
	err := m.db.Model(&APIToken{}).Where(apiToken).First(&apiToken).Error
	if err != nil {
		return nil, err
	}

	return apiToken.DTO(), nil
}

func (m MysqlClient) GetJobQueryTime(hash string) *time.Time {
	var res *QueryRegistry
	err := m.db.First(&res, `hash = ?`, hash).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		m.log.Errorf("failed to find QueryRegistry: %v", err)
	}

	if res != nil {
		return &res.UpdatedAt
	}
	return nil
}

func (m MysqlClient) AddQuery(hash string) {
	var res *QueryRegistry
	err := m.db.First(&res, `hash = ?`, hash).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		m.db.Create(&QueryRegistry{Hash: hash})
		return
	}

	res.UpdatedAt = time.Now()
	m.db.Save(res)
}

func (m MysqlClient) GetAuthor(key string, includeBooks bool) (author *dto.Author, err error) {
	var res *Author

	tx := m.db.Model(&Author{})
	if includeBooks {
		tx.Preload(`Books`)
	}
	err = tx.Where(`key`, key).First(&res).Error
	if err == nil {
		author = res.DTO(includeBooks)
	}

	return
}

func (m MysqlClient) GetAuthors(page, size int) (authors []*dto.Author, err error) {
	var res []*Author

	tx := m.db.Model(&Author{}).Preload(`Books`)
	m.paginate(tx, page, size)
	err = tx.Find(&res).Error
	if err != nil {
		return
	}

	authors = make([]*dto.Author, len(res))
	for i, a := range res {
		authors[i] = a.DTO(true)
	}

	return
}

func (m MysqlClient) SearchAuthors(params *dto.Author) (authors []*dto.Author, err error) {
	p := makeAuthor(params)

	var res []*Author

	tx := m.db.Model(&Author{}).Preload(`Books`).Where(p)
	err = tx.Find(&res).Error
	if err != nil {
		return
	}

	authors = make([]*dto.Author, len(res))
	for i, a := range res {
		authors[i] = a.DTO(true)
	}

	return
}

func (m MysqlClient) SaveAuthor(author *dto.Author) (err error) {
	authorModel := makeAuthor(author)

	var orig *Author
	err = m.db.Model(&Author{}).Where(&Author{Key: author.Key}).First(&orig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m.db.Create(authorModel).Error
		}

		return err
	}

	// set orig params
	authorModel.ID = orig.ID
	authorModel.CreatedAt = orig.CreatedAt
	authorModel.UpdatedAt = time.Now()

	return m.db.Save(authorModel).Error
}

func (m MysqlClient) GetBook(key string, includeAuthors bool) (book *dto.Book, err error) {
	var res *Book

	tx := m.db.Model(&Book{})
	if includeAuthors {
		tx.Preload(`Authors`)
	}
	err = tx.Where(`key`, key).First(&res).Error
	if err == nil {
		book = res.DTO(includeAuthors)
	}

	return
}

func (m MysqlClient) GetBooks(page, size int) (books []*dto.Book, err error) {
	var res []*Book

	tx := m.db.Model(&Book{}).Preload(`Authors`)
	m.paginate(tx, page, size)

	err = tx.Find(&res).Error
	if err != nil {
		return
	}

	books = make([]*dto.Book, len(res))
	for i, a := range res {
		books[i] = a.DTO(true)
	}

	return
}

func (m MysqlClient) SearchBooks(params *dto.Book) (books []*dto.Book, err error) {
	p := makeBook(params)
	var res []*Book

	tx := m.db.Model(&Book{}).Preload(`Authors`).Where(p)
	err = tx.Find(&res).Error
	if err != nil {
		return
	}

	books = make([]*dto.Book, len(res))
	for i, a := range res {
		books[i] = a.DTO(true)
	}

	return
}

func (m MysqlClient) SaveBook(book *dto.Book) (err error) {
	newBook := makeBook(book)

	var orig *Book
	err = m.db.Model(&Book{}).Where(&Book{Key: book.Key}).First(&orig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m.db.Create(newBook).Error
		}

		return err
	}

	// set orig params
	newBook.ID = orig.ID
	newBook.CreatedAt = orig.CreatedAt
	newBook.UpdatedAt = time.Now()

	return m.db.Save(newBook).Error
}

func (m MysqlClient) paginate(tx *gorm.DB, page, size int) {
	offset := size * (page - 1)
	if page == 1 {
		offset = 1
	}
	limit := size

	tx.Offset(offset).Limit(limit)
}

func NewMysqlClient(config *Config, logger *logrus.Logger) (*MysqlClient, error) {
	var (
		retries int
		err     error
		db      *gorm.DB
	)

	for retries < 5 {
		db, err = gorm.Open(mysql.Open(config.DSN()))
		if err == nil {
			break
		}

		logger.Info(`Waiting to retry MYSQL connect...`)
		time.Sleep(5 * time.Second)
		retries++
	}

	if err != nil {
		return nil, err
	}

	return &MysqlClient{
		db:  db,
		log: logger,
	}, nil
}
