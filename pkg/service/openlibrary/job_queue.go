package openlibrary

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

// jobQueue lives in the global space to be accessible by Jobs
var jobQueue *JobQueue
var mu *jobRegistry

type Queue interface {
	Work(ctx context.Context)
}

type jobRegistry struct {
	sync.RWMutex
	jobs map[string]bool
}

type JobQueue struct {
	client   *Client
	jobQueue chan Job
}

func (q *JobQueue) hasJob(hash string) bool {
	q.client.log.Debug(`has job`)
	mu.RLock()
	defer mu.RUnlock()

	_, ok := mu.jobs[hash]
	return ok
}

func (q *JobQueue) clearJob(hash string) {
	delete(mu.jobs, hash)
	q.client.log.Debugf(`cleared job %v`, hash)
}

func (q *JobQueue) queue(job Job) {
	q.client.log.Debugf(`queue job %s`, job.Hash())

	if q.hasJob(job.Hash()) {
		q.client.log.Debugf(`ignoring job %v`, job.Hash())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	q.client.log.Debugf(`starting new job %v %v`, reflect.TypeOf(job), job.Hash())
	mu.jobs[job.Hash()] = true

	//q.client.log.Debug(len(q.jobQueue))

	q.jobQueue <- job // fixme: queue locks here after a few passes
	//q.client.log.Debug(len(q.jobQueue))

	q.client.log.Debugf(`exit queue job %s`, job.Hash())
}

func (q *JobQueue) Work(ctx context.Context, resolver JobResolver) {
	// fixme: figure out reason for thread lock
	for {
		select {
		case <-ctx.Done():
			q.client.log.Info(`OpenLibrary worker stopped`)
			return
		case j := <-q.jobQueue:
			q.client.log.WithFields(logrus.Fields{`params`: j.queryParams()}).Debugf(`Resolving job: %v`, j.Hash())
			if err := q.client.executeJob(j); err != nil {
				q.client.log.Errorf(`error in response: %d %v`, j.StatusCode(), err)
				continue
			}
			_, err := j.decode()
			if err != nil {
				q.client.log.WithFields(logrus.Fields{`err`: err}).Error(`Failed to resolve job`)
			} else {
				resolver(j)
			}

			q.clearJob(j.Hash())
		}
	}

}
