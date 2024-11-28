package openlibrary

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	log    *logrus.Logger
	host   string
	//mu     jobRegistry
}

func (c *Client) StartWorker(ctx context.Context, resolver JobResolver) {
	jobQueue.Work(ctx, resolver)
}

// executeJob retrieves the job request and resolves the job data
// todo: httpretryable & rate limit
func (c *Client) executeJob(job Job) error {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/%s`, c.host, job.uri()), nil)
	q := req.URL.Query()
	for k, v := range job.queryParams() {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	c.log.Debug(req.URL.String())

	res, err := c.client.Do(req)
	if err != nil {
		job.setResult(500, nil)
		return err
	}
	job.setResult(res.StatusCode, res.Body)
	return nil
}

func NewClient(log *logrus.Logger, hostPath string, maxConcurrent int) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		log:  log,
		host: hostPath,
		//jobQueue: make(chan Job, maxConcurrent),
	}

	// fixme: possible override of jobQueue global
	jobQueue = &JobQueue{
		jobQueue: make(chan Job),
		client:   c,
	}

	mu = &jobRegistry{
		jobs: make(map[string]bool),
	}

	return c
}
