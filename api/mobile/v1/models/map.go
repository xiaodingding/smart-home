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

package models

type MapOptions struct {
	Zoom              float64 `json:"zoom"`
	ElementStateText  bool    `json:"element_state_text"`
	ElementOptionText bool    `json:"element_option_text"`
}

// swagger:model
type NewMap struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type UpdateMap struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type Map struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

// swagger:model
type MapFull struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Options     MapOptions  `json:"options"`
	Layers      []*MapLayer `json:"layers"`
}
