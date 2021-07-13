package mysql

import (
	"context"
	"database/sql"
	"tiki/internal/api/booking/storages"
	"tiki/internal/api/utils"
	"time"
)

// MySQLDB for working with mysql
type MySQLDB struct {
	DB *sql.DB
}

// RetrieveBookings returns bookings if match userID AND createDate.
func (l *MySQLDB) RetrieveBookings(ctx context.Context, userID, createdDate string, page, limit int) ([]*storages.Booking, error) {

	limit, offset := utils.GetLimitOffsetFormPageNutikier(page, limit)

	userIDSQL := sql.NullString{
		String: userID,
		Valid:  true,
	}
	createdDateSql := sql.NullString{
		String: createdDate,
		Valid:  true,
	}

	query := `SELECT id, content, user_id, created_date FROM bookings WHERE user_id = $1 AND created_date = $2 LIMIT $3 OFFSET $4`
	rows, err := l.DB.QueryContext(ctx, query, userIDSQL, createdDateSql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*storages.Booking
	var creDateType time.Time
	for rows.Next() {
		t := &storages.Booking{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &creDateType)
		if err != nil {
			return nil, err
		}
		t.CreatedDate = creDateType.Format(utils.DefaultLayout)
		bookings = append(bookings, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// AddBooking adds a new booking to DB
func (l *MySQLDB) AddBooking(ctx context.Context, t *storages.Booking) error {
	query := `INSERT INTO bookings (id, content, user_id, created_date) VALUES ($1, $2, $3, $4)`
	if _, err := l.DB.ExecContext(ctx, query, &t.ID, &t.Content, &t.UserID, &t.CreatedDate); err != nil {
		return err
	}
	return nil
}
