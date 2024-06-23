package databasefactory

import (
	"central-dummy-master/repository/databasefactory/drivers"
	"errors"
	"os"
)

func GetDatabase() (db Database, err error) {
	driverName := os.Getenv("app.database_driver")
	if driverName == "" {
		driverName = drivers.MYSQL
	}

	if driverName == drivers.MYSQL {
		return newMysqlImpl(), nil
	} else {
		return db, errors.New("unimplement database driver")
	}
}
