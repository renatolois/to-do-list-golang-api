package task

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(title string) (Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return Task{}, errors.New("title cannot be left blank")
	}

	if len(title) < 3 {
		return Task{}, errors.New("the title must be length >= 3")
	}

	t := Task{
		Id:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:       title,
		Completed:   false,
		CompletedAt: sql.NullTime{},
		CreatedAt:   time.Now(),
	}

	created, err := s.repo.Create(t)
	if err != nil {
		return Task{}, err
	}

	return created, nil
}

func (s *Service) ListTasks() ([]Task, error) {
	return s.repo.GetAll()
}

func (s *Service) GetTask(id string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, errors.New("invalid id")
	}

	return s.repo.GetByID(id)
}

func (s *Service) DeleteTask(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("invalid id")
	}

	return s.repo.DeleteByID(id)
}
