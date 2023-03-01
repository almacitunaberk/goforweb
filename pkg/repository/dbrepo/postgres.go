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

func (m *PostgresDBRepo) SearchAvailabilityByDatesByRoomID(startDate, endDate time.Time, roomID int) (bool, error) {
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

func (m *PostgresDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel();

	query := `SELECT r.id, r.room_name FROM rooms r WHERE r.id NOT IN
	(SELECT rr.room_id FROM room_restrictions rr WHERE rr.start_date < $1 and rr.end_date > $2)`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil

}


// Get a room by a specific ID
func (m *PostgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := `SELECT id, room_name, created_at, updated_at FROM rooms WHERE id = $1`

	var room models.Room

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}