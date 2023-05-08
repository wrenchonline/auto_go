package main

import (
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {
	//robotgo.KeySleep = 50
	robotgo.SetKeyDelay(0)
	for {
		robotgo.KeySleep = 1
		robotgo.EventHook(hook.KeyDown, []string{"b"}, func(e hook.Event) {
			robotgo.EventEnd()
		})
		s := robotgo.EventStart()
		<-robotgo.EventProcess(s)
		ok := robotgo.AddEvents("b")
		if ok {
			//狂战血怒
			robotgo.KeyDown(`x`)

			robotgo.KeyUp(`x`)

			robotgo.KeyDown(`t`)

			robotgo.KeyUp(`t`)
			robotgo.KeySleep = 10

		}
	}

}
