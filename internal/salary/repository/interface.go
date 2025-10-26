// internal/salary/repository/interface.go
package repository

import "salary-bot/internal/salary/model"

type Repository interface {
	Create(salary *model.Salary) error
	List() ([]*model.Salary, error)
	Filter(filter *model.FilterDTO) ([]*model.Salary, error)
}
