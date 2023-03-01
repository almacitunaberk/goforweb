package repository

import "github.com/almacitunaberk/goforweb/pkg/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(rest models.RoomRestriction) error
}