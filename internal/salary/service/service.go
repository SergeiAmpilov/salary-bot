// internal/salary/service/service.go
package service

import (
	"salary-bot/internal/salary/model"
	"salary-bot/internal/salary/repository"
)

type salaryService struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &salaryService{repo: repo}
}

func (s *salaryService) Create(dto *model.CreateSalaryDTO) error {
	salary := dto.ToEntity()
	return s.repo.Create(salary)
}

func (s *salaryService) GetAll() ([]*model.Salary, error) {
	return s.repo.List()
}

func (s *salaryService) Filter(filter *model.FilterDTO) ([]*model.Salary, error) {
	return s.repo.Filter(filter)
}
