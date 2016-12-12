package core

import (
	"log"
	"sync"
	"github.com/astaxie/beego/orm"
	cr "github.com/e154/cron"
	r "../../lib/rpc"
	"../models"
	"../scripts"
	"time"
	"fmt"
)

type Worker struct {
	Model     *models.Worker
	CronTasks map[int64]*cr.Task
	Devices   map[int64]*models.Device
}

func NewWorkflow(model *models.Workflow, nodes map[int64]*models.Node) (workflow *Workflow) {

	workflow = &Workflow{
		model: model,
		Nodes: nodes,
		Flows: make(map[int64]*Flow),
		Workers: make(map[int64]*Worker),
		mutex: &sync.Mutex{},
	}

	return
}

type Workflow struct {
	model   	*models.Workflow
	mutex   	*sync.Mutex
	Flows   	map[int64]*Flow
	Workers 	map[int64]*Worker
	Nodes   	map[int64]*models.Node
}

func (wf *Workflow) Run() (err error) {

	err = wf.InitFlows()

	if err != nil {
		return
	}

	err = wf.InitWorkers()

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, flow := range wf.Flows {
		wf.RemoveFlow(flow.Model)
	}

	//TODO check!!!
	for _, worker := range wf.Workers {
		wf.RemoveWorker(worker.Model)
	}

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	err = wf.Run()

	return
}

// ------------------------------------------------
// Flows
// ------------------------------------------------

// получаем все связанные процессы
func (wf *Workflow) InitFlows() (err error) {

	var flows []*models.Flow
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {
		return
	}

	for _, flow := range flows {
		wf.AddFlow(flow)
	}

	return
}

// Flow должен быть полный:
// с Connections
// с FlowElements
// с Cursor
// с Workers
func (wf *Workflow) AddFlow(flow *models.Flow) (err error) {

	if flow.Status != "enabled" {
		return
	}

	log.Println("Add flow:", flow.Name)

	wf.mutex.Lock()
	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}

	if wf.Flows[flow.Id], err = NewFlow(flow, wf); err != nil {
		return

	}
	wf.mutex.Unlock()

	workers, err := models.GetAllEnabledWorkersByFlow(flow)
	if err != nil {
		return
	}

	for _, worker := range workers {
		if err = wf.AddWorker(worker); err != nil {
			break
		}
	}

	return
}

func (wf *Workflow) UpdateFlow(flow *models.Flow) (err error) {

	err = wf.RemoveFlow(flow)
	if err != nil {
		return
	}

	err = wf.AddFlow(flow)

	return
}

func (wf *Workflow) RemoveFlow(flow *models.Flow) (err error) {

	log.Println("Remove flow:", flow.Name)

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	// to do first remove all workers
	var workers	[]*models.Worker
	if workers, err = flow.GetWorkers(); err != nil {
		return
	}

	for _, worker := range workers {
		wf.RemoveWorker(worker)
	}

	wf.Flows[flow.Id].Close()
	delete(wf.Flows, flow.Id)

	return
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (wf *Workflow) InitWorkers() (err error) {

	//log.Println("------------------- WORKERS ---------------------")

	var workers	[]*models.Worker
	if workers, err = wf.model.GetAllEnabledWorkers(); err != nil {
		return
	}

	for _, worker := range workers {
		if err = wf.AddWorker(worker); err != nil {
			log.Println("error:", err.Error())
			return
		}
	}

	return
}

func (wf *Workflow) AddWorker(worker *models.Worker) (err error) {

	if worker.Status != "enabled" {
		return
	}

	log.Printf("Add worker: \"%s\"", worker.Name)

	if _, ok := wf.Workers[worker.Id]; ok {
		return
	}

	wf.Workers[worker.Id] = &Worker{Model:worker,}

	// vars
	// ------------------------------------------------
	flow := wf.Flows[worker.Flow.Id]
	message := &Message{}
	//command := []byte{}
	////if command, err = hex.DecodeString(worker.DeviceAction.Command); err != nil {
	//if command, err = hex.DecodeString("0300000005"); err != nil {
	//	return
	//}

	// get device
	// ------------------------------------------------
	var devices []*models.Device
	if worker.DeviceAction.Device.Address != nil {
		devices = append(devices, worker.DeviceAction.Device)
	} else {
		// значит тут группа устройств
		var childs []*models.Device
		if childs, _, err = worker.DeviceAction.Device.GetChilds(); err != nil {
			return
		}

		for _, child := range childs {
			if child.Address == nil || child.Status != "enabled" {
				continue
			}

			device := &models.Device{}
			*device = *worker.DeviceAction.Device
			device.Id = child.Id
			device.Name = child.Name
			device.Address = new(int)
			*device.Address = *child.Address
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	var node *models.Node
	if _, ok := wf.Nodes[worker.DeviceAction.Device.Node.Id]; ok {
		node = wf.Nodes[worker.DeviceAction.Device.Node.Id]
	} else {
		// autoload nodes
		node, err = models.GetNodeById(worker.DeviceAction.Device.Node.Id)
		if err != nil {
			return
		}

		CorePtr().AddNode(node)
	}

	// cron worker
	// ------------------------------------------------
	for _, device := range devices {

		//var _command []byte
		//_command = append(_command, byte(*device.Address))
		//_command = append(_command, command...)

		args := r.Request{
			Baud: device.Baud,
			Result: true,
			//Command: _command,
			Device: device.Tty,
			Line: "",
			StopBits: int(device.StopBite),
			Time: time.Now(),
			Timeout: device.Timeout,
		}

		// device
		if wf.Workers[worker.Id].Devices == nil {
			wf.Workers[worker.Id].Devices = make(map[int64]*models.Device)
		}

		wf.Workers[worker.Id].Devices[device.Id] = device

		// get script
		// ------------------------------------------------
		o := orm.NewOrm()
		if _, err = o.LoadRelated(worker.DeviceAction, "Script"); err != nil {
			return
		}

		// add script
		script, _ := scripts.New(worker.DeviceAction.Script)


		script.PushStruct("device", device)
		script.PushStruct("flow", flow.Model)
		script.PushStruct("node", node)
		script.PushStruct("message", message)

		script.PushFunction("modbus_send", func(command []byte) (result r.Result) {

			args.Command = command
			if err := node.ModbusSend(args, &result); err != nil {
				result.Error = err.Error()
			}

			return
		})

		script.PushFunction("flow_new_message", func(v []byte) string {

			message.Variable = v
			if err = flow.NewMessage(message); err != nil {
				log.Println("error" , err.Error())
				return err.Error()
			}

			return ""
		})

		// set cron
		// ------------------------------------------------
		if wf.Workers[worker.Id].CronTasks == nil {
			wf.Workers[worker.Id].CronTasks = make(map[int64]*cr.Task)
		}

		wf.Workers[worker.Id].CronTasks[device.Id] = cron.NewTask(worker.Time, func() {
			script.Do()
		})

	}

	return
}

func (wf *Workflow) UpdateWorker(worker *models.Worker) (err error) {

	if _, ok := wf.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}

	if err = wf.RemoveWorker(worker); err != nil {
		log.Println("error:", err.Error())
	}

	if err = wf.AddWorker(worker); err != nil {
		log.Println("error:", err.Error())
	}

	return
}

func (wf *Workflow) RemoveWorker(worker *models.Worker) (err error) {

	log.Printf("Remove worker: \"%s\"", worker.Name)

	if _, ok := wf.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
		return
	}

	// stop cron task
	for _, task := range wf.Workers[worker.Id].CronTasks {

		task.Disable()

		// remove task from cron
		cron.RemoveTask(task)
	}

	// delete worker
	delete(wf.Workers, worker.Id)

	return
}
