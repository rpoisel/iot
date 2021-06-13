package nodered

import (
	"encoding/json"
)

const (
	typeMQTTBroker = "mqtt-broker"
)

type AllFields map[string]interface{}

func NodeREDExportParse(cfg string) (*ModelNodeRED, error) {
	n := struct {
		Elements []AllFields
	}{
		Elements: []AllFields{},
	}
	if err := json.Unmarshal([]byte(cfg), &n.Elements); err != nil {
		return nil, err
	}

	m := ModelNodeRED{
		elements: make(map[string]AllFields),
	}

	for _, element := range n.Elements {
		fieldID, ok := element["id"]
		if !ok {
			continue
		}
		id, ok := fieldID.(string)
		if !ok {
			continue
		}
		m.elements[id] = element
	}

	return &m, nil
}
