package utils

import "sync"

type ThreadSafeSlice struct {
	sync.Mutex
	Workers []*Worker
}

func (s *ThreadSafeSlice) Push(w *Worker) {
	s.Lock()
	defer s.Unlock()

	s.Workers = append(s.Workers, w)
}

func (s *ThreadSafeSlice) Iter(routine func(*Worker)) {
	s.Lock()
	defer s.Unlock()

	for _, worker := range s.Workers {
		routine(worker)
	}
}
