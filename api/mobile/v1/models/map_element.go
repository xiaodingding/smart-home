package models

import (
	"encoding/json"
	"github.com/e154/smart-home-old/api/log"
)

type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}
type MapElementGraphSettings struct {
	Width    *int64                          `json:"width"`
	Height   *int64                          `json:"height"`
	Position MapElementGraphSettingsPosition `json:"position"`
}

// swagger:model
type Prototype struct {
	*MapImage
	*MapText
	*MapDevice
}

func (n Prototype) MarshalJSON() (b []byte, err error) {

	switch {
	case n.MapText != nil:
		b, err = json.Marshal(n.MapText)
	case n.MapImage != nil:
		b, err = json.Marshal(n.MapImage)
	case n.MapDevice != nil && n.MapDevice.Device != nil:
		b, err = json.Marshal(n.MapDevice)
	case n.MapDevice != nil && n.MapDevice.DeviceId != 0:
		b, err = json.Marshal(n.MapDevice)
	default:
		b, err = json.Marshal(struct{}{})
		return
	}
	return
}

func (n *Prototype) UnmarshalJSON(data []byte) (err error) {

	device := &MapDevice{}
	if err = json.Unmarshal(data, device); err != nil {
		log.Error(err)
	}

	if (device.Device != nil && device.Device.Id != 0 && device.DeviceId != 0) ||
		device.DeviceId != 0 {
		n.MapDevice = device
		return
	}

	image := &MapImage{}
	err = json.Unmarshal(data, image)
	if image.ImageId != 0 {
		n.MapImage = image
		return
	}

	text := &MapText{}
	err = json.Unmarshal(data, text)
	n.MapText = text
	return
}

// swagger:model
type MapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
	Zone          *ZoneTag                `json:"zone,omitempty"`
}

type SortMapElement struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}