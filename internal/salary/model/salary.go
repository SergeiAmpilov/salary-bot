// internal/salary/model/salary.go
package model

// Salary — сущность для хранения в БД
type Salary struct {
	ID            int64  `json:"id"`
	Tech          string `json:"tech"`
	SalaryMin     int    `json:"salary_min"`
	SalaryMax     int    `json:"salary_max"`
	Type          string `json:"type"`
	ExperienceMin int    `json:"experience_min"`
	ExperienceMax int    `json:"experience_max"`
	CreatedAt     string `json:"created_at"` // можно использовать time.Time, но для простоты — string
}
