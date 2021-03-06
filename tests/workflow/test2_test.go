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

package workflow

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

//
// create workflow
//
// add workflow scenarios (wf_scenario_1 + script7)
//
// add flow (flow1)
// +----------+    +----------+
// | handler  |    |  emitter |
// | script6  +---->  script5 |
// |          |    |          |
// +----------+    +----------+
//
// run core
//
func Test2(t *testing.T) {

	var store interface{}
	Convey("add scripts", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			scriptService.PushFunctions("store", func(value interface{}) {
				store = value
			})
		})
	})

	Convey("run core", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			// stop core
			// ------------------------------------------------
			err := c.Stop()
			So(err, ShouldBeNil)

			// clear database
			// ------------------------------------------------
			migrations.Purge()

			// create scripts
			// ------------------------------------------------
			script1 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test1",
				Source:      coffeeScript1,
				Description: "test1",
			}
			script5 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test5",
				Source:      coffeeScript5,
				Description: "test5",
			}
			script6 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test6",
				Source:      coffeeScript6,
				Description: "test6",
			}
			script7 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test7",
				Source:      coffeeScript7,
				Description: "test7",
			}

			ok, _ := script1.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script5.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script6.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script7.Valid()
			So(ok, ShouldEqual, true)

			engine1, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)
			err = engine1.Compile()
			So(err, ShouldBeNil)
			script1Id, err := adaptors.Script.Add(script1)
			So(err, ShouldBeNil)
			script1, err = adaptors.Script.GetById(script1Id)
			So(err, ShouldBeNil)

			engine5, err := scriptService.NewEngine(script5)
			So(err, ShouldBeNil)
			err = engine5.Compile()
			So(err, ShouldBeNil)
			script5Id, err := adaptors.Script.Add(script5)
			So(err, ShouldBeNil)
			script5, err = adaptors.Script.GetById(script5Id)
			So(err, ShouldBeNil)

			engine6, err := scriptService.NewEngine(script6)
			So(err, ShouldBeNil)
			err = engine6.Compile()
			So(err, ShouldBeNil)
			script6Id, err := adaptors.Script.Add(script6)
			So(err, ShouldBeNil)
			script6, err = adaptors.Script.GetById(script6Id)
			So(err, ShouldBeNil)

			engine7, err := scriptService.NewEngine(script7)
			So(err, ShouldBeNil)
			err = engine7.Compile()
			So(err, ShouldBeNil)
			script7Id, err := adaptors.Script.Add(script7)
			So(err, ShouldBeNil)
			script7, err = adaptors.Script.GetById(script7Id)
			So(err, ShouldBeNil)

			// create workflow
			// ------------------------------------------------
			workflow := &m.Workflow{
				Name:        "main workflow",
				Description: "main workflow desc",
				Status:      "enabled",
			}

			ok, _ = workflow.Valid()
			So(ok, ShouldEqual, true)

			wfId, err := adaptors.Workflow.Add(workflow)
			So(err, ShouldBeNil)
			workflow.Id = wfId

			// add workflow scenario
			// ------------------------------------------------
			wfScenario1 := &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			ok, _ = wfScenario1.Valid()
			So(ok, ShouldEqual, true)

			wfScenarioId1, err := adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)
			wfScenario1.Id = wfScenarioId1

			err = adaptors.WorkflowScenario.AddScript(wfScenario1, script7)
			So(err, ShouldBeNil)

			workflow.Scenario = wfScenario1
			err = adaptors.Workflow.Update(workflow)
			So(err, ShouldBeNil)

			flow1 := &m.Flow{
				Name:               "flow1",
				Status:             Enabled,
				WorkflowId:         workflow.Id,
				WorkflowScenarioId: wfScenario1.Id,
			}

			ok, _ = flow1.Valid()
			So(ok, ShouldEqual, true)

			flow1.Id, err = adaptors.Flow.Add(flow1)
			So(err, ShouldBeNil)

			feEmitter := &m.FlowElement{
				Name:          "emitter",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageEmitter,
				ScriptId:      &script5.Id,
			}

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script6.Id,
			}

			ok, _ = feEmitter.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feEmitter.Uuid, err = adaptors.FlowElement.Add(feEmitter)
			So(err, ShouldBeNil)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			connect := &m.Connection{
				Name:        "con1",
				ElementFrom: feHandler.Uuid,
				ElementTo:   feEmitter.Uuid,
				FlowId:      flow1.Id,
				PointFrom:   1,
				PointTo:     1,
			}

			ok, _ = connect.Valid()
			So(ok, ShouldEqual, true)

			connect.Uuid, err = adaptors.Connection.Add(connect)
			So(err, ShouldBeNil)

			err = c.Run()
			So(err, ShouldBeNil)

			workflowCore, err := c.GetWorkflow(workflow.Id)
			So(err, ShouldBeNil)

			flowCore, err := workflowCore.GetFLow(flow1.Id)
			So(err, ShouldBeNil)

			message := core.NewMessage()

			// create context
			ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
			ctx = context.WithValue(ctx, "msg", message)

			err = flowCore.NewMessage(ctx)
			So(err, ShouldBeNil)

			So(store, ShouldEqual, "b")

			err = c.Stop()
			So(err, ShouldBeNil)
		})
	})
}
