package service

import (
	"sync"
	"time"

	"cloudpulse/backend/internal/model"
)

type Store struct {
	mu     sync.Mutex
	tasks  []model.Task
	nextID int
}

func NewStore() *Store {
	return &Store{
		tasks: []model.Task{
			{ID: 1, Title: "Create CloudPulse project", Done: true, CreatedAt: time.Now().Add(-3 * time.Hour).UTC()},
			{ID: 2, Title: "Connect frontend to Go API", Done: false, CreatedAt: time.Now().Add(-2 * time.Hour).UTC()},
			{ID: 3, Title: "Add monitoring and deploy later", Done: false, CreatedAt: time.Now().Add(-1 * time.Hour).UTC()},
		},
		nextID: 4,
	}
}

func (s *Store) ListTasks() []model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

func (s *Store) AddTask(title string) model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := model.Task{ID: s.nextID, Title: title, Done: false, CreatedAt: time.Now().UTC()}
	s.nextID++
	s.tasks = append(s.tasks, task)
	return task
}

func (s *Store) UpdateTask(id int, title *string, done *bool) (model.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID != id {
			continue
		}
		if title != nil {
			s.tasks[i].Title = *title
		}
		if done != nil {
			s.tasks[i].Done = *done
		}
		return s.tasks[i], true
	}

	return model.Task{}, false
}

func (s *Store) DeleteTask(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}

	return false
}

func (s *Store) Summary() map[string]any {
	tasks := s.ListTasks()
	completed := 0
	for _, task := range tasks {
		if task.Done {
			completed++
		}
	}

	return map[string]any{
		"app":            "CloudPulse",
		"health":         "ok",
		"generatedAt":    time.Now().UTC().Format(time.RFC3339),
		"tasksTotal":     len(tasks),
		"tasksCompleted": completed,
		"tasksOpen":      len(tasks) - completed,
		"status":         "demo-ready",
	}
}
