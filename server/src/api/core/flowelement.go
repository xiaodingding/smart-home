package core

import (
	"sync"
	"../models"
	"../scripts"
)

func NewFlowElement(model *models.FlowElement, flow *Flow, workflow *Workflow) (flowElement *FlowElement, err error) {

	flowElement = &FlowElement{
		Model:model,
		Flow: flow,
		Workflow: workflow,
	}

	if model.Script == nil {
		return
	}

	var script *models.Script
	if script, err = models.GetScriptById(model.Script.Id); err != nil {
		return
	}

	if flowElement.Script, err = scripts.New(script); err != nil {
		return
	}

	return
}

type FlowElement struct {
	Model 		*models.FlowElement
	Flow		*Flow
	Workflow	*Workflow
	Script		*scripts.Engine
	Prototype	ActionPrototypes
	status		Status
	mutex     	sync.Mutex
}

func (m *FlowElement) Before(message *Message) error {

	m.mutex.Lock()
	m.status = DONE
	m.mutex.Unlock()
	return m.Prototype.Before(message, m.Flow)
}

// run internal process
func (m *FlowElement) Run(message *Message) (err error) {

	m.mutex.Lock()
	m.status = IN_PROCESS
	m.mutex.Unlock()

	//m.Flow.PushCursor(m)
	err = m.Before(message)
	err = m.Prototype.Run(message, m.Flow)

	//run script if exist
	var isTrue, isScripted bool
	if m.Script != nil {
		isScripted = true
		m.Script.PushStruct("message", message)
		res, _ := m.Script.Do()
		isTrue = res == "true"
	}

	err = m.After(message)

	// each connections
	for _, conn := range m.Flow.Connections {
		if conn.ElementFrom != m.Model.Uuid || conn.ElementTo == m.Model.Uuid {
			continue
		}

		for _, element := range m.Flow.FlowElements {
			if conn.ElementTo != element.Model.Uuid {
				continue
			}

			if isScripted {
				if conn.Direction == "true" {
					if !isTrue {
						continue
					}
				} else if conn.Direction == "false" {
					if isTrue {
						continue
					}
				}
			}

			go element.Run(message)
		}
	}

	//m.Flow.PopCursor(m)

	m.mutex.Lock()
	m.status = ENDED
	m.mutex.Unlock()

	return
}

func (m *FlowElement) After(message *Message) error {

	m.mutex.Lock()
	m.status = STARTED
	m.mutex.Unlock()

	return  m.Prototype.After(message, m.Flow)
}

func (m *FlowElement) GetStatus() (status Status) {

	m.mutex.Lock()
	status = m.status
	m.mutex.Unlock()

	return
}