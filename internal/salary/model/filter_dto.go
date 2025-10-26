// internal/salary/model/filter_dto.go
package model

// FilterDTO — параметры фильтрации
type FilterDTO struct {
	Tech          *string `json:"tech,omitempty"`
	Type          *string `json:"type,omitempty"`
	SalaryMin     *int    `json:"salary_min,omitempty"`
	SalaryMax     *int    `json:"salary_max,omitempty"`
	ExperienceMin *int    `json:"experience_min,omitempty"`
	ExperienceMax *int    `json:"experience_max,omitempty"`
	CreatedAtFrom *string `json:"created_at,omitempty"` // интерпретируется как ">= created_at"
}
