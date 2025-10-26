// internal/salary/model/dto.go
package model

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type CreateSalaryDTO struct {
	Tech          string `json:"tech"`
	SalaryMin     int    `json:"salary_min"`
	SalaryMax     int    `json:"salary_max"`
	Type          string `json:"type"`
	ExperienceMin int    `json:"experience_min"`
	ExperienceMax int    `json:"experience_max"`
}

// Validate проверяет корректность DTO
func (dto *CreateSalaryDTO) Validate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Tech, validation.Required, validation.Length(1, 50)),
		validation.Field(&dto.Type, validation.Required, validation.In("remote", "relocate", "office")),
		validation.Field(&dto.SalaryMin, validation.Min(0)),
		validation.Field(&dto.SalaryMax, validation.Min(0)),
		validation.Field(&dto.ExperienceMin, validation.Min(0), validation.Max(50)),
		validation.Field(&dto.ExperienceMax, validation.Min(0), validation.Max(50)),
	)
}

// ToEntity преобразует DTO в сущность Salary
func (dto *CreateSalaryDTO) ToEntity() *Salary {
	return &Salary{
		Tech:          dto.Tech,
		SalaryMin:     dto.SalaryMin,
		SalaryMax:     dto.SalaryMax,
		Type:          dto.Type,
		ExperienceMin: dto.ExperienceMin,
		ExperienceMax: dto.ExperienceMax,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}
}
