// internal/salary/repository/repository.go
package repository

import (
	"database/sql"
	"salary-bot/internal/salary/model"
	"strings"
	"time"
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

func (r *salaryRepository) Filter(f *model.FilterDTO) ([]*model.Salary, error) {
	var query strings.Builder
	query.WriteString(`
		SELECT id, tech, salary_min, salary_max, type, experience_min, experience_max, created_at
		FROM salaries
		WHERE 1=1
	`)

	var args []interface{}

	if f.Tech != nil {
		query.WriteString(" AND tech = ?")
		args = append(args, *f.Tech)
	}

	if f.Type != nil {
		query.WriteString(" AND type = ?")
		args = append(args, *f.Type)
	}

	if f.SalaryMin != nil {
		query.WriteString(" AND salary_max >= ?") // перекрывается, если зарплата в диапазоне
		args = append(args, *f.SalaryMin)
	}

	if f.SalaryMax != nil {
		query.WriteString(" AND salary_min <= ?")
		args = append(args, *f.SalaryMax)
	}

	if f.ExperienceMin != nil {
		query.WriteString(" AND experience_max >= ?")
		args = append(args, *f.ExperienceMin)
	}

	if f.ExperienceMax != nil {
		query.WriteString(" AND experience_min <= ?")
		args = append(args, *f.ExperienceMax)
	}

	if f.CreatedAtFrom != nil {
		// Проверим формат (опционально)
		_, err := time.Parse("2006-01-02 15:04:05", *f.CreatedAtFrom)
		if err != nil {
			// Можно вернуть ошибку, но для простоты проигнорируем или вернём пустой результат
			return []*model.Salary{}, nil
		}
		query.WriteString(" AND created_at >= ?")
		args = append(args, *f.CreatedAtFrom)
	}

	rows, err := r.db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	salaries := make([]*model.Salary, 0)
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
