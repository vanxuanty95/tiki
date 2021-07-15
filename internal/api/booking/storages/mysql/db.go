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

// InsertScreen will create new screen.
func (l *MySQLDB) InsertScreen(ctx context.Context, s *storages.Screen) (int, error) {
	query := `INSERT INTO screen (number_seat_row, number_seat_column, created_date, user_id) VALUES (?, ?, ?, ?)`
	res, err := l.DB.ExecContext(ctx, query, &s.NumberSeatRow, &s.NumberSeatColumn, &s.CreatedDate, &s.UserID)
	if err != nil {
		return 0, err
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(newID), nil
}

func (l *MySQLDB) GetScreenByID(ctx context.Context, id int) (*storages.Screen, error) {
	query := `SELECT * FROM screen WHERE id = ?`
	rows, err := l.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	screen := &storages.Screen{}
	var createdDateStr string
	err = rows.Scan(&screen.ID, &screen.NumberSeatRow, &screen.NumberSeatColumn, &createdDateStr, &screen.UserID)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(utils.DefaultLayoutDB, createdDateStr)
	if err != nil {
		return nil, err
	}

	screen.CreatedDate = t

	return screen, nil
}

func (l *MySQLDB) GetAllSeatByScreenID(ctx context.Context, id int) ([]*storages.Seat, error) {
	query := `SELECT * FROM seat WHERE screen_id = ?`
	rows, err := l.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []*storages.Seat
	for rows.Next() {
		var createdDateStr string
		seat := &storages.Seat{}
		err = rows.Scan(&seat.ID, &seat.Row, &seat.Column, &seat.UserID, &seat.ScreenID, &createdDateStr)
		if err != nil {
			return nil, err
		}
		t, err := time.Parse(utils.DefaultLayoutDB, createdDateStr)
		if err != nil {
			return nil, err
		}
		seat.BookedDate = t
		seats = append(seats, seat)
	}
	return seats, nil
}

func (l *MySQLDB) InsertSeats(ctx context.Context, seats []*storages.Seat) error {
	sqlStr := `INSERT INTO seat (row_id, column_id, user_id, screen_id, booked_date) VALUES `
	var vals []interface{}

	for _, row := range seats {
		sqlStr += `(?, ?, ?, ?, ?),`
		vals = append(vals, row.Row, row.Column, row.UserID, row.ScreenID, row.BookedDate)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	stmt, err := l.DB.PrepareContext(ctx, sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, vals...)
	return err
}
