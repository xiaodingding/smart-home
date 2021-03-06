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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type Message struct {
	table *db.Messages
}

func GetMessageAdaptor(d *gorm.DB) *Message {
	return &Message{
		table: &db.Messages{Db: d},
	}
}

func (n *Message) Add(msg *m.Message) (id int64, err error) {
	id, err = n.table.Add(n.toDb(msg))
	return
}

func (n *Message) fromDb(dbVer *db.Message) (ver *m.Message) {
	ver = &m.Message{
		Id:           dbVer.Id,
		Type:         m.MessageType(dbVer.Type),
		EmailFrom:    dbVer.EmailFrom,
		EmailSubject: dbVer.EmailSubject,
		EmailBody:    dbVer.EmailBody,
		SmsText:      dbVer.SmsText,
		SlackText:    dbVer.SlackText,
		UiText:       dbVer.UiText,
		TelegramText: dbVer.TelegramText,
		CreatedAt:    dbVer.CreatedAt,
		UpdatedAt:    dbVer.UpdatedAt,
	}
	return
}

func (n *Message) toDb(ver *m.Message) (dbVer *db.Message) {
	dbVer = &db.Message{
		Id:           ver.Id,
		Type:         string(ver.Type),
		EmailFrom:    ver.EmailFrom,
		EmailSubject: ver.EmailSubject,
		EmailBody:    ver.EmailBody,
		SmsText:      ver.SmsText,
		SlackText:    ver.SlackText,
		UiText:       ver.UiText,
		TelegramText: ver.TelegramText,
		CreatedAt:    ver.CreatedAt,
		UpdatedAt:    ver.UpdatedAt,
	}
	return
}
