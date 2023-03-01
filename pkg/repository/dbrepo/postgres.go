package dbrepo

import (
	"context"
	"time"

	"github.com/almacitunaberk/goforweb/pkg/models"
)

func (m *PostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	var newId int

	query := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ID`

	err := m.DB.QueryRowContext(ctx, query, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (m *PostgresDBRepo) InsertRoomRestriction(rest models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, query, rest.StartDate, rest.EndDate, rest.RoomID, rest.ReservationID, rest.RestrictionID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil

}

func (m *PostgresDBRepo) SearchAvailabilityByDates(startDate, endDate time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel();

	query := `SELECT COUNT(ID) FROM room_restrictions WHERE start_date < $1 and end_date > $2 and roomID = $3`

	var numRows int

	row := m.DB.QueryRowContext(ctx,  query, startDate, endDate, roomID)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil

}