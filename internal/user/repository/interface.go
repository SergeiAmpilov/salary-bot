// internal/user/repository/interface.go
package repository

import "salary-bot/internal/user/model"

type Repository interface {
	Upsert(userID int64, username, firstName string) error
	IncrementCalculation(userID int64) error
	GetAll() ([]*model.User, error)
	GetNewUsersStats() (last24h, last7d int, err error)
}
