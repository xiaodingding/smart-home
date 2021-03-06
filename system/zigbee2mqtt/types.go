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

package zigbee2mqtt

import (
	m "github.com/e154/smart-home/models"
	"time"
)

type DeviceType string

type BridgeLog struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta"`
}

type BridgePairingMeta struct {
	FriendlyName string `json:"friendly_name"`
	Model        string `json:"model"`
	Vendor       string `json:"vendor"`
	Description  string `json:"description"`
	Supported    bool   `json:"supported"`
}

type BridgeConfigMeta struct {
	Transportrev int64 `json:"transportrev"`
	Product      int64 `json:"product"`
	Majorrel     int64 `json:"majorrel"`
	Minorrel     int64 `json:"minorrel"`
	Maintrel     int64 `json:"maintrel"`
	Revision     int64 `json:"revision"`
}

type BridgeConfigCoordinator struct {
	Type string           `json:"type"`
	Meta BridgeConfigMeta `json:"meta"`
}

type BridgeConfig struct {
	Version     string                  `json:"version"`
	Commit      string                  `json:"commit"`
	Coordinator BridgeConfigCoordinator `json:"coordinator"`
	LogLevel    string                  `json:"log_level"`
	PermitJoin  string                  `json:"permit_join"`
}

type AssistDeviceInfo struct {
	Name         string `json:"name"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
}

type AssistDevice struct {
	Device AssistDeviceInfo `json:"device"`
}

const (
	active  = "active"
	banned  = "banned"
	removed = "removed"
)

type Zigbee2mqttInfo struct {
	ScanInProcess bool          `json:"scan_in_process"`
	LastScan      time.Time     `json:"last_scan"`
	Networkmap    string        `json:"networkmap"`
	Status        string        `json:"status"`
	Model         m.Zigbee2mqtt `json:"model"`
}
