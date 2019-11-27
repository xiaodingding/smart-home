package debug

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/system/validation"
	"runtime"
	"strings"
)

// https://github.com/Unknwon/gcblog/blob/master/content/04-go-caller.md
func CallerName(skip int) (name, file string, line int, ok bool) {
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip + 1); !ok {
		return
	}
	name = runtime.FuncForPC(pc).Name()
	return
}

func Trace() (trace string) {

	i := 1 //0...
	for skip := i; ; skip++ {
		name, file, line, ok := CallerName(skip)
		if !ok {
			break
		}
		fn := strings.Title(strings.Split(name, ".")[1]) + "()"
		trace += "\n"
		trace += fmt.Sprintf("called: %s:%s line: %d", file, fn, line)
	}

	return
}

func Println(i interface{}) {
	b, _ := json.MarshalIndent(i, " ", "  ")
	fmt.Println(string(b))
}

func PrintValidationErrs(errs []*validation.Error) {
	for _, err := range errs {
		fmt.Printf("%s - %s", err.Name, err.String())
	}
}
