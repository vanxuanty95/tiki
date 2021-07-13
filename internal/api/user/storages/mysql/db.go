package mysql

import (
	"context"
	"database/sql"
	"tiki/internal/api/user/storages"
)

const (
	GetUserByID = "SELECT id, password FROM user WHERE id = ?"
)

// MySQLDB for working with mysql
type MySQLDB struct {
	DB *sql.DB
}

// GetUserByID returns bookings if match userID AND password
func (l *MySQLDB) GetUserByID(ctx context.Context, userID string) (*storages.User, error) {
	row := l.DB.QueryRowContext(ctx,
		GetUserByID,
		userID,
	)
	user := &storages.User{}
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
