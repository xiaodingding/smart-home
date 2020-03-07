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

package controllers

import (
	"bytes"
	"encoding/json"
	dashboardModel "github.com/e154/smart-home/api/websocket/controllers/dashboard_models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
	"sync"
	"time"
)

type ControllerDashboard struct {
	*ControllerCommon
	Nodes    *dashboardModel.Nodes
	Devices  *dashboardModel.Devices
	Workflow *dashboardModel.Workflow
	Gate     *dashboardModel.Gate
	Cpu      *dashboardModel.Cpu
	sendLock *sync.Mutex
	buf      *bytes.Buffer
	enc      *json.Encoder
}

func NewControllerDashboard(common *ControllerCommon) (dashboard *ControllerDashboard) {
	dashboard = &ControllerDashboard{
		ControllerCommon: common,
		Nodes:            dashboardModel.NewNode(common.metric),
		Devices:          dashboardModel.NewDevices(common.metric),
		Workflow:         dashboardModel.NewWorkflow(common.metric),
		Gate:             dashboardModel.NewGate(common.metric),
		Cpu:              dashboardModel.NewCpu(common.metric),
		buf:              bytes.NewBuffer(nil),
		sendLock:         &sync.Mutex{},
	}
	dashboard.enc = json.NewEncoder(dashboard.buf)
	return dashboard
}

func (c *ControllerDashboard) Start() {
	c.metric.Subscribe("dashboard", c)
	c.stream.Subscribe("dashboard.get.nodes.status", c.Nodes.NodesStatus)
	c.stream.Subscribe("dashboard.get.gate.status", c.Gate.Status)
	//c.stream.Subscribe("t.get.flows.status", dashboardModel.FlowsStatus)
	c.stream.Subscribe("dashboard.get.telemetry", c.Telemetry)
}

func (c *ControllerDashboard) Stop() {
	c.metric.UnSubscribe("dashboard")
	c.stream.UnSubscribe("dashboard.get.nodes.status")
	c.stream.UnSubscribe("dashboard.get.gate.status")
	//c.stream.UnSubscribe("t.get.flows.status")
	c.stream.UnSubscribe("dashboard.get.telemetry")
}

func (t *ControllerDashboard) Broadcast(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch v := param.(type) {
	case string:
		switch v {
		case "workflow":
			body, ok = t.Workflow.Broadcast()
		case "node":
			body, ok = t.Nodes.Broadcast()
		case "device":
			body, ok = t.Devices.Broadcast()
		case "gate":
			body, ok = t.Gate.Broadcast()
		case "cpu":
			body, ok = t.Cpu.Broadcast()

		}
	case metrics.MapElementCursor:

	}

	if ok {
		go t.sendMsg(body)
	}
}

func (t *ControllerDashboard) sendMsg(payload map[string]interface{}) (err error) {

	t.sendLock.Lock()
	defer t.sendLock.Unlock()

	msg := stream.Message{
		Command: "dashboard.telemetry",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: payload,
	}

	t.buf.Reset()
	if err = t.enc.Encode(msg); err != nil {
		return
	}

	data := make([]byte, t.buf.Len())
	copy(data, t.buf.Bytes())
	t.stream.Broadcast(data)

	return
}

// only on request: 'dashboard.get.telemetry'
//
func (t *ControllerDashboard) Telemetry(client stream.IStreamClient, message stream.Message) {

	msg := stream.Message{
		Id:      message.Id,
		Command: "dashboard.telemetry",
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"memory":  t.metric.Memory.Snapshot(),
			"cpu":     t.metric.Cpu.Snapshot(),
			"time":    time.Now(),
			"uptime":  t.metric.Uptime.Snapshot(),
			"disk":    t.metric.Disk.Snapshot(),
			"nodes":   t.metric.Node.Snapshot(),
			"devices": t.metric.Device.Snapshot(),
			"gate":    t.metric.Gate.Snapshot(),
		},
	}

	client.Write(msg.Pack())
}
