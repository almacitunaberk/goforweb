package repository

import (
	"time"

	"github.com/almacitunaberk/goforweb/pkg/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(rest models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(startDate, endDate time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserById(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)
}