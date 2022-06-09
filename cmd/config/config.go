package config

import (
	"log"

	"github.com/dshurubtsov/pkg/mongodb"
)

// struct for storage dependencies of logs and others
type Application struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	UserModel *mongodb.UserModel
}
