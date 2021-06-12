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

const s = `[{
	"id": "f368fc00.a992a",
	"type": "tab",
	"label": "Flow 1",
	"disabled": false,
	"info": ""
}, {
	"id": "2e85ad93.0f5202",
	"type": "ui_button",
	"z": "f368fc00.a992a",
	"name": "",
	"group": "440313ac.6f2e7c",
	"order": 3,
	"width": 0,
	"height": 0,
	"passthru": false,
	"label": "",
	"tooltip": "",
	"color": "",
	"bgcolor": "",
	"icon": "fa-chevron-circle-down",
	"payload": "{\"src\": \"buttonDownSR\", \"value\": true }",
	"payloadType": "json",
	"topic": "",
	"x": 575,
	"y": 300,
	"wires": [
		["ba8d17d5.61c578"]
	],
	"l": false
}, {
	"id": "756aeb0c.5d2074",
	"type": "ui_button",
	"z": "f368fc00.a992a",
	"name": "",
	"group": "440313ac.6f2e7c",
	"order": 1,
	"width": 0,
	"height": 0,
	"passthru": false,
	"label": "",
	"tooltip": "",
	"color": "",
	"bgcolor": "",
	"icon": "fa-chevron-circle-up",
	"payload": "{\"src\": \"buttonUpSR\", \"value\": true }",
	"payloadType": "json",
	"topic": "",
	"x": 575,
	"y": 260,
	"wires": [
		["ba8d17d5.61c578"]
	],
	"l": false
}, {
	"id": "3e29935d.e99dbc",
	"type": "inject",
	"z": "f368fc00.a992a",
	"name": "",
	"repeat": "",
	"crontab": "",
	"once": true,
	"onceDelay": 0.1,
	"topic": "",
	"payload": "Init",
	"payloadType": "str",
	"x": 170,
	"y": 60,
	"wires": [
		["f8a6d26a.da20a"]
	]
}, {
	"id": "f8a6d26a.da20a",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "Init",
	"func": "// Library\nclass MAX7311\n{\n    constructor(bus, addr)\n    {\n        this.bus = bus;\n        this.addr = addr;\n        this.pinValues = 0b00000000;\n        \n        this.bus.writeByteSync(this.addr, 0x06, 0x00);\n    }\n    writePinSync(pin, value)\n    {\n        writePin(pin, value);\n        commit();\n    }\n    writePin(pin, value)\n    {\n        if (value)\n        {\n            this.pinValues |= (0x1 << pin);\n        }\n        else\n        {\n            this.pinValues &= (~(0x01 << pin) & 255);\n        }\n    }\n    commit()\n    {\n        var inverted = (~this.pinValues & 0xff);\n        this.bus.writeByteSync(this.addr, 0x02, inverted);\n    }\n}\n\nclass PCF8574\n{\n    constructor(bus, addr)\n    {\n        this.bus = bus;\n        this.addr = addr;\n        this.pinValues = i2c1.receiveByteSync(this.addr);\n    }\n    readPins()\n    {\n        this.pinValues = i2c1.receiveByteSync(this.addr);\n    }\n    getPins()\n    {\n        return this.pinValues;\n    }\n    getPin(pin)\n    {\n        return (this.pinValues & (0x01 << pin)) == (0x01 << pin);\n    }\n}\n\nvar Blind = function (name) {\n    this.name = name;\n    this.state = new StateStop(this);\n    this.change = function(state) {\n        this.state = state;   \n    };\n    this.button1 = function() {\n        return this.state.button1();\n    };\n    this.button2 = function() {\n        return this.state.button2();\n    };\n}\nvar StateUp = function(blind) {\n    this.blind = blind;\n    this.button2 = this.button1 = function() {\n        var stoppedState = new StateStop(blind);\n        blind.change(stoppedState);\n        return stoppedState.getOutputEvent();\n    };\n    this.getOutputEvent = function() {\n        return {\n            payload: {\n                cmd: 'up'\n            }\n        };\n    };\n};\nvar StateDown = function(blind) {\n    this.blind = blind;\n    this.button2 = this.button1 = function() {\n        var stoppedState = new StateStop(blind);\n        blind.change(stoppedState);\n        return stoppedState.getOutputEvent();\n    };\n    this.getOutputEvent = function() {\n        return {\n            payload: {\n                cmd: 'down'\n            }\n        };\n    };\n};\nvar StateStop = function(blind) {\n    this.blind = blind;\n    this.button1 = function() {\n        var upState = new StateUp(blind);\n        blind.change(upState);\n        return upState.getOutputEvent();\n    };\n    this.button2 = function() {\n        var downState = new StateDown(blind);\n        blind.change(downState);\n        return downState.getOutputEvent();\n    };\n    this.getOutputEvent = function() {\n        return {\n            payload: {\n                cmd: 'stop'\n            }\n        };\n    };\n};\n\n// Logic\nvar i2cbus = global.get(\"i2cbus\");\nvar i2c1 = i2cbus.openSync(1);\nflow.set(\"i2c1\", i2c1);\n\nnode.on('close', function() {\n    i2c1.closeSync(); \n});\n\nflow.set('MAX7311_1', new MAX7311(i2c1, 0x20));\nflow.set('PCF8574_1', new PCF8574(i2c1, 0x3b));\nflow.set('blindSR', new Blind('blindSR'));\nflow.set('blindKiZi2', new Blind('blindKiZi2'));\n\nreturn msg;",
	"outputs": 1,
	"noerr": 0,
	"x": 770,
	"y": 60,
	"wires": [
		[]
	]
}, {
	"id": "587860b.f3857a",
	"type": "ui_button",
	"z": "f368fc00.a992a",
	"name": "",
	"group": "7c3d9abe.f8c944",
	"order": 3,
	"width": 0,
	"height": 0,
	"passthru": false,
	"label": "",
	"tooltip": "",
	"color": "",
	"bgcolor": "",
	"icon": "fa-chevron-circle-down",
	"payload": "{\"src\": \"buttonDownKiZi2\", \"value\": true }",
	"payloadType": "json",
	"topic": "",
	"x": 575,
	"y": 620,
	"wires": [
		["647b7ca5.24cb44"]
	],
	"l": false
}, {
	"id": "8df35a3b.0976c8",
	"type": "ui_button",
	"z": "f368fc00.a992a",
	"name": "",
	"group": "7c3d9abe.f8c944",
	"order": 1,
	"width": 0,
	"height": 0,
	"passthru": false,
	"label": "",
	"tooltip": "",
	"color": "",
	"bgcolor": "",
	"icon": "fa-chevron-circle-up",
	"payload": "{\"src\": \"buttonUpKiZi2\", \"value\": true }",
	"payloadType": "json",
	"topic": "",
	"x": 575,
	"y": 580,
	"wires": [
		["647b7ca5.24cb44"]
	],
	"l": false
}, {
	"id": "c7aa4dca.0897f",
	"type": "inject",
	"z": "f368fc00.a992a",
	"name": "",
	"repeat": ".1",
	"crontab": "",
	"once": false,
	"onceDelay": 0.1,
	"topic": "",
	"payload": "",
	"payloadType": "date",
	"x": 110,
	"y": 420,
	"wires": [
		["25075186.b3ab9e"]
	]
}, {
	"id": "25075186.b3ab9e",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "IdentifyPCF8574Events",
	"func": "var getBitValue = function(values, pos)\n{\n    return (values & (0x01 << pos)) != (0x01 << pos);\n};\nconst pcf8574_1 = flow.get('PCF8574_1');\nvar oldPinValues = context.get('oldPinValues') || 0b11111111;\n\npcf8574_1.readPins();\n\nif (pcf8574_1.getPins() == oldPinValues)\n{\n    return null;\n}\n\nvar retVal = [];\nfor (var i = 0; i < 8; i++)\n{\n    var curPinVal = getBitValue(pcf8574_1.getPins(), i);\n    if (curPinVal != getBitValue(oldPinValues, i))\n    {\n        retVal.push({\n            payload: {\n                value: curPinVal\n            }\n        });\n    }\n    else\n    {\n        retVal.push(null);\n    }\n}\n\ncontext.set('oldPinValues', pcf8574_1.getPins());\n\nreturn retVal;",
	"outputs": 8,
	"noerr": 0,
	"x": 350,
	"y": 420,
	"wires": [
		["3e07059f.03aa4a"],
		["73175301.dd862c"],
		["c5f1b83d.2bad38"],
		["8ee6a0ea.970dc"],
		[],
		[],
		[],
		[]
	]
}, {
	"id": "3e07059f.03aa4a",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "map1",
	"func": "return {\n    payload: {\n        src: 'buttonUpSR',\n        value: msg.payload.value\n    }\n};",
	"outputs": 1,
	"noerr": 0,
	"x": 610,
	"y": 380,
	"wires": [
		["ba8d17d5.61c578"]
	]
}, {
	"id": "49a45935.e5b3a8",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "remap1",
	"func": "switch(msg.payload.cmd)\n{\n    case 'up':\n        return {\n            payload: [{\n                pin: 0,\n                value: true\n            },\n            {\n                pin: 1,\n                value: false\n            }]\n        }\n    case 'down':\n        return {\n            payload: [{\n                pin: 0,\n                value: false\n            },\n            {\n                pin: 1,\n                value: true\n            }]\n        }\n    case 'stop':\n        return {\n            payload: [{\n                pin: 0,\n                value: false\n            },\n            {\n                pin: 1,\n                value: false\n            }]\n        }\n    default:\n        break;\n}\nreturn null;",
	"outputs": 1,
	"noerr": 0,
	"x": 1020,
	"y": 420,
	"wires": [
		["a71eba91.eee258"]
	]
}, {
	"id": "ba8d17d5.61c578",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "BlindSR",
	"func": "var blindSR = flow.get('blindSR');\n\nswitch (msg.payload.src)\n{\n    case 'buttonUpSR':\n        if (msg.payload.value) {\n            return blindSR.button2();\n        }\n        break;\n    case 'buttonDownSR':\n        if (msg.payload.value) {\n            return blindSR.button1();\n        }\n        break;\n    default:\n        return null;\n}\n\nreturn null;",
	"outputs": 1,
	"noerr": 0,
	"x": 860,
	"y": 420,
	"wires": [
		["49a45935.e5b3a8"]
	]
}, {
	"id": "a71eba91.eee258",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "WriteMAX7311Outputs",
	"func": "var max7311_1 = flow.get('MAX7311_1');\n\nmsg.payload.forEach(function(item, index)\n{\n   max7311_1.writePin(item.pin, item.value);\n});\n\nmax7311_1.commit();\n\nreturn msg;",
	"outputs": 1,
	"noerr": 0,
	"x": 1280,
	"y": 420,
	"wires": [
		[]
	]
}, {
	"id": "73175301.dd862c",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "map2",
	"func": "return {\n    payload: {\n        src: 'buttonDownSR',\n        value: msg.payload.value\n    }\n};",
	"outputs": 1,
	"noerr": 0,
	"x": 610,
	"y": 420,
	"wires": [
		["ba8d17d5.61c578"]
	]
}, {
	"id": "c5f1b83d.2bad38",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "map3",
	"func": "return {\n    payload: {\n        src: 'buttonUpKiZi2',\n        value: msg.payload.value\n    }\n};",
	"outputs": 1,
	"noerr": 0,
	"x": 610,
	"y": 460,
	"wires": [
		["647b7ca5.24cb44"]
	]
}, {
	"id": "8ee6a0ea.970dc",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "map4",
	"func": "return {\n    payload: {\n        src: 'buttonDownKiZi2',\n        value: msg.payload.value\n    }\n};",
	"outputs": 1,
	"noerr": 0,
	"x": 610,
	"y": 500,
	"wires": [
		["647b7ca5.24cb44"]
	]
}, {
	"id": "647b7ca5.24cb44",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "BlindKiZi2",
	"func": "var blindKiZi2 = flow.get('blindKiZi2');\n\nswitch (msg.payload.src)\n{\n    case 'buttonUpKiZi2':\n        if (msg.payload.value) {\n            return blindKiZi2.button2();\n        }\n        break;\n    case 'buttonDownKiZi2':\n        if (msg.payload.value) {\n            return blindKiZi2.button1();\n        }\n        break;\n    default:\n        return null;\n}\n\nreturn null;",
	"outputs": 1,
	"noerr": 0,
	"x": 860,
	"y": 460,
	"wires": [
		["57632c8a.f20564"]
	]
}, {
	"id": "57632c8a.f20564",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "remap2",
	"func": "switch(msg.payload.cmd)\n{\n    case 'up':\n        return {\n            payload: [{\n                pin: 2,\n                value: true\n            },\n            {\n                pin: 3,\n                value: false\n            }]\n        }\n    case 'down':\n        return {\n            payload: [{\n                pin: 2,\n                value: false\n            },\n            {\n                pin: 3,\n                value: true\n            }]\n        }\n    case 'stop':\n        return {\n            payload: [{\n                pin: 2,\n                value: false\n            },\n            {\n                pin: 3,\n                value: false\n            }]\n        }\n    default:\n        break;\n}\nreturn null;",
	"outputs": 1,
	"noerr": 0,
	"x": 1020,
	"y": 460,
	"wires": [
		["a71eba91.eee258"]
	]
}, {
	"id": "435c84b6.64d6a4",
	"type": "mqtt in",
	"z": "f368fc00.a992a",
	"name": "",
	"topic": "/homeautomation/blinds/SR",
	"qos": "2",
	"datatype": "utf8",
	"broker": "7a6b0a36.7a8694",
	"x": 200,
	"y": 160,
	"wires": [
		["5d071f24.19021", "c5c000da.d460e8"]
	]
}, {
	"id": "5d071f24.19021",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "",
	"func": "var src = \"button\";\nif (msg.payload == \"up\")\n{\n    src += \"Up\";\n}\nelse if (msg.payload == \"down\")\n{\n    src += \"Down\";\n}\nelse\n{\n    return null;\n}\nsrc += \"SR\";\n\nreturn {\n    \"payload\": {\n    \"src\": src,\n    \"value\": true\n  }\n};",
	"outputs": 1,
	"noerr": 0,
	"initialize": "",
	"finalize": "",
	"x": 700,
	"y": 160,
	"wires": [
		["52151b6.e537ae4", "ba8d17d5.61c578"]
	]
}, {
	"id": "bce0e55a.ad844",
	"type": "mqtt in",
	"z": "f368fc00.a992a",
	"name": "",
	"topic": "/homeautomation/blinds/KiZi2",
	"qos": "2",
	"datatype": "utf8",
	"broker": "7a6b0a36.7a8694",
	"x": 180,
	"y": 680,
	"wires": [
		["48f37efb.83dff8"]
	]
}, {
	"id": "48f37efb.83dff8",
	"type": "function",
	"z": "f368fc00.a992a",
	"name": "",
	"func": "var src = \"button\";\nif (msg.payload == \"up\")\n{\n    src += \"Up\";\n}\nelse if (msg.payload == \"down\")\n{\n    src += \"Down\";\n}\nelse\n{\n    return null;\n}\nsrc += \"KiZi2\";\n\nreturn {\n    \"payload\": {\n    \"src\": src,\n    \"value\": true\n  }\n};",
	"outputs": 1,
	"noerr": 0,
	"initialize": "",
	"finalize": "",
	"x": 540,
	"y": 680,
	"wires": [
		["647b7ca5.24cb44"]
	]
}, {
	"id": "52151b6.e537ae4",
	"type": "debug",
	"z": "f368fc00.a992a",
	"name": "",
	"active": false,
	"tosidebar": true,
	"console": false,
	"tostatus": false,
	"complete": "false",
	"statusVal": "",
	"statusType": "auto",
	"x": 890,
	"y": 220,
	"wires": []
}, {
	"id": "c5c000da.d460e8",
	"type": "debug",
	"z": "f368fc00.a992a",
	"name": "",
	"active": false,
	"tosidebar": true,
	"console": false,
	"tostatus": false,
	"complete": "false",
	"statusVal": "",
	"statusType": "auto",
	"x": 520,
	"y": 120,
	"wires": []
}, {
	"id": "440313ac.6f2e7c",
	"type": "ui_group",
	"name": "Raffstore Schrankraum",
	"tab": "e2ec0056.67b71",
	"disp": true,
	"width": 5,
	"collapse": false
}, {
	"id": "7c3d9abe.f8c944",
	"type": "ui_group",
	"name": "Raffstore KiZi2",
	"tab": "e2ec0056.67b71",
	"order": 2,
	"disp": true,
	"width": 5,
	"collapse": false
}, {
	"id": "7a6b0a36.7a8694",
	"type": "mqtt-broker",
	"name": "somehost.local",
	"broker": "somehost.local",
	"port": "1883",
	"clientid": "",
	"usetls": false,
	"compatmode": false,
	"keepalive": "60",
	"cleansession": true,
	"birthTopic": "",
	"birthQos": "0",
	"birthPayload": "",
	"closeTopic": "",
	"closeQos": "0",
	"closePayload": "",
	"willTopic": "",
	"willQos": "0",
	"willPayload": ""
}, {
	"id": "e2ec0056.67b71",
	"type": "ui_tab",
	"name": "HomeAutomation",
	"icon": "dashboard",
	"order": 1,
	"disabled": false,
	"hidden": false
}]
`
