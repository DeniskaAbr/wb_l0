package config

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Postgres_host string
	Postgres_port string
	Postgres_db   string

	Postgres_user     string
	Postgres_password string

	Nats_cluster_id string
	Nats_hostname   string

	Http_port string
}

func (c *Config) InitFromEnv() error {
	errorMessage := ""

	cv := reflect.ValueOf(*c)
	typeOfS := cv.Type()

	for i := 0; i < cv.NumField(); i++ {

		// fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, cv.Field(i).Interface())

		obj := reflect.ValueOf(c).Elem()

		v := os.Getenv(strings.ToUpper(typeOfS.Field(i).Name))

		if os.Getenv(strings.ToUpper(typeOfS.Field(i).Name)) == "" {

			errorMessage = errorMessage + typeOfS.Field(i).Name

		}

		obj.FieldByName(typeOfS.Field(i).Name).SetString(v)

	}

	if len(errorMessage) > 0 {
		return errors.New("value for parameters is missing: " + errorMessage)
	}

	return nil
}
