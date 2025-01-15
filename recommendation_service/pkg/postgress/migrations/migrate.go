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
	if _, err := tx.Exec(createTable); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(insertTestActionListen); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(insertTestActionLike); err != nil {
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

	if _, err := tx.Exec(dropTable); err != nil {
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
