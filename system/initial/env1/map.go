package env1

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/null"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

func addMaps(adaptors *adaptors.Adaptors,
	scripts map[string]*m.Script,
	devices []*m.Device,
	imageList map[string]*m.Image,
	deviceActions map[string]*m.DeviceAction,
	deviceStates map[string]*m.DeviceState) (maps []*m.Map) {

	var err error

	// zones
	// ------------------------------------------------
	mainHallZone := &m.MapZone{
		Name: "Main Hall",
	}
	mainHallZone.Id, err = adaptors.MapZone.Add(mainHallZone)
	So(err, ShouldBeNil)

	kitchenZone := &m.MapZone{
		Name: "Kitchen",
	}
	kitchenZone.Id, err = adaptors.MapZone.Add(kitchenZone)
	So(err, ShouldBeNil)

	// map 1
	// ------------------------------------------------
	map1 := &m.Map{
		Name:        "office1",
		Description: "офис на ул. Красный проспект, д.22",
		Options: m.MapOptions{
			Zoom:              1,
			ElementStateText:  false,
			ElementOptionText: false,
		},
	}
	ok, _ := map1.Valid()
	So(ok, ShouldEqual, true)
	map1.Id, err = adaptors.Map.Add(map1)
	So(err, ShouldBeNil)

	// background layer
	// ------------------------------------------------
	backgroundLayer := &m.MapLayer{
		Name:        "background",
		Status:      "enabled",
		Description: "фон",
		MapId:       map1.Id,
		Weight:      1,
	}
	ok, _ = backgroundLayer.Valid()
	So(ok, ShouldEqual, true)
	backgroundLayer.Id, err = adaptors.MapLayer.Add(backgroundLayer)
	So(err, ShouldBeNil)

	// background image
	// ------------------------------------------------
	backgroundImage1 := &m.MapImage{
		ImageId: imageList["map-schematic-original"].Id,
	}

	ok, _ = backgroundImage1.Valid()
	So(ok, ShouldEqual, true)

	mapElementBackgroundImage1 := &m.MapElement{
		Name: "background schematic map",
		Prototype: m.Prototype{
			MapImage: backgroundImage1,
		},
		MapId:   map1.Id,
		LayerId: backgroundLayer.Id,
		Status:  Frozen,
		GraphSettings: m.MapElementGraphSettings{
			Position: m.MapElementGraphSettingsPosition{
				Top:  0,
				Left: 0,
			},
		},
	}
	ok, _ = mapElementBackgroundImage1.Valid()
	So(ok, ShouldEqual, true)
	mapElementBackgroundImage1.Id, err = adaptors.MapElement.Add(mapElementBackgroundImage1)
	So(err, ShouldBeNil)

	// base layer
	// ------------------------------------------------
	baseLayer := &m.MapLayer{
		Name:        "base",
		Status:      "enabled",
		Description: "базовый слой",
		MapId:       map1.Id,
		Weight:      0,
	}
	ok, _ = baseLayer.Valid()
	So(ok, ShouldEqual, true)
	baseLayer.Id, err = adaptors.MapLayer.Add(baseLayer)
	So(err, ShouldBeNil)

	// light1
	// ------------------------------------------------
	devLight1 := &m.MapDevice{
		SystemName: "DEV1_LIGHT1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light1_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light1_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light1_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light1_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}

	ok, _ = devLight1.Valid()
	mapElementLight1 := &m.MapElement{
		Name: "dev1_light1",
		Description: "Light1 in the hall",
		Prototype: m.Prototype{
			MapDevice: devLight1,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  644,
				Left: 329,
			},
		},
		Zone: mainHallZone,
	}
	ok, _ = mapElementLight1.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight1.Id, err = adaptors.MapElement.Add(mapElementLight1)
	So(err, ShouldBeNil)

	// light2
	// ------------------------------------------------
	devLight2 := &m.MapDevice{
		SystemName: "DEV1_LIGHT2",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light2_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light2_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light2_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light2_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight2 := &m.MapElement{
		Name: "dev1_light2",
		Description: "Light2 in the hall",
		Prototype: m.Prototype{
			MapDevice: devLight2,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  487,
				Left: 541,
			},
		},
		Zone: mainHallZone,
	}
	ok, _ = mapElementLight2.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight2.Id, err = adaptors.MapElement.Add(mapElementLight2)
	So(err, ShouldBeNil)

	// light3
	// ------------------------------------------------
	devLight3 := &m.MapDevice{
		SystemName: "DEV1_LIGHT3",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light3_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light3_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light3_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light3_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight3 := &m.MapElement{
		Name: "dev1_light3",
		Description: "Light3 in the hall",
		Prototype: m.Prototype{
			MapDevice: devLight3,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  779,
				Left: 630,
			},
		},
		Zone: mainHallZone,
	}
	ok, _ = mapElementLight3.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight3.Id, err = adaptors.MapElement.Add(mapElementLight3)
	So(err, ShouldBeNil)

	// light4
	// ------------------------------------------------
	devLight4 := &m.MapDevice{
		SystemName: "DEV1_LIGHT4",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["lamp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_light4_on"].Id,
				ImageId:       imageList["lamp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_light4_off"].Id,
				ImageId:       imageList["lamp_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["lamp_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_light4_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_light4_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementLight4 := &m.MapElement{
		Name: "dev1_light4",
		Description: "Light in the kitchen",
		Prototype: m.Prototype{
			MapDevice: devLight4,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  191,
				Left: 564,
			},
		},
		Zone: kitchenZone,
	}
	ok, _ = mapElementLight4.Valid()
	So(ok, ShouldEqual, true)
	mapElementLight4.Id, err = adaptors.MapElement.Add(mapElementLight4)
	So(err, ShouldBeNil)

	// fan5
	// ------------------------------------------------
	devFan1 := &m.MapDevice{
		SystemName: "DEV1_FAN1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["fan_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_fan1_on"].Id,
				ImageId:       imageList["fan_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_fan1_off"].Id,
				ImageId:       imageList["fan_v1_def"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_error"].Id,
				ImageId:       imageList["fan_v1_r"].Id,
			},
		},
		Actions: []*m.MapDeviceAction{
			{
				DeviceActionId: deviceActions["mb_dev1_turn_on_fan1_v1"].Id,
				ImageId:        imageList["button_v1_on"].Id,
			},
			{
				DeviceActionId: deviceActions["mb_dev1_turn_off_fan1_v1"].Id,
				ImageId:        imageList["button_v1_off"].Id,
			},
		},
	}
	mapElementFan1 := &m.MapElement{
		Name: "dev1_fan1",
		Description: "fan in the kitchen",
		Prototype: m.Prototype{
			MapDevice: devFan1,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  113,
				Left: 734,
			},
		},
		Zone: kitchenZone,
	}
	ok, _ = mapElementFan1.Valid()
	So(ok, ShouldEqual, true)
	mapElementFan1.Id, err = adaptors.MapElement.Add(mapElementFan1)
	So(err, ShouldBeNil)

	// temp1
	// ------------------------------------------------
	dev1Temp1 := &m.MapDevice{
		SystemName: "DEV1_TEMP1",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["temp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_temp1_on"].Id,
				ImageId:       imageList["temp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_temp1_off"].Id,
				ImageId:       imageList["temp_v1_r"].Id,
			},
		},
	}
	mapElementTemp1 := &m.MapElement{
		Name:        "dev1_temp1",
		Description: "temp sensor in the kitchen",
		Prototype: m.Prototype{
			MapDevice: dev1Temp1,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  388,
				Left: 288,
			},
		},
		Zone: kitchenZone,
	}
	ok, _ = mapElementTemp1.Valid()
	So(ok, ShouldEqual, true)
	mapElementTemp1.Id, err = adaptors.MapElement.Add(mapElementTemp1)
	So(err, ShouldBeNil)

	// temp2
	// ------------------------------------------------
	dev1Temp2 := &m.MapDevice{
		SystemName: "DEV1_TEMP2",
		DeviceId:   devices[0].Id,
		ImageId:    imageList["temp_v1_def"].Id,
		States: []*m.MapDeviceState{
			{
				DeviceStateId: deviceStates["dev1_temp2_on"].Id,
				ImageId:       imageList["temp_v1_y"].Id,
			},
			{
				DeviceStateId: deviceStates["dev1_temp2_off"].Id,
				ImageId:       imageList["temp_v1_r"].Id,
			},
		},
	}
	mapElementTemp2 := &m.MapElement{
		Name:        "dev1_temp2",
		Description: "temp sensor in the hall",
		Prototype: m.Prototype{
			MapDevice: dev1Temp2,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  113,
				Left: 468,
			},
		},
		Zone: mainHallZone,
	}
	ok, _ = mapElementTemp2.Valid()
	So(ok, ShouldEqual, true)
	mapElementTemp2.Id, err = adaptors.MapElement.Add(mapElementTemp2)
	So(err, ShouldBeNil)

	// map element text1
	// ------------------------------------------------
	mapText1 := &m.MapText{
		Text: "workflow:",
	}

	mapElementText1 := &m.MapElement{
		Name: "text1",
		Prototype: m.Prototype{
			MapText: mapText1,
		},
		MapId:   map1.Id,
		LayerId: baseLayer.Id,
		Status:  Enabled,
		GraphSettings: m.MapElementGraphSettings{
			Width:  null.NewInt64(33),
			Height: null.NewInt64(33),
			Position: m.MapElementGraphSettingsPosition{
				Top:  30,
				Left: 30,
			},
		},
	}

	ok, _ = mapElementText1.Valid()
	So(ok, ShouldEqual, true)
	mapElementText1.Id, err = adaptors.MapElement.Add(mapElementText1)
	So(err, ShouldBeNil)

	return
}

// ('dev1_light1', '', '1', 'device', '{"width": 33, "height": 33, "position": {"top": 644, "left": 329}}')
// ('dev1_light2', '', '2', 'device', '{"width": 33, "height": 33, "position": {"top": 487, "left": 541}}')
// ('dev1_light3', '', '3', 'device', '{"width": 33, "height": 33, "position": {"top": 779, "left": 630}}')
// ('dev1_light4', '', '4', 'device', '{"width": 33, "height": 33, "position": {"top": 191, "left": 564}}')