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

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	"strconv"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

type ControllerDeviceState struct {
	*ControllerCommon
}

func NewControllerDeviceState(common *ControllerCommon) *ControllerDeviceState {
	return &ControllerDeviceState{ControllerCommon: common}
}

// swagger:operation POST /device_state deviceStateAdd
// ---
// parameters:
// - description: device state params
//   in: body
//   name: device_state
//   required: true
//   schema:
//     $ref: '#/definitions/NewDeviceState'
//     type: object
// summary: add new device state
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) Add(ctx *gin.Context) {

	var params models.NewDeviceState
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	state := &m.DeviceState{}
	_ = common.Copy(&state, &params)

	if params.Device != nil && params.Device.Id != 0 {
		state.DeviceId = params.Device.Id
	}

	_, id, errs, err := c.endpoint.DeviceState.Add(state)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	state, err = c.endpoint.DeviceState.GetById(id)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.DeviceState{}
	_ = common.Copy(&result, &state, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /device_state/{id} deviceStateGetById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: get device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	state, err := c.endpoint.DeviceState.GetById(int64(aid))
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	result := &models.DeviceState{}
	_ = common.Copy(&result, &state, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation PUT /device_state/{id} deviceStateUpdateById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// - description: Update device state params
//   in: body
//   name: device_state
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateDeviceState'
//     type: object
// summary: update device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/DeviceState'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) Update(ctx *gin.Context) {

	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	var params models.UpdateDeviceState
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	state := &m.DeviceState{}
	_ = common.Copy(&state, &params, common.JsonEngine)

	if params.Device != nil && params.Device.Id != 0 {
		state.DeviceId = params.Device.Id
	}

	state.Id = int64(aid)

	state, errs, err := c.endpoint.DeviceState.Update(state)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	result := &models.DeviceState{}
	_ = common.Copy(&result, &state, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation DELETE /device_state/{id} deviceStateDeleteById
// ---
// parameters:
// - description: DeviceState ID
//   in: path
//   name: id
//   required: true
//   type: integer
// summary: delete device state by id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) Delete(ctx *gin.Context) {

	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if err := c.endpoint.DeviceState.Delete(int64(aid)); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// swagger:operation GET /device_states/{id} deviceStateList
// ---
// summary: get device state list by device id
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - device_state
// parameters:
// - description: Device ID
//   in: path
//   name: id
//   required: true
//   type: integer
// responses:
//   "200":
//     description: OK
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/DeviceState'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerDeviceState) GetStateList(ctx *gin.Context) {

	id := ctx.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	items, err := c.endpoint.DeviceState.GetList(int64(deviceId))
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	result := make([]*models.DeviceState, 0)
	_ = common.Copy(&result, &items, common.JsonEngine)

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}
