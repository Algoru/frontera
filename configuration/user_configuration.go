package configuration

// UserConfiguration
type UserConfiguration struct {
	AllowDuplicateEmail   bool
	AllowNullPayload      bool
	AllowMultipleSessions bool
	Payload               map[string]UserPayloadConfiguration
}

// UserPayloadConfiguration
type UserPayloadConfiguration struct {
	Required     bool
	Format       string
	DefaultValue interface{}
	Max          int64
	Min          int64
}
