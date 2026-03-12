package config

type Auth struct {
	JwtSecret               string `yaml:"jwt_secret"`
	RefreshTokenSecret      string `yaml:"refresh_token_secret"`
	TokenExpiration         int64  `yaml:"token_expiration"`
	RefreshTokenExpiration  int64  `yaml:"refresh_token_expiration"`
	EnforcerDurationSeconds int64  `yaml:"enforcer_duration_seconds"`
}
