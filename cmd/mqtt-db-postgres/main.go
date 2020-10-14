package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/rpoisel/IoT/cmd/mqtt-db-postgres/power/public/model"
	. "github.com/rpoisel/IoT/cmd/mqtt-db-postgres/power/public/table"

	UTIL "github.com/rpoisel/IoT/internal/util"
)

type Configuration struct {
	Mqtt     UTIL.MqttConfiguration
	Postgres struct {
		Host     string
		Port     int16
		User     string
		Password string
		DbName   string
	}
}

func main() {
	var configPath = flag.String("c", "/etc/homeautomation.json", "Path to the configuration file")
	flag.Parse()

	configuration := Configuration{}
	err := UTIL.ReadConfig(*configPath, &configuration)
	if err != nil {
		panic(err)
	}

	var connectString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configuration.Postgres.Host,
		configuration.Postgres.Port,
		configuration.Postgres.User,
		configuration.Postgres.Password,
		configuration.Postgres.DbName)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// using table data types
	stmt := SELECT(Power.Modtime, Power.Solar, Power.Total).
		FROM(Power).
		WHERE(Power.Modtime.GT(Timestamp(2020, 6, 30, 12, 20, 0)).
			AND(Power.Modtime.LT(Timestamp(2020, 6, 30, 12, 30, 0)))).
		LIMIT(20)
	var dest []model.Power
	err = stmt.Query(db, &dest)
	if err != nil {
		panic(err)
	}

	// using model data types
	for _, record := range dest {
		fmt.Printf("%s: Solar = %d, Total = %d\n", record.Modtime.Format(time.UnixDate), record.Solar, record.Total)
	}
}
