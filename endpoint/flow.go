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

package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/validation"
)

type FlowEndpoint struct {
	*CommonEndpoint
}

func NewFlowEndpoint(common *CommonEndpoint) *FlowEndpoint {
	return &FlowEndpoint{
		CommonEndpoint: common,
	}
}

func (f *FlowEndpoint) Add(params *m.Flow) (result *m.Flow, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = f.adaptors.Flow.Add(params); err != nil {
		return
	}

	if result, err = f.adaptors.Flow.GetById(id); err != nil {
		return
	}

	if err = f.core.AddFlow(result); err != nil {
		return
	}

	f.metric.Update(metrics.FlowAdd{TotalNum: 1})

	return
}

func (f *FlowEndpoint) GetById(id int64) (flow *m.Flow, err error) {

	flow, err = f.adaptors.Flow.GetById(id)

	return
}

func (f *FlowEndpoint) GetRedactor(flowId int64) (redactorFlow *m.RedactorFlow, err error) {

	var flow *m.Flow
	if flow, err = f.adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	redactorFlow, err = f.ExportToRedactor(flow)

	return
}

func (f *FlowEndpoint) GetList(limit, offset int64, order, sortBy string) (list []*m.Flow, total int64, err error) {

	list, total, err = f.adaptors.Flow.List(limit, offset, order, sortBy)

	return
}

func (f *FlowEndpoint) Search(query string, limit, offset int) (list []*m.Flow, total int64, err error) {

	if list, total, err = f.adaptors.Flow.Search(query, limit, offset); err != nil {
		return
	}

	return
}

func (f *FlowEndpoint) Update(params *m.Flow) (result *m.Flow, errs []*validation.Error, err error) {

	var flow *m.Flow
	if flow, err = f.adaptors.Flow.GetById(flow.Id); err != nil {
		return
	}

	_, errs = flow.Valid()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("%s %s", err.Key, err.Message)
		}
		return
	}

	if err = f.adaptors.Flow.Update(flow); err != nil {
		return
	}

	if flow, err = f.adaptors.Flow.GetById(flow.Id); err != nil {
		return
	}

	err = f.core.UpdateFlow(flow)

	return
}

func (f *FlowEndpoint) Delete(flowId int64) (err error) {

	var flow *m.Flow
	if flow, err = f.adaptors.Flow.GetById(flowId); err != nil {
		return
	}

	if err = f.core.RemoveFlow(flow); err != nil {
		return
	}

	if err = f.adaptors.Flow.Delete(flowId); err != nil {
		return
	}

	f.metric.Update(metrics.FlowDelete{TotalNum: 1})

	return
}

func (f *FlowEndpoint) UpdateRedactor(params *m.RedactorFlow) (result *m.RedactorFlow,
	errs []*validation.Error, err error) {

	//debug.Println(params)
	//fmt.Println("------")

	var oldFlow *m.Flow
	if oldFlow, err = f.adaptors.Flow.GetById(params.Id); err != nil {
		return
	}

	newFlow := &m.Flow{}
	if err = common.Copy(&newFlow, &params, common.JsonEngine); err != nil {
		return
	}
	if params.Scenario != nil {
		newFlow.WorkflowScenarioId = params.Scenario.Id
	}
	if params.Workflow != nil {
		newFlow.WorkflowId = params.Workflow.Id
	}

	_, errs = newFlow.Valid()
	if len(errs) > 0 {
		return
	}

	tx := f.adaptors.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Flow.Update(newFlow); err != nil {
		return
	}

	//update flow elements
	flowTodoRemove := make([]uuid.UUID, 0)
	for _, element := range oldFlow.FlowElements {
		exist := false
		for _, object := range params.Objects {
			if object.Id.String() == element.Uuid.String() {
				exist = true
				break
			}
		}

		if !exist {
			flowTodoRemove = append(flowTodoRemove, element.Uuid)
		}
	}

	if len(flowTodoRemove) > 0 {
		if err = tx.FlowElement.Delete(flowTodoRemove); err != nil {
			return
		}
	}

	// insert or update flow elements
	for _, element := range params.Objects {

		fl := &m.FlowElement{}
		common.Copy(&fl, &element)
		common.Copy(&fl.GraphSettings.Position, &element.Position)
		fl.Uuid.Scan(element.Id)
		fl.FlowId = newFlow.Id
		fl.Name = element.Title

		if element.FlowLink != nil && element.FlowLink.Id != 0 {
			fl.FlowLink = &element.FlowLink.Id
		}

		if element.Script != nil {
			fl.ScriptId = &element.Script.Id
		}

		switch element.Type.Name {
		case "event":
			if element.Type.Start != nil {
				fl.PrototypeType = common.FlowElementsPrototypeMessageHandler
			} else if element.Type.End != nil {
				fl.PrototypeType = common.FlowElementsPrototypeMessageEmitter
			}
		case "task":
			fl.PrototypeType = common.FlowElementsPrototypeTask
		case "gateway":
			fl.PrototypeType = common.FlowElementsPrototypeGateway
		case "flow":
			fl.PrototypeType = common.FlowElementsPrototypeFlow
		default:
			fl.PrototypeType = common.FlowElementsPrototypeDefault
		}

		_, errs = fl.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if err = tx.FlowElement.AddOrUpdateFlowElement(fl); err != nil {
			return
		}
	}

	// connectors
	connTodoRemove := make([]uuid.UUID, 0)
	for _, oldConn := range oldFlow.Connections {
		exist := false
		for _, newConn := range params.Connectors {
			if oldConn.Uuid.String() == newConn.Id.String() {
				exist = true
				break
			}
		}

		if !exist {
			connTodoRemove = append(connTodoRemove, oldConn.Uuid)
		}
	}

	if len(connTodoRemove) > 0 {
		if err = tx.Connection.Delete(connTodoRemove); err != nil {
			return
		}
	}

	for _, c := range params.Connectors {

		conn := &m.Connection{
			Name:      c.Title,
			PointFrom: c.Start.Point,
			PointTo:   c.End.Point,
			FlowId:    newFlow.Id,
			Direction: c.Direction,
		}
		conn.Uuid.Scan(c.Id)
		conn.ElementFrom.Scan(c.Start.Object)
		conn.ElementTo.Scan(c.End.Object)

		_, errs = conn.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if err = tx.Connection.AddOrUpdateConnection(conn); err != nil {
			return
		}
	}

	// workers
	workersTodoRemove := make([]*m.Worker, 0)
	for _, oldWorker := range oldFlow.Workers {
		exist := false
		for _, newWorker := range params.Workers {
			if newWorker.Id == oldWorker.Id {
				exist = true
				break
			}
		}

		if !exist {
			workersTodoRemove = append(workersTodoRemove, oldWorker)
		}
	}

	for _, worker := range workersTodoRemove {
		if err = f.core.RemoveWorker(worker); err == nil {
			if err = tx.Worker.Delete([]int64{worker.Id}); err != nil {
				return
			}
		}
	}

	for _, w := range params.Workers {
		worker := &m.Worker{}
		common.Copy(&worker, &w)
		worker.WorkflowId = newFlow.Workflow.Id
		worker.FlowId = newFlow.Id
		worker.DeviceActionId = w.DeviceAction.Id

		_, errs = worker.Valid()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%s %s", err.Key, err.Message)
			}
			return
		}

		if worker.Id == 0 {
			if _, err = tx.Worker.Add(worker); err != nil {
				return
			}
		} else {
			if err = tx.Worker.Update(worker); err != nil {
				return
			}
		}
	}

	// mqtt topic subscriptions
	subscriptionsTodoRemove := make([]int64, 0)
	for _, oldSub := range oldFlow.Subscriptions {
		exist := false
		for _, newSub := range params.Subscriptions {
			if newSub.Topic == oldSub.Topic {
				exist = true
			}
		}
		if !exist {
			subscriptionsTodoRemove = append(subscriptionsTodoRemove, oldSub.Id)
		}
	}

	if err := tx.FlowSubscription.Remove(subscriptionsTodoRemove); err != nil {
		log.Error(err.Error())
	}

	for _, sub := range params.Subscriptions {
		if sub.Id == 0 && sub.Topic != "" {
			flowSubscription := &m.FlowSubscription{
				Topic:  sub.Topic,
				FlowId: newFlow.Id,
			}
			if err = tx.FlowSubscription.Add(flowSubscription); err != nil {
				log.Error(err.Error())
			}
		}
	}

	// zigbee2mqtt subscriptions
	var exist bool
	zigbe2mqttDeviceTodoRemove := make([]string, 0)
	for _, oldDev := range oldFlow.Zigbee2mqttDevices {
		exist = false
		for _, newDev := range params.Zigbee2mqttDevices {
			if newDev.Id == oldDev.Id {
				exist = true
			}
		}
		if !exist {
			zigbe2mqttDeviceTodoRemove = append(zigbe2mqttDeviceTodoRemove, oldDev.Id)
		}
	}

	if err := tx.FlowZigbee2mqttDevice.Remove(oldFlow.Id, zigbe2mqttDeviceTodoRemove); err != nil {
		log.Error(err.Error())
	}

	for _, dev := range params.Zigbee2mqttDevices {
		exist = false
		for _, oldDev := range oldFlow.Zigbee2mqttDevices {
			if oldDev.Id == dev.Id {
				exist = true
			}
		}
		if !exist {
			flowZigbee2mqttDevice := &m.FlowZigbee2mqttDevice{
				FlowId:              newFlow.Id,
				Zigbee2mqttDeviceId: dev.Id,
			}
			if err = tx.FlowZigbee2mqttDevice.Add(flowZigbee2mqttDevice); err != nil {
				log.Error(err.Error())
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	// exit
	if newFlow, err = f.adaptors.Flow.GetById(params.Id); err != nil {
		return
	}

	if err = f.core.UpdateFlow(newFlow); err != nil {
		return
	}

	result, err = f.ExportToRedactor(newFlow)

	return
}

func (n *FlowEndpoint) ExportToRedactor(f *m.Flow) (redactorFlow *m.RedactorFlow, err error) {

	if f == nil {
		err = errors.New("ExportToRedactor: Nil point")
		return
	}

	var scenario *m.WorkflowScenario
	if scenario, err = n.adaptors.WorkflowScenario.GetById(f.WorkflowScenarioId); err != nil {
		return
	}

	redactorFlow = &m.RedactorFlow{
		Id:                 f.Id,
		Name:               f.Name,
		Status:             f.Status,
		Description:        f.Description,
		Workflow:           f.Workflow,
		Subscriptions:      f.Subscriptions,
		Zigbee2mqttDevices: f.Zigbee2mqttDevices,
		Scenario:           scenario,
		Workers:            make([]*m.Worker, 0),
		Objects:            make([]*m.RedactorObject, 0),
		Connectors:         make([]*m.RedactorConnector, 0),
		CreatedAt:          f.CreatedAt,
		UpdatedAt:          f.UpdatedAt,
	}

	// elements
	for _, el := range f.FlowElements {
		object := &m.RedactorObject{
			Id:            el.Uuid,
			Title:         el.Name,
			Description:   el.Description,
			PrototypeType: el.PrototypeType,
			Script:        el.Script,
		}

		if el.FlowLink != nil {
			var flow *m.Flow
			if flow, err = n.adaptors.Flow.GetById(*el.FlowLink); err != nil {
				return
			}
			object.FlowLink = flow
		}

		switch el.PrototypeType {
		case "MessageHandler":
			object.Type.Name = "event"
			object.Type.Start = map[int64]interface{}{0: &map[int64]interface{}{0: true}}
		case "MessageEmitter":
			object.Type.Name = "event"
			object.Type.End = map[string]interface{}{"simply": &map[string]interface{}{"top_level": true}}
		case "Task":
			object.Type.Name = "task"
		case "Flow":
			object.Type.Name = "flow"
		case "Gateway":
			object.Type.Name = "gateway"
			object.Type.Start = map[int64]interface{}{0: &map[int64]interface{}{0: true}}
		default:

		}

		common.Copy(&object.Position, &el.GraphSettings.Position)
		redactorFlow.Objects = append(redactorFlow.Objects, object)
	}

	// connections
	for _, con := range f.Connections {
		connector := &m.RedactorConnector{
			Id:        con.Uuid,
			FlowType:  "default",
			Title:     con.Name,
			Direction: con.Direction,
		}
		connector.Start.Object = con.ElementFrom
		connector.Start.Point = con.PointFrom

		connector.End.Object = con.ElementTo
		connector.End.Point = con.PointTo

		redactorFlow.Connectors = append(redactorFlow.Connectors, connector)
	}

	// workers
	common.Copy(&redactorFlow.Workers, &f.Workers)

	return
}
