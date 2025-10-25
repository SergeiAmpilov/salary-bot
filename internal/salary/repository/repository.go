// internal/salary/repository/repository.go
package repository

import (
	"database/sql"
	"salary-bot/internal/salary/model"
)

type salaryRepository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &salaryRepository{db: db}
}

func (r *salaryRepository) Create(s *model.Salary) error {
	query := `
		INSERT INTO salaries (tech, salary_min, salary_max, type, experience_min, experience_max)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		s.Tech,
		s.SalaryMin,
		s.SalaryMax,
		s.Type,
		s.ExperienceMin,
		s.ExperienceMax,
	)
	return err
}

func (r *salaryRepository) List() ([]*model.Salary, error) {
	query := `SELECT id, tech, salary_min, salary_max, type, experience_min, experience_max, created_at FROM salaries`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salaries []*model.Salary
	for rows.Next() {
		s := &model.Salary{}
		err := rows.Scan(
			&s.ID,
			&s.Tech,
			&s.SalaryMin,
			&s.SalaryMax,
			&s.Type,
			&s.ExperienceMin,
			&s.ExperienceMax,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		salaries = append(salaries, s)
	}

	return salaries, rows.Err()
}
