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
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test10(t *testing.T) {

	type Foo struct {
		Bar string
		Foo *Foo
	}

	var script1 *m.Script
	Convey("require external library", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService) {

			scriptService.PushFunctions("So", func(actual interface{}, assert string, expected interface{}) {
				//fmt.Printf("actual(%v), expected(%v)\n", actual, expected)
				switch assert {
				case "ShouldEqual":
					So(fmt.Sprintf("%v", actual), ShouldEqual, expected)
				}

			})

			script1 = &m.Script{
				Lang:        "coffeescript",
				Name:        "test10",
				Source:      coffeeScripts["coffeeScript25"],
				Description: "test10",
			}

			bar := &Foo{
				Bar: "bar",
				Foo: &Foo{
					Bar: "foo",
				},
			}

			scriptService.PushStruct("bar2", bar)

			scriptService.PushFunctions("external", func(varName string, f *Foo) {
				//fmt.Printf("varName: %v\n", varName)
				switch varName {
				case "bar":
					So(f.Foo.Bar, ShouldEqual, "")
				case "bar2":
					So(f.Foo.Bar, ShouldEqual, "foo")
				case "IC.bar2":
					So(f.Foo.Bar, ShouldEqual, "foo")
				}
			})

			engine, err := scriptService.NewEngine(script1)
			So(err, ShouldBeNil)

			_, err = engine.PushStruct("bar", bar)
			So(err, ShouldBeNil)

			counter := engine.PushGlobalProxy("bar2", bar)
			//fmt.Println(counter)
			So(counter, ShouldEqual, 12)

			counter = engine.PushGlobalProxy("bar2", bar)
			//fmt.Println(counter)
			So(counter, ShouldEqual, 12)

			err = engine.Compile()
			So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)

			engine.Gc()

			counter = engine.PushGlobalProxy("bar2", bar)
			//fmt.Println(counter)
			So(counter, ShouldEqual, 15)

			err = engine.Compile()
			So(err, ShouldBeNil)

			_, err = engine.Do()
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 2)
		})
	})
}
