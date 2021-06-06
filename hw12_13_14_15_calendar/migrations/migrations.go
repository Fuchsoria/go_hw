package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE 'events' (
		'id' TEXT NOT NULL,
		'title' TEXT NOT NULL,
		'date' BIGINT NOT NULL,
		'duration_until' BIGINT NOT NULL,
		'description' TEXT NOT NULL,
		'owner_id' TEXT NOT NULL,
		'notice_before' BIGINT NOT NULL,
		UNIQUE KEY 'id' ('id') USING BTREE,
		PRIMARY KEY ('id')
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE 'notices' (
		'event_id' TEXT NOT NULL,
		'title' TEXT NOT NULL,
		'date' BIGINT NOT NULL,
		'user_id' TEXT NOT NULL,
		UNIQUE KEY 'event_id' ('event_id') USING BTREE,
		PRIMARY KEY ('event_id')
	);`)
	if err != nil {
		return err
	}

	return nil
}

func Down(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE 'events';")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DROP TABLE 'notices';")
	if err != nil {
		return err
	}

	return nil
}
