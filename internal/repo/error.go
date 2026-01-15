package repo

import "errors"

var (
	DB_FAILURE            = errors.New("database operation failed")
	DB_URL_NOT_SET        = errors.New("database URL not set in environment")
	DB_CONNECTION_FAILURE = errors.New("database connection failed")
	REC_NOT_FOUND         = errors.New("record not found")
	DUPICATE_RECORD       = errors.New("duplicate record")
)
