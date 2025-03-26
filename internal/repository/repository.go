package repository

import (
	"time"

	"github.com/abhidanviji/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailbilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error)
	SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error)
}
