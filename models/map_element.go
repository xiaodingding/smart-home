package models

import (
	"encoding/json"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/null"
	"github.com/e154/smart-home/system/validation"
	"time"
)

type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}
type MapElementGraphSettings struct {
	Width    null.Int64                      `json:"width"`
	Height   null.Int64                      `json:"height"`
	Position MapElementGraphSettingsPosition `json:"position"`
}

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
	case n.MapDevice != nil:
		b, err = json.Marshal(n.MapDevice)
	default:
		b, err = json.Marshal(struct{}{})
		return
	}
	return
}

func (n *Prototype) UnmarshalJSON(data []byte) (err error) {

	device := &MapDevice{}
	err = json.Unmarshal(data, device)
	if device.Device != nil && device.Device.Id != 0 && device.DeviceId != 0 {
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

type MapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name" valid:"Required"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType PrototypeType           `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype" valid:"Required"`
	MapId         int64                   `json:"map_id" valid:"Required"`
	LayerId       int64                   `json:"layer_id" valid:"Required"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        StatusType              `json:"status" valid:"Required"`
	Weight        int64                   `json:"weight"`
	Zone          *ZoneTag                `json:"zone"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

func (m *MapElement) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

type SortMapElementByWeight []*MapElement

func (l SortMapElementByWeight) Len() int           { return len(l) }
func (l SortMapElementByWeight) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SortMapElementByWeight) Less(i, j int) bool { return l[i].Weight < l[j].Weight }

type SortMapElement struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}