package models

/*
Config is used to store the configuration key value pair. it is used to store the configuration in the database.
*/
type Config struct {
	Key   string
	Value string
	// UpdatedBy is the user who updated the configuration. It will taken from header (X-User-Id) of the request. it will be empty if header is not present.
	UpdatedBy string
	UpdatedAt int64
	CreatedAt int64
}
