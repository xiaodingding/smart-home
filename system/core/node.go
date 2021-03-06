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

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"sync"
	"time"
)

type Nodes []*Node

type Node struct {
	*m.Node
	sync.Mutex
	mqttClient *mqtt.Client
	stat       NodeStat
	quit       chan struct{}
	ch         map[int64]chan *NodeResponse
	metric     *metrics.MetricManager
}

func NewNode(model *m.Node,
	mqtt *mqtt.Mqtt,
	metric *metrics.MetricManager) *Node {

	node := &Node{
		Node: model,
		//connStatus: "disabled",
		stat: NodeStat{
			ConnStatus: "disabled",
			LastPing:   time.Now(),
		},
		ch:         make(map[int64]chan *NodeResponse, 0),
		quit:       make(chan struct{}),
		metric:     metric,
		mqttClient: mqtt.NewClient(fmt.Sprintf("node_%v", model.Name)),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				node.Lock()
				node.stat.IsConnected = time.Now().Sub(node.stat.LastPing).Seconds() < 2

				if node.Node.Status == "enabled" {
					if node.stat.IsConnected {
						node.stat.ConnStatus = "connected"
					} else {
						node.stat.ConnStatus = "wait"
					}
				} else {
					node.stat.ConnStatus = "disabled"
				}
				go node.metric.Update(metrics.NodeUpdateStatus{Id: node.Id, Status: node.stat.ConnStatus})
				node.Unlock()

			case _, ok := <-node.quit:
				if !ok {
					return
				}
				close(node.quit)
				return
			}
		}
	}()

	return node
}

func (n *Node) Remove() {
	n.Lock()
	defer n.Unlock()

	log.Infof("Remove node %v", n.Id)

	n.quit <- struct{}{}
}

func (n *Node) Send(device *m.Device, command []byte) (result NodeResponse, err error) {

	//log.Debugf("send device(%v) command(%v)", device.Id, command)

	// time metric
	startTime := time.Now()

	ch := make(chan *NodeResponse)
	n.addCh(device.Id, ch)
	defer n.delCh(device.Id)

	// send message to node
	msg := &NodeMessage{
		DeviceId:   device.Id,
		DeviceType: device.Type,
		Properties: device.Properties,
		Command:    command,
	}

	n.MqttPublish(msg)

	// wait response
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var done bool
	for ; ; {
		if done {
			break
		}
		select {
		case <-ticker.C:
			//log.Debugf("request timeout device(%d)", device.Id)
			err = errors.New("timeout")
			done = true
		case resp := <-ch:
			if resp == nil {
				return
			}
			if resp.DeviceId != device.Id {
				continue
			}

			// response from node
			result = NodeResponse{
				DeviceId:   resp.DeviceId,
				Status:     resp.Status,
				DeviceType: resp.DeviceType,
				Response:   resp.Response,
				Properties: resp.Properties,
			}
			done = true
		}
	}

	result.Time = time.Since(startTime).Seconds()

	return
}

func (n *Node) addCh(deviceId int64, ch chan *NodeResponse) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; ok {
		return
	}

	n.ch[deviceId] = ch
}

func (n *Node) delCh(deviceId int64) {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[deviceId]; !ok {
		return
	}

	close(n.ch[deviceId])
	delete(n.ch, deviceId)
}

func (n *Node) Connect() *Node {

	n.mqttClient.Subscribe(n.topic("resp"), n.onPublish)
	n.mqttClient.Subscribe(n.topic("ping"), n.ping)

	return n
}

func (n *Node) IsConnected() bool {
	n.Lock()
	defer n.Unlock()
	return n.stat.IsConnected
}

func (n *Node) onPublish(client *mqtt.Client, msg mqtt.Message) {

	resp := &NodeResponse{}
	if err := json.Unmarshal(msg.Payload, resp); err != nil {
		log.Error(err.Error())
		return
	}

	n.Lock()
	defer n.Unlock()
	if _, ok := n.ch[resp.DeviceId]; !ok {
		return
	}

	n.ch[resp.DeviceId] <- resp
}

func (n *Node) ping(client *mqtt.Client, msg mqtt.Message) {

	var stat NodeStatModel
	_ = json.Unmarshal(msg.Payload, &stat)

	n.Lock()

	//n.stat.Status = stat.Status //????
	n.stat.Thread = stat.Thread
	n.stat.Rps = stat.Rps
	n.stat.Min = stat.Min
	n.stat.Max = stat.Max
	n.stat.StartedAt = stat.StartedAt
	n.stat.LastPing = time.Now()

	n.Unlock()

	return
}

func (n *Node) MqttPublish(msg interface{}) {

	data, _ := json.Marshal(msg)

	n.Lock()
	defer n.Unlock()

	if err := n.mqttClient.Publish(n.topic("req"), data); err != nil {
		log.Error(err.Error())
		return
	}
}

func (n *Node) topic(r string) string {
	return fmt.Sprintf("/home/node/%s/%s", n.Node.Name, r)
}

func (n *Node) GetConnStatus() string {
	n.Lock()
	defer n.Unlock()
	return n.stat.ConnStatus
}

func (n *Node) GetStat() NodeStat {
	n.Lock()
	defer n.Unlock()
	return n.stat
}

func (n *Node) UpdateClientParams(params *m.Node) {
	n.Lock()
	n.Node = params
	n.Unlock()

	// unsubscribe all mqtt client
	if n.mqttClient != nil {
		n.mqttClient.UnsubscribeAll()
	}

	if n.Status != "disabled" {
		n.Connect()
	}

}
