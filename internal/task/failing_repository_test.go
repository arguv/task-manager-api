package task

import "errors"

type failingRepository struct{ err error }

func (r failingRepository) List() ([]Task, error)     { return nil, r.err }
func (r failingRepository) Get(int) (Task, error)     { return Task{}, r.err }
func (r failingRepository) Create(Task) (Task, error) { return Task{}, r.err }
func (r failingRepository) Update(Task) (Task, error) { return Task{}, r.err }
func (r failingRepository) Delete(int) error          { return r.err }

var _ Repository = failingRepository{}

func newFailingHandler() *Handler {
	return NewHandler(NewService(failingRepository{err: errors.New("database unavailable")}))
}
