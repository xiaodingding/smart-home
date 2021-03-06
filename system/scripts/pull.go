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

package scripts

import (
	"sync"
)

type Pull struct {
	sync.Mutex
	heap map[string]interface{}
}

func NewPull() *Pull {
	return &Pull{
		heap: make(map[string]interface{}),
	}
}

func (p *Pull) Get(name string) (value interface{}, ok bool) {
	p.Lock()
	value, ok = p.heap[name]
	p.Unlock()
	return
}

func (p *Pull) Add(name string, s interface{}) {
	p.Lock()
	defer p.Unlock()

	p.heap[name] = s
}
