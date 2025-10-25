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
	salary := &model.Salary{
		Tech:          dto.Tech,
		SalaryMin:     dto.SalaryMin,
		SalaryMax:     dto.SalaryMax,
		Type:          dto.Type,
		ExperienceMin: dto.ExperienceMin,
		ExperienceMax: dto.ExperienceMax,
		// CreatedAt будет установлен автоматически в БД
	}

	return s.repo.Create(salary)
}

func (s *salaryService) GetAll() ([]*model.Salary, error) {
	return s.repo.List()
}
