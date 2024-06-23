package databasefactory

var AppDb Database

type Database interface {
	Connect() error
	Ping() error
	GetConnection() interface{}
	GetDriverName() string
	SetEnvironmentVariablePrefix(string)
	Close()
}
