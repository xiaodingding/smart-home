package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"strings"
	"github.com/e154/smart-home/system/access_list"
)

func roles(adaptors *adaptors.Adaptors,
	accessList *access_list.AccessListService) {

	// admin role
	adminRole := &m.Role{
		Name: "admin",
	}
	err := adaptors.Role.Add(adminRole)
	So(err, ShouldBeNil)

	// demo role
	demoRole := &m.Role{
		Name: "demo",
	}
	err = adaptors.Role.Add(demoRole)
	So(err, ShouldBeNil)

	for pack, item := range *accessList.List {
		for right := range item {
			if strings.Contains(right, "read") ||
				strings.Contains(right, "view") {
				permission := &m.Permission{
					RoleName:    demoRole.Name,
					PackageName: pack,
					LevelName:   right,
				}

				_, err = adaptors.Permission.Add(permission)
				So(err, ShouldBeNil)
			}
		}
	}

	// user role
	userRole := &m.Role{
		Name:   "user",
		Parent: demoRole,
	}
	err = adaptors.Role.Add(userRole)
	So(err, ShouldBeNil)

	// add admin
	adminUser := &m.User{
		Nickname:          "admin",
		RoleName:          adminRole.Name,
		Email:             "admin@e154.ru",
		Lang:              "en",
		Status:            "active",
	}
	err = adminUser.SetPass("admin")
	So(err, ShouldBeNil)

	ok, _ := adminUser.Valid()
	So(ok, ShouldEqual, true)

	adminUser.Id, err = adaptors.User.Add(adminUser)
	So(err, ShouldBeNil)

	// add demo user
	demoUser := &m.User{
		Nickname:          "demo",
		RoleName:          demoRole.Name,
		Email:             "demo@e154.ru",
		Lang:              "en",
		Status:            "active",
	}
	err = demoUser.SetPass("demo")
	So(err, ShouldBeNil)

	ok, _ = demoUser.Valid()
	So(ok, ShouldEqual, true)

	demoUser.Id, err = adaptors.User.Add(demoUser)
	So(err, ShouldBeNil)

	// add base user
	baseUser := &m.User{
		Nickname:          "user",
		RoleName:          userRole.Name,
		Email:             "user@e154.ru",
		Lang:              "en",
		Status:            "active",
	}
	err = baseUser.SetPass("user")
	So(err, ShouldBeNil)

	ok, _ = baseUser.Valid()
	So(ok, ShouldEqual, true)

	baseUser.Id, err = adaptors.User.Add(baseUser)
	So(err, ShouldBeNil)
}
