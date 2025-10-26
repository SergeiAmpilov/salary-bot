// internal/salary/service/interface.go
package service

import "salary-bot/internal/salary/model"

type Service interface {
	Create(dto *model.CreateSalaryDTO) error
	GetAll() ([]*model.Salary, error)
	Filter(filter *model.FilterDTO) ([]*model.Salary, error)
}
