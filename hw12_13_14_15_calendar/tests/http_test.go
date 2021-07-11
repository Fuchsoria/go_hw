package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var (
	HTTPHost           = os.Getenv("TESTS_HTTP_HOST")
	PostgresDSN        = os.Getenv("TESTS_POSTGRES_DSN")
	createdMessageText = "Created new"
	updatedMessageText = "Updated"
	removedMessageText = "Deleted"
)

type Response = struct {
	Message string `json:"message"`
}

type ByDateResponse = struct {
	Results []storage.Event `json:"results"`
}

type CreateBody = struct {
	Title string `json:"title"`
	Date  int    `json:"date"`
}

type DeleteBody = struct {
	ID string `json:"id"`
}

type ByDateBody = struct {
	Date int64 `json:"date"`
}

func init() {
	if HTTPHost == "" {
		HTTPHost = "http://0.0.0.0:5555"
	}

	if PostgresDSN == "" {
		PostgresDSN = "host=0.0.0.0 port=5432 user=postgres password=example dbname=calendar sslmode=disable"
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestHTTPCreate(t *testing.T) {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", PostgresDSN)
	if err != nil {
		panicOnErr(err)
	}

	httpUrlCreate := HTTPHost + "/event/create"
	httpUrlUpdate := HTTPHost + "/event/update"
	httpUrlDelete := HTTPHost + "/event/delete"
	httpUrlDaily := HTTPHost + "/event/daily"
	httpUrlWeekly := HTTPHost + "/event/weekly"
	httpUrlMonth := HTTPHost + "/event/month"

	t.Run("test event create", func(t *testing.T) {
		title := fmt.Sprintf("Test_%d", time.Now().Unix())

		values := CreateBody{
			Title: title,
			Date:  1624193015,
		}

		jsonData, err := json.Marshal(values)
		if err != nil {
			panicOnErr(err)
		}

		resp, err := http.Post(httpUrlCreate, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body Response

		json.NewDecoder(resp.Body).Decode(&body)

		var events []storage.Event

		err = db.Select(&events, "SELECT * FROM events WHERE title=$1", title)
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "new event should be added")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.Equal(t, createdMessageText, body.Message, "response message should be equal")
	})

	t.Run("test event update", func(t *testing.T) {
		currentUnix := time.Now().Unix()
		uuid := uuid.NewString()
		title := fmt.Sprintf("Test_%d", currentUnix)
		newTitle := fmt.Sprintf("New_%d", currentUnix)

		_, err := db.Exec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid, title, 1624193015, 0, "", "main_owner_id", 0)
		if err != nil {
			panicOnErr(err)
		}

		updatedData := storage.Event{
			ID:            uuid,
			Title:         newTitle,
			Date:          1624193015,
			DurationUntil: 0,
			Description:   "Test description",
			OwnerID:       "main_owner_id",
			NoticeBefore:  0,
		}

		jsonData, err := json.Marshal(updatedData)
		if err != nil {
			panicOnErr(err)
		}

		respUpdate, err := http.Post(httpUrlUpdate, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body Response

		json.NewDecoder(respUpdate.Body).Decode(&body)

		var events []storage.Event

		err = db.Select(&events, "SELECT * FROM events WHERE id=$1 AND title=$2", uuid, newTitle)
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "updated event should be found")
		require.Equal(t, http.StatusOK, respUpdate.StatusCode, "response statuscode should be ok")
		require.Equal(t, updatedMessageText, body.Message, "response message should be equal")
	})

	t.Run("test event delete", func(t *testing.T) {
		currentUnix := time.Now().Unix()
		uuid := uuid.NewString()
		title := fmt.Sprintf("Test_%d", currentUnix)

		_, err := db.Exec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid, title, 1624193015, 0, "", "main_owner_id", 0)
		if err != nil {
			panicOnErr(err)
		}

		deleteBody := DeleteBody{ID: uuid}

		jsonData, err := json.Marshal(deleteBody)
		if err != nil {
			panicOnErr(err)
		}

		respUpdate, err := http.Post(httpUrlDelete, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body Response

		json.NewDecoder(respUpdate.Body).Decode(&body)

		var events []storage.Event

		err = db.Select(&events, "SELECT * FROM events WHERE id=$1", uuid)
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 0, "event should be removed")
		require.Equal(t, http.StatusOK, respUpdate.StatusCode, "response statuscode should be ok")
		require.Equal(t, removedMessageText, body.Message, "response message should be equal")
	})

	t.Run("test daily events", func(t *testing.T) {
		currentUnix := time.Now().Unix()
		uuid := uuid.NewString()
		title := fmt.Sprintf("Test_%d", currentUnix)

		_, err := db.Exec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid, title, currentUnix, 0, "", "main_owner_id", 0)
		if err != nil {
			panicOnErr(err)
		}

		dateBody := ByDateBody{Date: currentUnix}

		jsonData, err := json.Marshal(dateBody)
		if err != nil {
			panicOnErr(err)
		}

		respUpdate, err := http.Post(httpUrlDaily, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body ByDateResponse

		json.NewDecoder(respUpdate.Body).Decode(&body)

		createdEventExist := false
		for _, event := range body.Results {
			if event.ID == uuid {
				createdEventExist = true
			}
		}

		require.GreaterOrEqual(t, len(body.Results), 1, "results length should be greater or equal than 1")
		require.True(t, createdEventExist, "created item should be found by date")
	})

	t.Run("test weekly events", func(t *testing.T) {
		currentUnix := time.Now().Unix()
		uuid := uuid.NewString()
		title := fmt.Sprintf("Test_%d", currentUnix)

		_, err := db.Exec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid, title, currentUnix, 0, "", "main_owner_id", 0)
		if err != nil {
			panicOnErr(err)
		}

		dateBody := ByDateBody{Date: currentUnix}

		jsonData, err := json.Marshal(dateBody)
		if err != nil {
			panicOnErr(err)
		}

		respUpdate, err := http.Post(httpUrlWeekly, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body ByDateResponse

		json.NewDecoder(respUpdate.Body).Decode(&body)

		createdEventExist := false
		for _, event := range body.Results {
			if event.ID == uuid {
				createdEventExist = true
			}
		}

		require.GreaterOrEqual(t, len(body.Results), 1, "results length should be greater or equal than 1")
		require.True(t, createdEventExist, "created item should be found by date")
	})

	t.Run("test month events", func(t *testing.T) {
		currentUnix := time.Now().Unix()
		uuid := uuid.NewString()
		title := fmt.Sprintf("Test_%d", currentUnix)

		_, err := db.Exec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid, title, currentUnix, 0, "", "main_owner_id", 0)
		if err != nil {
			panicOnErr(err)
		}

		dateBody := ByDateBody{Date: currentUnix}

		jsonData, err := json.Marshal(dateBody)
		if err != nil {
			panicOnErr(err)
		}

		respUpdate, err := http.Post(httpUrlMonth, "application/json",
			bytes.NewBuffer(jsonData))
		if err != nil {
			panicOnErr(err)
		}

		var body ByDateResponse

		json.NewDecoder(respUpdate.Body).Decode(&body)

		createdEventExist := false
		for _, event := range body.Results {
			if event.ID == uuid {
				createdEventExist = true
			}
		}

		require.GreaterOrEqual(t, len(body.Results), 1, "results length should be greater or equal than 1")
		require.True(t, createdEventExist, "created item should be found by date")
	})
}
