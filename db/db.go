package db

import (
	"database/sql"

	"fmt"

	"github.com/ricardoaat/mass-blocker/config"
	log "github.com/sirupsen/logrus"
	//Driver init
	_ "github.com/lib/pq"
)

//Ledb Returns a DB connection
var Ledb *sql.DB

func conToDB(d config.Database) *sql.DB {
	p := "dbname=%s user=%s password=%s host=%s port=%s sslmode=disable"
	p = fmt.Sprintf(p,
		d.DBname,
		d.Username,
		d.Password,
		d.Host,
		d.Port)
	log.Info("Connecting to " + p)

	db, err := sql.Open("postgres", p)
	if err != nil {
		log.Panic(err)
		panic(err)
	}

	return db
}

/*InitDatabases Starts databases connections
 */
func InitDatabases() {
	Ledb = conToDB(config.Conf.Databases["ledb"])
}

/*CloseDatabases Stop databases connections
 */
func CloseDatabases() {
	log.Info("Closing databases")
	Lebd.Close()
}
