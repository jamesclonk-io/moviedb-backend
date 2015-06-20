package migration

import (
	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/mattes/migrate/migrate"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.GetLogger()
}

func RunMigrations(uri string) {
	errors, ok := migrate.UpSync(uri, "./migrations")
	if !ok {
		for _, err := range errors {
			log.Error(err)
		}
		log.Fatal("Could not migrate up database")
	}
}
