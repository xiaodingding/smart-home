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

package endpoint

import (
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/validation"
)

type RoleEndpoint struct {
	*CommonEndpoint
}

func NewRoleEndpoint(common *CommonEndpoint) *RoleEndpoint {
	return &RoleEndpoint{
		CommonEndpoint: common,
	}
}

func (n *RoleEndpoint) Add(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Role.GetByName(params.Name)

	return
}

func (n *RoleEndpoint) GetByName(name string) (result *m.Role, err error) {

	result, err = n.adaptors.Role.GetByName(name)

	return
}

func (n *RoleEndpoint) Update(params *m.Role) (result *m.Role, errs []*validation.Error, err error) {

	role, err := n.adaptors.Role.GetByName(params.Name)
	if err != nil {
		return
	}

	if params.Parent.Name == "" {
		role.Parent = nil
	} else {
		role.Parent = &m.Role{
			Name: params.Parent.Name,
		}
	}

	// validation
	_, errs = role.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Role.Update(role); err != nil {
		return
	}

	role, err = n.adaptors.Role.GetByName(params.Name)

	return
}

func (n *RoleEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.List(limit, offset, order, sortBy)

	return
}

func (n *RoleEndpoint) Delete(name string) (err error) {

	if name == "admin" {
		err = errors.New("admin is base role")
		return
	}
	err = n.adaptors.Role.Delete(name)

	return
}

func (n *RoleEndpoint) Search(query string, limit, offset int) (result []*m.Role, total int64, err error) {

	result, total, err = n.adaptors.Role.Search(query, limit, offset)

	return
}

func (n *RoleEndpoint) GetAccessList(roleName string,
	accessListService *access_list.AccessListService) (accessList access_list.AccessList, err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	accessList, err = accessListService.GetFullAccessList(role)

	return
}

func (n *RoleEndpoint) UpdateAccessList(roleName string, accessListDif map[string]map[string]bool) (err error) {

	var role *m.Role
	if role, err = n.adaptors.Role.GetByName(roleName); err != nil {
		return
	}

	tx := n.adaptors.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	addPerms := make([]*m.Permission, 0)
	delPerms := make([]string, 0)
	for packName, pack := range accessListDif {
		for levelName, dir := range pack {
			if dir {
				addPerms = append(addPerms, &m.Permission{
					RoleName:    role.Name,
					PackageName: packName,
					LevelName:   levelName,
				})
			} else {
				delPerms = append(delPerms, levelName)
			}

			if len(delPerms) > 0 {
				if err = tx.Permission.Delete(packName, delPerms); err != nil {
					return
				}
				delPerms = []string{}
			}
		}
	}

	if len(addPerms) == 0 {
		return
	}

	for _, perm := range addPerms {
		tx.Permission.Add(perm)
	}

	err = tx.Commit()

	return
}
