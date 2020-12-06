package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/rpoisel/iot/cmd/mqtt-db-postgres/power/public/model"
	"github.com/rpoisel/iot/cmd/mqtt-db-postgres/power/public/table"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	UTIL "github.com/rpoisel/iot/internal/util"
)

type configuration struct {
	Mqtt     UTIL.MqttConfiguration
	Postgres struct {
		Host     string
		Port     int16
		User     string
		Password string
		DbName   string
	}
}

func powerPublishHandler(_ MQTT.Client, msg MQTT.Message) {
	r, err := UTIL.NewReadings(msg.Payload())
	if err != nil {
		fmt.Printf("Unhandled message: %s", err)
	}
	fmt.Printf("%s: Solar = %d, Total = %d\n", time.Now().Format(time.UnixDate), r.Solar, r.Total)
}

func main() {
	var configPath = flag.String("c", "/etc/homeautomation.yaml", "Path to the configuration file")
	flag.Parse()

	configuration := configuration{}
	err := UTIL.ReadConfig(*configPath, &configuration)
	if err != nil {
		panic(err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			configuration.Postgres.Host,
			configuration.Postgres.Port,
			configuration.Postgres.User,
			configuration.Postgres.Password,
			configuration.Postgres.DbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mqttClient := UTIL.SetupMqtt(configuration.Mqtt, func(_ MQTT.Client, msg MQTT.Message) {
		log.Print("Unhandled MQTT message ", msg)
	}, func(client MQTT.Client) {
		client.Subscribe("/homeautomation/power/cumulative", 0 /* qos */, powerPublishHandler)
	})
	defer mqttClient.Disconnect(250)

	// using table data types
	stmt := postgres.SELECT(table.Power.Modtime, table.Power.Solar, table.Power.Total).
		FROM(table.Power).
		WHERE(table.Power.Modtime.GT(postgres.Timestamp(2020, 6, 30, 12, 20, 0)).
			AND(table.Power.Modtime.LT(postgres.Timestamp(2020, 6, 30, 12, 30, 0)))).
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

	<-stopChan
	fmt.Println("Good bye!")
}
