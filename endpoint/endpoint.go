package endpoint

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("endpoint")
)

type Endpoint struct {
	Auth             *AuthEndpoint
	Device           *DeviceEndpoint
	DeviceAction     *DeviceActionEndpoint
	DeviceState      *DeviceStateEndpoint
	Flow             *FlowEndpoint
	Image            *ImageEndpoint
	Log              *LogEndpoint
	Map              *MapEndpoint
	MapElement       *MapElementEndpoint
	MapLayer         *MapLayerEndpoint
	Node             *NodeEndpoint
	Role             *RoleEndpoint
	Script           *ScriptEndpoint
	Workflow         *WorkflowEndpoint
	WorkflowScenario *WorkflowScenarioEndpoint
	User             *UserEndpoint
	Gate             *GateEndpoint
}

func NewEndpoint(adaptors *adaptors.Adaptors,
	core *core.Core,
	scriptService *scripts.ScriptService,
	accessList *access_list.AccessListService,
	gate *gate_client.GateClient) *Endpoint {
	common := NewCommonEndpoint(adaptors, core, accessList, scriptService, gate)
	return &Endpoint{
		Auth:             NewAuthEndpoint(common),
		Device:           NewDeviceEndpoint(common),
		DeviceAction:     NewDeviceActionEndpoint(common),
		DeviceState:      NewDeviceStateEndpoint(common),
		Flow:             NewFlowEndpoint(common),
		Image:            NewImageEndpoint(common),
		Log:              NewLogEndpoint(common),
		Map:              NewMapEndpoint(common),
		MapElement:       NewMapElementEndpoint(common),
		MapLayer:         NewMapLayerEndpoint(common),
		Node:             NewNodeEndpoint(common),
		Role:             NewRoleEndpoint(common),
		Script:           NewScriptEndpoint(common),
		Workflow:         NewWorkflowEndpoint(common),
		WorkflowScenario: NewWorkflowScenarioEndpoint(common),
		User:             NewUserEndpoint(common),
		Gate:             NewGateEndpoint(common),
	}
}