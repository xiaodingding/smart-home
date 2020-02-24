// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package zigbee2mqtt

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("zigbee2mqtt")
)

type Zigbee2mqtt struct {
	graceful    *graceful_service.GracefulService
	mqtt        *mqtt.Mqtt
	adaptors    *adaptors.Adaptors
	isStarted   bool
	bridgesLock *sync.Mutex
	bridges     map[int64]*Bridge
}

func NewZigbee2mqtt(graceful *graceful_service.GracefulService,
	mqtt *mqtt.Mqtt,
	adaptors *adaptors.Adaptors) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		graceful:    graceful,
		mqtt:        mqtt,
		adaptors:    adaptors,
		bridgesLock: &sync.Mutex{},
		bridges:     make(map[int64]*Bridge),
	}
}

func (z *Zigbee2mqtt) Start() {
	if z.isStarted {
		return
	}
	z.isStarted = true

	models, _, err := z.adaptors.Zigbee2mqtt.List(99, 0)
	if err != nil {
		log.Error(err.Error())
	}

	if len(models) == 0 {
		model := &m.Zigbee2mqtt{
			Name:       "zigbee2mqtt",
			BaseTopic:  "zigbee2mqtt",
			PermitJoin: true,
		}
		model.Id, err = z.adaptors.Zigbee2mqtt.Add(model)
		if err != nil {
			log.Error(err.Error())
			return
		}
		models = append(models, model)
	}

	for _, model := range models {
		bridge := NewBridge(z.mqtt, z.adaptors, model)
		bridge.Start()
		z.bridges[model.Id] = bridge
	}
}

func (z *Zigbee2mqtt) Shutdown() {
	if !z.isStarted {
		return
	}
	z.isStarted = false
	for _, bridge := range z.bridges {
		bridge.Stop(context.Background())
	}
}

func (z *Zigbee2mqtt) AddBridge(model *m.Zigbee2mqtt) (err error) {

	model.Id, err = z.adaptors.Zigbee2mqtt.Add(model)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if model, err = z.adaptors.Zigbee2mqtt.GetById(model.Id); err != nil {
		return
	}

	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	bridge := NewBridge(z.mqtt, z.adaptors, model)
	bridge.Start()
	z.bridges[model.Id] = bridge
	return
}

func (z *Zigbee2mqtt) GetBridgeById(id int64) (*m.Zigbee2mqtt, error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	if br, ok := z.bridges[id]; ok {
		return br.model, nil
	}
	return nil, adaptors.ErrRecordNotFound
}

func (z *Zigbee2mqtt) ListBridges(limit, offset int64, order, sortBy string) (models []m.Zigbee2mqtt, total int64, err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	total = int64(len(z.bridges))

	for _, br := range z.bridges {
		models = append(models, *br.model)
	}

	return
}

func (z *Zigbee2mqtt) UpdateBridge(model *m.Zigbee2mqtt) (result *m.Zigbee2mqtt, err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	if br, ok := z.bridges[model.Id]; ok {
		br.UpdateModel(model)
	} else {
		err = adaptors.ErrRecordNotFound
		return
	}

	if err = z.adaptors.Zigbee2mqtt.Update(model); err != nil {
		return
	}

	result, err = z.adaptors.Zigbee2mqtt.GetById(model.Id)

	return
}

func (z *Zigbee2mqtt) DeleteBridge(id int64) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	if br, ok := z.bridges[id]; ok {
		br.Stop(context.Background())
		delete(z.bridges, id)
	} else {
		err = adaptors.ErrRecordNotFound
		return
	}

	err = z.adaptors.Zigbee2mqtt.Delete(id)

	return
}