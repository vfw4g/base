package vfwrepository

type RDSConfig interface {
	UserName() string
	Password() string
	Host() string
	Port() int
	Schema() string
	DriverName() string
	TablePrefix() string
	ShowSQL() bool
	LogLevel() string
}
