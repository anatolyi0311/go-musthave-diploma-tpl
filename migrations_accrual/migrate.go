package migrationsaccrual

import (
	"database/sql"
	"fmt"
)

func Up(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db not init")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(createEnum); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(createTableOrders); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(createTableCatalog); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(createTableAccrual); err != nil {
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
	if db == nil {
		return fmt.Errorf("db not init")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(dropTableCatalog); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(dropTableOrders); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(dropTableAccrual); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(dropEnum); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
