package config

import (
	"log"
)

// struct for storage dependencies of logs and others
type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
