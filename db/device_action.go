package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

type DeviceActions struct {
	Db *gorm.DB
}

type DeviceAction struct {
	Id          int64 `gorm:"primary_key"`
	Device      *Device
	DeviceId    int64
	Name        string
	Description string
	Script      *Script
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *DeviceAction) TableName() string {
	return "device_actions"
}

func (n DeviceActions) Add(action *DeviceAction) (id int64, err error) {
	if err = n.Db.Create(&action).Error; err != nil {
		return
	}
	id = action.Id
	return
}

func (n DeviceActions) GetById(actionId int64) (action *DeviceAction, err error) {
	action = &DeviceAction{Id: actionId}
	err = n.Db.First(&action).Error
	return
}

func (n DeviceActions) Update(m *DeviceAction) (err error) {
	err = n.Db.Model(&DeviceAction{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"script":      m.Script,
	}).Error
	return
}

func (n DeviceActions) Delete(actionId int64) (err error) {
	err = n.Db.Delete(&DeviceAction{Id: actionId}).Error
	return
}

func (n *DeviceActions) List(limit, offset int64, orderBy, sort string) (list []*DeviceAction, total int64, err error) {

	if err = n.Db.Model(DeviceAction{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*DeviceAction, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	return
}