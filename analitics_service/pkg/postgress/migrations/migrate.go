package migrations

import (
	"database/sql"
)

func Up(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(createType); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(createMusicEvents); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(createUserEvents); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func Down(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(dropMusicEvents); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(dropUserEvents); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(dropType); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
