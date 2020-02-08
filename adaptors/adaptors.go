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

package adaptors

import (
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/migrations"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("adaptors")
)

type Adaptors struct {
	Node             *Node
	Script           *Script
	Workflow         *Workflow
	WorkflowScenario *WorkflowScenario
	Device           *Device
	DeviceAction     *DeviceAction
	DeviceState      *DeviceState
	Flow             *Flow
	FlowElement      *FlowElement
	Connection       *Connection
	Worker           *Worker
	Role             *Role
	Permission       *Permission
	User             *User
	UserMeta         *UserMeta
	Image            *Image
	Variable         *Variable
	Map              *Map
	MapLayer         *MapLayer
	MapText          *MapText
	MapImage         *MapImage
	MapDevice        *MapDevice
	MapElement       *MapElement
	MapDeviceState   *MapDeviceState
	MapDeviceAction  *MapDeviceAction
	Log              *Log
	MapZone          *MapZone
	Template         *Template
	Message          *Message
	MessageDelivery  *MessageDelivery
	FlowSubscription *FlowSubscription
}

func NewAdaptors(db *gorm.DB,
	cfg *config.AppConfig,
	migrations *migrations.Migrations) (adaptors *Adaptors) {

	if cfg.AutoMigrate {
		migrations.Up()
	}

	adaptors = &Adaptors{
		Node:             GetNodeAdaptor(db),
		Script:           GetScriptAdaptor(db),
		Workflow:         GetWorkflowAdaptor(db),
		WorkflowScenario: GetWorkflowScenarioAdaptor(db),
		Device:           GetDeviceAdaptor(db),
		DeviceAction:     GetDeviceActionAdaptor(db),
		DeviceState:      GetDeviceStateAdaptor(db),
		Flow:             GetFlowAdaptor(db),
		FlowElement:      GetFlowElementAdaptor(db),
		Connection:       GetConnectionAdaptor(db),
		Worker:           GetWorkerAdaptor(db),
		Role:             GetRoleAdaptor(db),
		Permission:       GetPermissionAdaptor(db),
		User:             GetUserAdaptor(db),
		UserMeta:         GetUserMetaAdaptor(db),
		Image:            GetImageAdaptor(db),
		Variable:         GetVariableAdaptor(db),
		Map:              GetMapAdaptor(db),
		MapLayer:         GetMapLayerAdaptor(db),
		MapText:          GetMapTextAdaptor(db),
		MapImage:         GetMapImageAdaptor(db),
		MapDevice:        GetMapDeviceAdaptor(db),
		MapElement:       GetMapElementAdaptor(db),
		MapDeviceState:   GetMapDeviceStateAdaptor(db),
		MapDeviceAction:  GetMapDeviceActionAdaptor(db),
		Log:              GetLogAdaptor(db),
		MapZone:          GetMapZoneAdaptor(db),
		Template:         GetTemplateAdaptor(db),
		Message:          GetMessageAdaptor(db),
		MessageDelivery:  GetMessageDeliveryAdaptor(db),
		FlowSubscription: GetFlowSubscriptionAdaptor(db),
	}

	return
}
