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

package mqtt

import (
	"github.com/DrmagicE/gmqtt"
	"github.com/e154/smart-home/system/config"
	"time"
)

type MqttConfig struct {
	Port                       int
	RetryInterval              time.Duration
	RetryCheckInterval         time.Duration
	SessionExpiryInterval      time.Duration
	SessionExpireCheckInterval time.Duration
	QueueQos0Messages          bool
	MaxInflight                int
	MaxAwaitRel                int
	MaxMsgQueue                int
	DeliverMode                gmqtt.DeliverMode
}

func NewMqttConfig(cfg *config.AppConfig) *MqttConfig {
	return &MqttConfig{
		Port:                       cfg.MqttPort,
		RetryInterval:              cfg.MqttRetryInterval,
		RetryCheckInterval:         cfg.MqttRetryCheckInterval,
		SessionExpiryInterval:      cfg.MqttSessionExpiryInterval,
		SessionExpireCheckInterval: cfg.MqttSessionExpireCheckInterval,
		QueueQos0Messages:          cfg.MqttQueueQos0Messages,
		MaxInflight:                cfg.MqttMaxInflight,
		MaxAwaitRel:                cfg.MqttMaxAwaitRel,
		MaxMsgQueue:                cfg.MqttMaxMsgQueue,
		DeliverMode:                gmqtt.DeliverMode(cfg.MqttDeliverMode),
	}
}
