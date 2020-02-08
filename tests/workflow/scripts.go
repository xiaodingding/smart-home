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

package workflow

var coffeeScripts = map[string]string{
	"coffeeScript1": coffeeScript1,
	"coffeeScript2": coffeeScript2,
	"coffeeScript3": coffeeScript3,
	"coffeeScript4": coffeeScript4,
	"coffeeScript5": coffeeScript5,
	"coffeeScript6": coffeeScript6,
	"coffeeScript7": coffeeScript7,
	"coffeeScript8": coffeeScript8,
	"coffeeScript9": coffeeScript9,
	"coffeeScript10": coffeeScript10,
	"coffeeScript11": coffeeScript11,
	"coffeeScript12": coffeeScript12,
	"coffeeScript13": coffeeScript13,
	"coffeeScript14": coffeeScript14,
	//"coffeeScript15": coffeeScript15,
	"coffeeScript16": coffeeScript16,
	"coffeeScript17": coffeeScript17,
	"coffeeScript18": coffeeScript18,
	"coffeeScript19": coffeeScript19,
	"coffeeScript20": coffeeScript20,
	"coffeeScript21": coffeeScript21,
	"coffeeScript22": coffeeScript22,
	"coffeeScript23": coffeeScript23,
	"coffeeScript24": coffeeScript24,
}

// test1, test2
// ------------------------------------------------
const coffeeScript1 = `
main =->
    wf = IC.Workflow()
    #print wf.getName() + " script1"
`

const coffeeScript2 = `
main =->
    wf = IC.Workflow()
    #print wf.getName() + " script2"
`

const coffeeScript3 = `
main =->
    #print "workflow scenario script 3"

on_enter =->
    wf = IC.Workflow()
    #print "enter to " + wf.getName() + " " + wf.getScenarioName()
    wf.setVar("var1", 123)
    wf.setScenario("wf_scenario_2")

on_exit =->
    wf = IC.Workflow()

    #print "exit from " + wf.getName() + " " + wf.getScenarioName()
    #print "description " + wf.getDescription()
`

const coffeeScript4 = `
main =->
    #print "workflow scenario script 4"

on_enter =->
    wf = IC.Workflow()

    #print "enter to " + wf.getName() + " " + wf.getScenarioName()
    var1 = wf.getVar("var1")
    #print "var: " + var1

on_exit =->
    wf = IC.Workflow()

    #print "exit from " + wf.getName() + " " + wf.getScenarioName()
    #print "description " + wf.getDescription()
`

const coffeeScript5 = `
# test 2

#print "run message emitter script (script 5)"
#print "message:", message
#print "ENV", ENV.a

#print "message.ENV:", message.getVar "ENV"

IC.store ENV.a
`

const coffeeScript6 = `
# test 2

ENV = {"a": "b"}
#print "run message handler script (script 6)"

message.setVar("ENV", ENV) 
`

const coffeeScript7 = `
# test 2
# workflow script

#print "run workflow script (script 7)"

on_enter =->
	#print 'on_enter'
on_exit =->
`

// test3
// ------------------------------------------------
const coffeeScript8 = `
#print "run workflow script (script 8)"
`

const coffeeScript9 = `
#print "run workflow script (script 9)"
`

// test 4
// ------------------------------------------------
const coffeeScript10 = `
#print "run workflow script (script 10)"
`

const coffeeScript11 = `
#print "run workflow script (script 11)"
`

// test6, test7
// ------------------------------------------------
const coffeeScript12 = `
#print "run workflow script (script 12)"
IC.store('script12')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript13 = `
#print "run workflow script (script 13)"
IC.store('script13')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript14 = `
#print "run workflow script (script 14)"
IC.store('script14')
#print "message:", message.getVar('val')
`

// test8, test9, test10
// ------------------------------------------------
const coffeeScript16 = `
#print "run workflow script (script 16)"
IC.store('script16')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript17 = `
#print "run workflow script (script 17)"

IC.store('script17')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript18 = `
#print "run workflow script (script 18)"

IC.store('script18')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript19 = `
#print "run workflow script (script 19)"

IC.store('script19')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript20 = `
#print "run workflow script (script 20)"

IC.store('script20')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript21 = `
#print "run workflow script (script 21)"

IC.store('script21')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

// test11
// ------------------------------------------------
const coffeeScript22 = `
#print "run workflow script (script 22)"
IC.store('script22')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript23 = `
#print "run workflow script (script 23)"
IC.store('script23')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

const coffeeScript24 = `
#print "run workflow script (script 24)"
IC.store('script24')
#print "message:", message.getVar('val')
c = message.getVar('val') + 1
message.setVar('val', c)
IC.store2(c)
`

// test8...
// ------------------------------------------------
