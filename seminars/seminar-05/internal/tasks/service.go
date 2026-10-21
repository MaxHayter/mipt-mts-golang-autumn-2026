// Package tasks stores the seminar's tasks in memory.
package tasks

import "sync"

type Task struct {
	ID   int    `json:"id" xml:"id"`
	Text string `json:"text" xml:"text"`
}

type Service struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

func New(initial []Task) *Service {
	service := &Service{tasks: append([]Task(nil), initial...)}
	for _, task := range initial {
		if task.ID >= service.nextID {
			service.nextID = task.ID + 1
		}
	}
	if service.nextID == 0 {
		service.nextID = 1
	}
	return service
}

func (s *Service) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Task(nil), s.tasks...)
}

func (s *Service) Create(text string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{ID: s.nextID, Text: text}
	s.nextID++
	s.tasks = append(s.tasks, task)
	return task
}
