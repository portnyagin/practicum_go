package base_db_handler

type DbHandler interface {
	Execute(statement string, args ...interface{}) error
	Query(statement string, args ...interface{}) (Rows, error)
	QueryRow(statement string, args ...interface{}) (Row, error)
	Close()
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type Row interface {
	Scan(dest ...interface{}) error
}
