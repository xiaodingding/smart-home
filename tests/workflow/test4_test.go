package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
)

// create node
//
// create device
//
// create device-actions (+script)
//
func Test4(t *testing.T) {

	Convey("add scripts", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			// clear database
			migrations.Purge()

			// add node
			node := &m.Node{
				Name: "node",
				Ip: "127.0.0.1",
				Port: 3001,
				Status: "enabled",
			}
			ok, _ := node.Valid()
			So(ok, ShouldEqual, true)

			var err error
			node.Id, err = adaptors.Node.Add(node)
			So(err, ShouldBeNil)

			// add device
			device := &m.Device{
				Name: "device",
				Status: "enabled",
				Type: "default",
				Node: node,
				Properties: []byte("{}"),
			}

			ok, _ = device.Valid()
			So(ok, ShouldEqual, true)

			deviceId, err := adaptors.Device.Add(device)
			So(err, ShouldBeNil)
			device.Id = deviceId

			// add script
			script8 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test8",
				Source:      coffeeScript8,
				Description: "test8",
			}
			script9 := &m.Script{
				Lang:        "coffeescript",
				Name:        "test9",
				Source:      coffeeScript9,
				Description: "test9",
			}

			ok, _ = script8.Valid()
			So(ok, ShouldEqual, true)

			ok, _ = script9.Valid()
			So(ok, ShouldEqual, true)

			engine8, err := scriptService.NewEngine(script8)
			So(err, ShouldBeNil)
			err = engine8.Compile()
			So(err, ShouldBeNil)
			script8Id, err := adaptors.Script.Add(script8)
			So(err, ShouldBeNil)
			script8, err = adaptors.Script.GetById(script8Id)
			So(err, ShouldBeNil)

			engine9, err := scriptService.NewEngine(script9)
			So(err, ShouldBeNil)
			err = engine9.Compile()
			So(err, ShouldBeNil)
			script9Id, err := adaptors.Script.Add(script9)
			So(err, ShouldBeNil)
			script9, err = adaptors.Script.GetById(script9Id)
			So(err, ShouldBeNil)

			// add device action
			deviceAction := &m.DeviceAction{
				Name: "deviceAction",
				DeviceId: deviceId,
				ScriptId: script8Id,
			}
			deviceAction.Id, err = adaptors.DeviceAction.Add(deviceAction)
			So(err, ShouldBeNil)

			// add workflow
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
			wfScenario1 := &m.WorkflowScenario{
				Name:       "wf scenario 1",
				SystemName: "wf_scenario_1",
				WorkflowId: workflow.Id,
			}

			ok, _ = wfScenario1.Valid()
			So(ok, ShouldEqual, true)

			wfScenario1.Id, err = adaptors.WorkflowScenario.Add(wfScenario1)
			So(err, ShouldBeNil)

			// add flow1
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

			feHandler := &m.FlowElement{
				Name:          "handler",
				FlowId:        flow1.Id,
				Status:        Enabled,
				PrototypeType: FlowElementsPrototypeMessageHandler,
				ScriptId:      &script9.Id,
			}

			ok, _ = feHandler.Valid()
			So(ok, ShouldEqual, true)

			feHandler.Uuid, err = adaptors.FlowElement.Add(feHandler)
			So(err, ShouldBeNil)

			// add worker
			worker := &m.Worker{
				Name: "worker",
				Time: "* * * * * *",
				Status: "enabled",
				WorkflowId: workflow.Id,
				FlowId: flow1.Id,
				DeviceActionId: deviceAction.Id,
			}

			ok, _ = worker.Valid()
			So(ok, ShouldEqual, true)

			worker.Id, err = adaptors.Worker.Add(worker)
			So(err, ShouldBeNil)

			// get flow
			flow1, err = adaptors.Flow.GetById(flow1.Id)
			So(err, ShouldBeNil)

			//fmt.Println("----")
			//debug.Println(flow1)
			//fmt.Println("----")
		})
	})
}