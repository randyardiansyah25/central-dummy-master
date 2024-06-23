package databasefactory

import (
	"central-dummy-master/repository/databasefactory/drivers"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/randyardiansyah25/libpkg/util/env"

	_ "github.com/go-sql-driver/mysql"
)

func newMysqlImpl() Database {
	return &mysqlImpl{}
}

type mysqlImpl struct {
	conn   *sql.DB
	prefix string
}

func (m *mysqlImpl) Connect() error {

	//Mysql Connection String Format : [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//Example : root:123456@tcp(127.0.0.1:3306)/employees?charset=utf8

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		os.Getenv(m.prefix+"mysql.username"),
		os.Getenv(m.prefix+"mysql.password"),
		os.Getenv(m.prefix+"mysql.address"),
		os.Getenv(m.prefix+"mysql.port"),
		os.Getenv(m.prefix+"mysql.name"),
	)
	maxPool := env.GetInt(m.prefix + "mysql.maxpoolsize")
	maxIdleConn := env.GetInt(m.prefix + "mysql.maxidleconn")
	maxLifeTime := env.GetInt(m.prefix + "mysql.maxconnlifetime")
	var err error
	if m.conn, err = sql.Open(drivers.MYSQL, connectionString); err != nil {
		return err
	}

	m.conn.SetMaxOpenConns(maxPool)
	m.conn.SetMaxIdleConns(maxIdleConn)
	m.conn.SetConnMaxLifetime(time.Minute * time.Duration(maxLifeTime))

	return nil
}

func (m *mysqlImpl) Ping() error {
	return m.conn.Ping()
}

func (m *mysqlImpl) GetConnection() interface{} {
	return m.conn
}

func (m *mysqlImpl) GetDriverName() string {
	return drivers.MYSQL
}

func (m *mysqlImpl) SetEnvironmentVariablePrefix(prefix string) {
	m.prefix = prefix
}

func (m *mysqlImpl) Close() {
	_ = m.conn.Close()
}
