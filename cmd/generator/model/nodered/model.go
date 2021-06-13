package nodered

import "github.com/rpoisel/iot/cmd/generator/model"

type ModelNodeRED struct {
	elements map[string]AllFields
}

func (m *ModelNodeRED) Visit(visitor model.Visitor) {
	for _, el := range m.elements {
		fieldElementType, ok := el["type"]
		if !ok {
			continue
		}
		elementType, ok := fieldElementType.(string)
		if !ok {
			continue
		}
		if elementType == "mqtt in" {
			visitor.MQTTIn(
				&MQTTIn{
					name:   "foo",
					broker: "bar",
				},
			)
		}
	}
}
