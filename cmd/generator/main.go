package main

import (
	"encoding/json"
	"fmt"
)

type Element struct {
	ID     string `json:"id"`
	Type   string
	Broker string
	Port   string
	UseTLS bool
	Wires  [][]string
}

type Fields map[string]interface{}

type Foo struct {
	Elements []Element
	Fields   []Fields
}

func main() {
	f := Foo{}
	if err := json.Unmarshal([]byte(s), &f.Elements); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(s), &f.Fields); err != nil {
		panic(err)
	}

	for idx, element := range f.Elements {
		if element.Type == "mqtt-broker" {
			fmt.Printf("MQTT Broker | Host: %s, Port: %s, UseTLS: %v\n",
				element.Broker, element.Port, element.UseTLS)
			continue
		}
		fmt.Printf("Node | ID: %s, Type: %s, Wires: %v", element.ID, element.Type, element.Wires)
		if element.Type == "mqtt in" {
			fmt.Printf(", Topic: %s", f.Fields[idx]["topic"])
		}
		fmt.Println()
	}
}
