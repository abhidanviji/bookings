package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/abhidanviji/bookings/internal/models"
)

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room:   models.Room{Id: 1, RoomName: "General's Quarters"},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room:   models.Room{Id: 1, RoomName: "General's Quarters"},
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("phone", "1234567890")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test without session
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned wrong response code for processing without session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned wrong response code for missing request body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	postedData.Set("first_name", "J")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	reservation.RoomId = 2
	postedData.Set("first_name", "John")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert room restriction into database
	reservation.RoomId = 1000
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler failed when trying to fail inserting room restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJson(t *testing.T) {
	// case 1: rooms are not available

	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJson)

	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	// for start date after 2049-21-31 rooms should not be available
	if j.Ok {
		t.Error("got availability when there should be no availability")
	}

	// rooms are available
	postedData.Set("start", "2030-01-01")
	postedData.Set("end", "2030-01-02")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJson)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if !j.Ok {
		t.Error("got no availability when there are some")
	}

	// missing body test
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJson)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.Ok || j.Message != "Internal server error" {
		t.Error("got availability when there is no req body")
	}

	// database error
	postedData.Set("start", "2060-01-01")
	postedData.Set("end", "2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJson)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.Ok || j.Message != "Error connecting to database" {
		t.Error("got availability when there is database error")
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// rooms not available
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error("post availability returns rooms when there are none")
	}

	// search availability with empty req body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post availability fails with empty request body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid start date
	postedData.Set("start", "invalid")
	postedData.Set("end", "2050-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post availability returns wrong response when parsing invalid date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid end date
	postedData.Set("start", "2050-01-01")
	postedData.Set("end", "invalid")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post availability returns wrong response when parsing invalid date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// search all rooms in the date that fails
	postedData.Set("start", "2060-01-01")
	postedData.Set("end", "2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post availability returns wrong response when no rooms are available: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// search rooms in the available date range
	postedData.Set("start", "2030-01-01")
	postedData.Set("end", "2030-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("post availability returns wrong response when rooms are available: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			Id:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation summary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// not in session
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation summary handler returned wrong response code when no session in request: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			Id:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	req.RequestURI = "/choose-room/1"
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("choose room handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// invalid room id
	req, _ = http.NewRequest("GET", "/choose-room/a", nil)
	req.RequestURI = "/choose-room/a"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("choose room handler returned wrong response code for invalid roomid: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// without session
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	req.RequestURI = "/choose-room/1"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("choose room handler returned wrong response code when no session in request: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			Id:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room?id=1&s=2050-01-01&e=2050-01-02", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("book room handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// can't get room from db
	req, _ = http.NewRequest("GET", "/book-room?id=3&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("book room handler returned wrong response code for unavailable room: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
