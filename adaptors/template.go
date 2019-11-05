package adaptors

import (
	"errors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type Template struct {
	table *db.Templates
	db    *gorm.DB
}

func GetTemplateAdaptor(d *gorm.DB) *Template {
	return &Template{
		table: &db.Templates{Db: d},
		db:    d,
	}
}

func (n *Template) UpdateOrCreate(ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateOrCreate(dbVer); err != nil {
		return
	}

	return
}

func (n *Template) Create(ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.Create(dbVer); err != nil {
		return
	}

	return
}

func (n *Template) UpdateStatus(ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateStatus(dbVer); err != nil {
		return
	}

	return
}

func (n *Template) GetList(templateType m.TemplateType) (items []*m.Template, err error) {

	var dbItems []*db.Template
	if dbItems, err = n.table.GetList(templateType.String()); err != nil {
		return
	}

	items = make([]*m.Template, 0, len(dbItems))
	for _, dbVer := range dbItems {
		items = append(items, n.fromDb(dbVer))
	}

	return
}

func (n *Template) GetByName(name string) (ver *m.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(name, "template"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

func (n *Template) GetItemByName(name string) (ver *m.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(name, "item"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

func (n *Template) GetItemsSortedList() (count int64, items []string, err error) {
	count, items, err = n.table.GetItemsSortedList()
	return
}

func (n *Template) Delete(name string) (err error) {
	err = n.table.Delete(name)
	return
}

func (n *Template) GetItemsTree() (tree []*m.TemplateTree, err error) {

	var dbTree []*db.TemplateTree
	if dbTree, err = n.table.GetItemsTree(); err != nil {
		return
	}

	tree = make([]*m.TemplateTree, 0, len(dbTree))
	err = common.Copy(&tree, &dbTree, common.JsonEngine)

	return
}

func (n *Template) UpdateItemsTree(tree []*m.TemplateTree) (err error) {

	dbTree := make([]*db.TemplateTree, 0)
	if err = common.Copy(&dbTree, &tree, common.JsonEngine); err != nil {
		err = errors.New(err.Error())
		return
	}

	if err = n.table.UpdateItemsTree(dbTree, ""); err != nil {
		return
	}

	return
}

func (n *Template) Search(query string, limit, offset int) (list []*m.Template, total int64, err error) {
	var dbList []*db.Template
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Template, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *Template) Render(name string, params map[string]interface{}) (render *m.TemplateRender, err error) {

	var item *m.Template
	var template *m.TemplateContent
	var items m.Templates

	if item, err = n.GetByName(name); err != nil {
		return
	}

	if template, err = item.GetTemplate(); err != nil {
		return
	}

	if items, err = n.GetList(m.TemplateTypeItem); err != nil {
		return
	}

	render, err = m.RenderTemplate(items, template, params)

	return
}

func (n *Template) fromDb(dbVer *db.Template) (ver *m.Template) {
	ver = &m.Template{
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Content:     dbVer.Content,
		Status:      m.TemplateStatus(dbVer.Status),
		Type:        m.TemplateType(dbVer.Type),
		ParentName:  dbVer.ParentName,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	return
}

func (n *Template) toDb(ver *m.Template) (dbVer *db.Template) {
	dbVer = &db.Template{
		Name:        ver.Name,
		Description: ver.Description,
		Content:     ver.Content,
		Status:      ver.Status.String(),
		Type:        ver.Type.String(),
		ParentName:  ver.ParentName,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}
	return
}