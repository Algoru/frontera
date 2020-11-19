package configuration

// SecurityConfiguration
type SecurityConfiguration struct {
	MinPasswordLength uint16
	BCryptCost        uint8
	AuthUseField      string
	TokenSigningKey   string
	TokenIssuer       string
	TokenLifetime     int64
}

// TODO (@Algoru): add encrypt algorithm of choice
