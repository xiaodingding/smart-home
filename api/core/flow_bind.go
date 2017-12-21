package core

import "github.com/e154/smart-home/api/models"

// Javascript Binding
//
// IC.Flow()
//	 .getName()
//	 .getDescription()
//	 .setVar(string, interface)
//	 .getVar(string)
//	 .node()
//
type FlowBind struct {
	flow *Flow
}

func (f *FlowBind) GetName() string {
	return f.flow.Model.Name
}

func (f *FlowBind) GetDescription() string {
	return f.flow.Model.Description
}

func (f *FlowBind) SetVar(key string, value interface{}) {
	f.flow.SetVar(key, value)
}

func (f *FlowBind) GetVar(key string) interface{} {
	return f.flow.GetVar(key)
}

func (f *FlowBind) Node() *models.Node {
	return f.flow.Node
}