// internal/user/repository/repository.go
package repository

import (
	"database/sql"
	"salary-bot/internal/user/model"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Upsert(userID int64, username, firstName string) error {
	query := `
		INSERT INTO users (user_id, username, first_name) 
		VALUES (?, ?, ?) 
		ON CONFLICT(user_id) 
		DO UPDATE SET 
			username = excluded.username,
			first_name = excluded.first_name,
			last_active_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, username, firstName)
	return err
}

func (r *userRepository) IncrementCalculation(userID int64) error {
	query := `
		UPDATE users 
		SET calculation_count = calculation_count + 1, last_active_at = CURRENT_TIMESTAMP
		WHERE user_id = ?
	`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *userRepository) GetAll() ([]*model.User, error) {
	rows, err := r.db.Query(`
		SELECT user_id, username, first_name, joined_at, last_active_at, calculation_count 
		FROM users 
		ORDER BY last_active_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.UserID, &u.Username, &u.FirstName, &u.JoinedAt, &u.LastActiveAt, &u.CalculationCount)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) GetNewUsersStats() (last24h, last7d int, err error) {
	now := time.Now()

	// 24 часа назад
	last24hTime := now.Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE joined_at >= ?
	`, last24hTime).Scan(&last24h)
	if err != nil {
		return 0, 0, err
	}

	// 7 дней назад
	last7dTime := now.AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE joined_at >= ?
	`, last7dTime).Scan(&last7d)
	if err != nil {
		return 0, 0, err
	}

	return last24h, last7d, nil
}
