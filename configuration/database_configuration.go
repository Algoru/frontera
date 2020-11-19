package configuration

// DatabaseConfiguration
type DatabaseConfiguration struct {
	Host       string
	User       string
	Password   string
	Database   string
	Timeout    int64
	AuthSource string
}
