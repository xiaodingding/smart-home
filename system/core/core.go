package core

import (
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("core")
)

type Core struct {
	sync.Mutex
	nodes         map[int64]*Node
	workflows     map[int64]*Workflow
	adaptors      *adaptors.Adaptors
	scripts       *scripts.ScriptService
	cron          *cr.Cron
	mqtt          *mqtt.Mqtt
	telemetry     telemetry.ITelemetry
	streamService *stream.StreamService
	Map           *Map
	isRunning     bool
}

func NewCore(adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	graceful *graceful_service.GracefulService,
	cron *cr.Cron,
	mqtt *mqtt.Mqtt,
	telemetry telemetry.ITelemetry,
	streamService *stream.StreamService) (core *Core, err error) {

	core = &Core{
		nodes:         make(map[int64]*Node),
		workflows:     make(map[int64]*Workflow),
		adaptors:      adaptors,
		scripts:       scripts,
		cron:          cron,
		mqtt:          mqtt,
		telemetry:     telemetry,
		streamService: streamService,
		Map: &Map{
			telemetry: telemetry,
		},
	}

	graceful.Subscribe(core)

	scripts.PushStruct("Map", &MapBind{Map: core.Map})

	return
}

func (c *Core) Run() (err error) {

	if c.safeIsRunning() {
		return
	}

	c.isRunning = true

	if err = c.initNodes(); err != nil {
		return
	}

	if err = c.InitWorkflows(); err != nil {
		return
	}

	return
}

func (b *Core) Stop() (err error) {

	if !b.safeIsRunning() {
		return
	}

	defer func() {
		b.isRunning = false
	}()

	// unregister steam actions
	b.streamService.UnSubscribe("do.worker")
	b.streamService.UnSubscribe("do.action")

	for _, workflow := range b.workflows {
		if err = b.DeleteWorkflow(workflow.model); err != nil {
			return
		}
	}

	for _, node := range b.nodes {
		if err = b.RemoveNode(&m.Node{Id: node.Id, Name: node.Name}); err != nil {
			return
		}
	}

	return
}

func (b *Core) Shutdown() {
	if err := b.Stop(); err != nil {
		log.Error(err.Error())
	}
}

// ------------------------------------------------
// Nodes
// ------------------------------------------------

func (c *Core) initNodes() (err error) {

	var nodes []*m.Node
	if nodes, err = c.adaptors.Node.GetAllEnabled(); err != nil {
		return
	}

	for _, modelNode := range nodes {
		if _, err = c.AddNode(modelNode); err != nil {
			log.Error(err.Error())
		}
	}

	return
}

func (c *Core) AddNode(node *m.Node) (n *Node, err error) {

	if _, exist := c.safeGetNode(node.Id); exist {
		err = c.ReloadNode(node)
		return
	}

	log.Infof("Add node: \"%s\"", node.Name)

	n = NewNode(node, c.mqtt)
	c.safeUpdateNodeMap(node.Id, n.Connect())

	go c.telemetry.Broadcast(telemetry.Node{})

	return
}

func (b *Core) RemoveNode(node *m.Node) (err error) {

	log.Infof("Remove node: \"%s\"", node.Name)

	err = b.removeNode(node)

	b.telemetry.Broadcast(telemetry.Node{})

	return
}

func (c *Core) removeNode(node *m.Node) (err error) {

	n, exist := c.safeGetNode(node.Id)
	if !exist {
		err = errors.New("not found")
		return
	}

	n.Disconnect()

	c.Lock()
	delete(c.nodes, node.Id)
	c.Unlock()

	return
}

func (c *Core) ReloadNode(node *m.Node) (err error) {

	log.Infof("Reload node: \"%s\"", node.Name)

	n, exist := c.safeGetNode(node.Id)
	if !exist {
		if _, err = c.AddNode(node); err != nil {
			log.Error(err.Error())
		}
		return
	}

	n.Status = node.Status

	if n.Status == "disabled" {
		n.Disconnect()
	} else {
		n.Connect()
	}

	return
}

func (c *Core) ConnectNode(node *m.Node) (err error) {

	log.Infof("Connect to node: \"%s\"", node.Name)

	if n, exist := c.safeGetNode(node.Id); exist {
		n.Connect()
	}

	c.telemetry.Broadcast(telemetry.Node{})

	return
}

func (c *Core) DisconnectNode(node *m.Node) (err error) {

	log.Infof("Disconnect from node: \"%s\"", node.Name)

	if n, exist := c.safeGetNode(node.Id); exist {
		n.Disconnect()
	}

	c.telemetry.Broadcast(telemetry.Node{})

	return
}

func (b *Core) GetNodes() (nodes map[int64]*Node) {

	nodes = make(map[int64]*Node)

	b.Lock()
	for id, node := range b.nodes {
		nodes[id] = node
	}
	b.Unlock()

	return
}

func (c *Core) GetNodeById(nodeId int64) *Node {

	if n, exist := c.safeGetNode(nodeId); exist {
		return n
	}

	return nil
}

// ------------------------------------------------
// Workflows
// ------------------------------------------------

// инициализация всего рабочего процесса, с запуском
// дочерни подпроцессов
func (b *Core) InitWorkflows() (err error) {

	workflows, err := b.adaptors.Workflow.GetAllEnabled()
	if err != nil {
		return
	}

	for _, workflow := range workflows {
		if err = b.AddWorkflow(workflow); err != nil {
			return
		}
	}

	return
}

// добавление рабочего процесс
func (b *Core) AddWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Add workflow: '%s'", workflow.Name)

	if _, ok := b.safeGetWorkflow(workflow.Id); ok {
		return
	}

	wf := NewWorkflow(workflow, b.adaptors, b.scripts, b.cron, b, b.mqtt, b.telemetry)

	if err = wf.Run(); err != nil {
		return
	}

	b.safeUpdateWorkflowMap(workflow.Id, wf)

	return
}

func (b *Core) GetWorkflow(workflowId int64) (workflow *Workflow, err error) {

	log.Infof("GetWorkflow: id(%v)", workflowId)

	var ok bool
	if workflow, ok = b.safeGetWorkflow(workflowId); !ok {
		err = errors.New("not found")
		return
	}

	return
}

func (b *Core) GetStatusAllWorkflow() (statusList []m.DashboardWorkflowStatus) {

	b.Lock()
	defer b.Unlock()

	statusList = make([]m.DashboardWorkflowStatus, 0, len(b.workflows))
	for _, workflow := range b.workflows {
		statusList = append(statusList, m.DashboardWorkflowStatus{
			Id:         workflow.model.Id,
			ScenarioId: workflow.model.Scenario.Id,
		})
	}

	return
}

func (c *Core) GetStatusWorkflow(workflowId int64) (status m.DashboardWorkflowStatus, err error) {

	workflow, ok := c.safeGetWorkflow(workflowId)
	if !ok {
		err = errors.New("not found")
		return
	}

	status = m.DashboardWorkflowStatus{
		Id:         workflow.model.Id,
		ScenarioId: workflow.model.Scenario.Id,
	}

	return
}

// нельзя удалить workflow, если присутствуют связанные сущности
func (c *Core) DeleteWorkflow(workflow *m.Workflow) (err error) {

	log.Infof("Remove workflow: %s", workflow.Name)

	wf, ok := c.safeGetWorkflow(workflow.Id)
	if !ok {
		err = errors.New("not found")
		return
	}

	if err = wf.Stop(); err != nil {
		log.Error(err.Error())
	}

	c.Lock()
	delete(c.workflows, workflow.Id)
	c.Unlock()

	return
}

func (c *Core) UpdateWorkflowScenario(workflow *m.Workflow) (err error) {

	wf, ok := c.safeGetWorkflow(workflow.Id)
	if !ok {
		err = errors.New("not found")
		return
	}

	err = wf.UpdateScenario()

	return
}

func (c *Core) UpdateWorkflow(workflow *m.Workflow) (err error) {

	if workflow.Status == "enabled" {
		if _, ok := c.safeGetWorkflow(workflow.Id); ok {
			err = c.AddWorkflow(workflow)
			return
		}
	} else {
		if _, ok := c.safeGetWorkflow(workflow.Id); ok {
			err = c.DeleteWorkflow(workflow)
			return
		}
	}

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

func (c *Core) AddFlow(flow *m.Flow) (err error) {

	wf, ok := c.safeGetWorkflow(flow.WorkflowId)
	if !ok {
		err = errors.New("not found")
		return
	}

	if err = wf.AddFlow(flow); err != nil {
		return
	}

	return
}

func (c *Core) GetFlow(id int64) (*Flow, error) {

	var flow *m.Flow
	var err error
	if flow, err = c.adaptors.Flow.GetById(id); err != nil {
		return nil, err
	}

	wf, ok := c.safeGetWorkflow(flow.WorkflowId)
	if !ok {
		return nil, nil
	}

	return wf.GetFLow(id)
}

func (c *Core) UpdateFlow(flow *m.Flow) error {

	wf, ok := c.safeGetWorkflow(flow.WorkflowId)
	if !ok {
		return nil
	}

	return wf.UpdateFlow(flow)
}

func (c *Core) RemoveFlow(flow *m.Flow) error {

	wf, ok := c.safeGetWorkflow(flow.WorkflowId)
	if !ok {
		return nil
	}

	return wf.RemoveFlow(flow)
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (b *Core) UpdateFlowFromDevice(device *m.Device) (err error) {

	//	var flows map[int64]*m.Flow
	//	flows = make(map[int64]*m.Flow)
	//	childs, _, _ := device.GetChilds()
	//
	//	for _, workflow := range b.workflows {
	//		for _, flow := range workflow.Flows {
	//			for _, worker := range flow.Workers {
	//				for _, action := range worker.actions {
	//					//if action.Device.Id == device.Id {
	//					//	workflow.UpdateFlow(flow.Model)
	//					//	continue
	//					//}
	//
	//					if action.Device != nil && action.Device.Id == device.Id {
	//						//workflow.UpdateFlow(flow.Model)
	//						flows[flow.Model.Id] = flow.Model
	//						continue
	//					}
	//
	//					for _, child := range childs {
	//						if action.Device != nil && action.Device.Id == child.Id {
	//							flows[flow.Model.Id] = flow.Model
	//						}
	//					}
	//				}
	//
	//				if device.Device != nil && worker.Model.DeviceAction.Device.Id == device.Device.Id {
	//					//workflow.UpdateFlow(flow.Model)
	//					flows[flow.Model.Id] = flow.Model
	//					continue
	//				}
	//			}
	//		}
	//
	//		for _, flow := range flows {
	//			workflow.UpdateFlow(flow)
	//		}
	//
	//		flows = make(map[int64]*m.Flow)
	//	}

	return
}

func (b *Core) UpdateWorker(_worker *m.Worker) (err error) {

	//b.Lock()
	//defer b.Unlock()

	for _, workflow := range b.workflows {
		_ = workflow.UpdateWorker(_worker)
	}

	return
}

func (b *Core) RemoveWorker(worker *m.Worker) (err error) {

	//b.Lock()
	//defer b.Unlock()

	for _, workflow := range b.workflows {
		_ = workflow.RemoveWorker(worker)
	}

	return
}

func (b *Core) DoWorker(worker *m.Worker) (err error) {

	//b.Lock()
	//defer b.Unlock()

	for _, workflow := range b.workflows {
		_ = workflow.DoWorker(worker)
	}

	return
}

// ------------------------------------------------
// safe methods
// ------------------------------------------------

func (b *Core) safeIsRunning() bool {
	b.Lock()
	defer b.Unlock()
	return b.isRunning
}

func (b *Core) safeSetIsRunning(v bool) {
	b.Lock()
	b.isRunning = v
	b.Unlock()
}

func (c *Core) safeGetWorkflow(k int64) (w *Workflow, ok bool) {
	c.Lock()
	w, ok = c.workflows[k]
	c.Unlock()
	return
}

func (c *Core) safeUpdateWorkflowMap(k int64, w *Workflow) {
	c.Lock()
	c.workflows[k] = w
	c.Unlock()
}

func (c *Core) safeGetNode(k int64) (w *Node, ok bool) {
	c.Lock()
	w, ok = c.nodes[k]
	c.Unlock()
	return
}

func (c *Core) safeUpdateNodeMap(k int64, n *Node) {
	c.Lock()
	c.nodes[k] = n
	c.Unlock()
}
