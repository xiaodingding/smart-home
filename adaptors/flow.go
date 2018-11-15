package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Flow struct {
	table *db.Flows
	db    *gorm.DB
}

func GetFlowAdaptor(d *gorm.DB) *Flow {
	return &Flow{
		table: &db.Flows{Db: d},
		db:    d,
	}
}

func (n *Flow) Add(flow *m.Flow) (id int64, err error) {

	dbFlow := n.toDb(flow)
	if id, err = n.table.Add(dbFlow); err != nil {
		return
	}

	return
}

func (n *Flow) GetAllEnabled() (list []*m.Flow, err error) {

	var dbList []*db.Flow
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Flow, 0)
	for _, dbFlow := range dbList {
		flow := n.fromDb(dbFlow)
		list = append(list, flow)
	}

	return
}

func (n *Flow) GetById(flowId int64) (flow *m.Flow, err error) {

	var dbFlow *db.Flow
	if dbFlow, err = n.table.GetById(flowId); err != nil {
		return
	}

	flow = n.fromDb(dbFlow)

	return
}

func (n *Flow) Update(flow *m.Flow) (err error) {
	dbFlow := n.toDb(flow)
	err = n.table.Update(dbFlow)
	return
}

func (n *Flow) Delete(flowId int64) (err error) {
	err = n.table.Delete(flowId)
	return
}

func (n *Flow) List(limit, offset int64, orderBy, sort string) (list []*m.Flow, total int64, err error) {
	var dbList []*db.Flow
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Flow, 0)
	for _, dbFlow := range dbList {
		flow := n.fromDb(dbFlow)
		list = append(list, flow)
	}

	return
}

func (n *Flow) fromDb(dbFlow *db.Flow) (flow *m.Flow) {
	flow = &m.Flow{
		Id:                 dbFlow.Id,
		Name:               dbFlow.Name,
		Status:             dbFlow.Status,
		Description:        dbFlow.Description,
		WorkflowId:         dbFlow.WorkflowId,
		WorkflowScenarioId: dbFlow.WorkflowScenarioId,
		Connections:        make([]*m.Connection, 0),
		FlowElements:       make([]*m.FlowElement, 0),
		CreatedAt:          dbFlow.CreatedAt,
		UpdatedAt:          dbFlow.UpdatedAt,
	}

	return
}

func (n *Flow) toDb(flow *m.Flow) (dbFlow *db.Flow) {
	dbFlow = &db.Flow{
		Id:                 flow.Id,
		Name:               flow.Name,
		Status:             flow.Status,
		Description:        flow.Description,
		WorkflowScenarioId: flow.WorkflowScenarioId,
		WorkflowId:         flow.WorkflowId,
		CreatedAt:          flow.CreatedAt,
		UpdatedAt:          flow.UpdatedAt,
	}
	return
}