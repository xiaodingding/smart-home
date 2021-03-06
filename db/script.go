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

package db

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"time"
	. "github.com/e154/smart-home/common"
)

type Scripts struct {
	Db *gorm.DB
}

type Script struct {
	Id          int64 `gorm:"primary_key"`
	Lang        ScriptLang
	Name        string
	Source      string
	Description string
	Compiled    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *Script) TableName() string {
	return "scripts"
}

func (n Scripts) Add(node *Script) (id int64, err error) {
	if err = n.Db.Create(&node).Error; err != nil {
		return
	}
	id = node.Id
	return
}

func (n Scripts) GetById(nodeId int64) (node *Script, err error) {
	node = &Script{Id: nodeId}
	err = n.Db.First(&node).Error
	return
}

func (n Scripts) Update(m *Script) (err error) {
	err = n.Db.Model(&Script{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"lang":        m.Lang,
		"source":      m.Source,
		"compiled":    m.Compiled,
	}).Error
	return
}

func (n Scripts) Delete(nodeId int64) (err error) {
	err = n.Db.Delete(&Script{Id: nodeId}).Error
	return
}

func (n *Scripts) List(limit, offset int64, orderBy, sort string) (list []*Script, total int64, err error) {

	if err = n.Db.Model(Script{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Script, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}

func (n *Scripts) Search(query string, limit, offset int) (list []*Script, total int64, err error) {

	q := n.Db.Model(&Script{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Script, 0)
	err = q.Find(&list).Error

	return
}

