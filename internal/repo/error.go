package repo

type RepoError = string

const (
	DB_URL_NOT_SET        RepoError = "Database URL not set"
	DB_CONNECTION_FAILURE RepoError = "Database connection failure"
)
