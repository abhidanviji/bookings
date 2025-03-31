package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/abhidanviji/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if room id is 2, then fail; otherwise pass
	if res.RoomId == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomId == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailbilityByDatesByRoomId returns true if availability exists for roomId and false if no availability exists
func (m *testDBRepo) SearchAvailbilityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	layout := "2006-01-02"
	str := "2049-12-31"

	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// test to fai the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if start date is after 2049-12-31, then return false,
	// indicating no availability
	if start.After(t) {
		return false, nil
	}

	// otherwise we have availability
	return true, nil
}

// SearchAvailibilityForAllRooms returns a slice of availabile rooms if any for give ndate range
func (m *testDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	layout := "2006-01-02"
	str := "2049-12-31"

	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// test to fai the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	// if start date is after 2049-12-31, then return false,
	// indicating no availability
	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		Id: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomById gets a room by id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
