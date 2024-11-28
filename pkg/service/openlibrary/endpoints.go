package openlibrary

import "fmt"

func (c *Client) GetBook(key string) Job {
	return NewJob(JobTypeBook, map[string]string{
		"key": key,
	})
}

func (c *Client) GetBookByISBN(isbn uint64) Job {
	return NewJob(JobTypeISBN, map[string]string{
		"isbn": fmt.Sprint(isbn),
	})
}

func (c *Client) GetAuthor(key string) Job {
	return NewJob(JobTypeAuthor, map[string]string{
		"key": key,
	})
}

func (c *Client) SearchBooks(title, language string) Job {
	params := map[string]string{
		"q": title,
	}
	if language != `` {
		params[`language`] = language
	}

	return NewJob(JobTypeBookSearch, params)
}

func (c *Client) SearchAuthors(name string) Job {
	return NewJob(JobTypeAuthorSearch, map[string]string{
		"q": name,
	})
}

func (c *Client) GetWork(key string) Job {
	return NewJob(JobTypeWork, map[string]string{
		"key": key,
	})
}
