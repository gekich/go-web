package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Jwt struct {
	Secret      string        `required:"true"`
	TokenExpiry time.Duration `split_words:"true" default:"24h"`
}

func JWT() Jwt {
	var jwt Jwt
	envconfig.MustProcess("JWT", &jwt)

	return jwt
}
