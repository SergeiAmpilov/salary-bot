// internal/salary/model/dto.go
package model

// CreateSalaryDTO — данные, приходящие в POST /salary
type CreateSalaryDTO struct {
	Tech          string `json:"tech" validate:"required"`
	SalaryMin     int    `json:"salary_min"`
	SalaryMax     int    `json:"salary_max"`
	Type          string `json:"type" validate:"required"`
	ExperienceMin int    `json:"experience_min"`
	ExperienceMax int    `json:"experience_max"`
}
