package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"../mouse"
	"github.com/tajtiattila/xinput"
)

func main() {

	if err := xinput.Load(); err != nil {
		panic(err)
	}

	state := &xinput.State{}
	vibration := &xinput.Vibration{}

	lmbDown := false
	rmbDown := false
	lmbTurboDown := false
	ch := make(chan int, 1)
	for {
		if err := xinput.GetState(0, state); err != nil {
			panic(fmt.Errorf("Controller 0 error: %w", err))
		}

		var point *mouse.POINT
		if math.Abs(float64(state.Gamepad.ThumbLX)) > xinput.LEFT_THUMB_DEADZONE || math.Abs(float64(state.Gamepad.ThumbLY)) > xinput.LEFT_THUMB_DEADZONE {
			point = handleLStick(state)
		} else {
			point = handleDpad(state)
		}
		if math.Abs(float64(state.Gamepad.ThumbRY)) > xinput.RIGHT_THUMB_DEADZONE {
			handleScrolling(state)
		}

		// Handle vibration
		vibration.LeftMotorSpeed = uint16(state.Gamepad.LeftTrigger) * 257
		vibration.RightMotorSpeed = uint16(state.Gamepad.RightTrigger) * 257
		xinput.SetState(0, vibration)

		// LMB
		if lmbDown == false && (xinput.BUTTON_X&xinput.Button(state.Gamepad.Buttons) == xinput.BUTTON_X) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTDOWN, point)
			lmbDown = true
		}
		if lmbDown == true && (xinput.BUTTON_X&xinput.Button(state.Gamepad.Buttons) != xinput.BUTTON_X) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTUP, point)
			lmbDown = false
		}

		// RMB
		if rmbDown == false && (xinput.BUTTON_B&xinput.Button(state.Gamepad.Buttons) == xinput.BUTTON_B) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_RIGHTDOWN, point)
			rmbDown = true
		}
		if rmbDown == true && (xinput.BUTTON_B&xinput.Button(state.Gamepad.Buttons) != xinput.BUTTON_B) {
			mouse.MouseEvent(mouse.MOUSEEVENTF_RIGHTUP, point)
			rmbDown = false
		}

		if lmbTurboDown == false && (xinput.BUTTON_Y&xinput.Button(state.Gamepad.Buttons) == xinput.BUTTON_Y) {
			go func(done chan int, point *mouse.POINT) {
				for {

					select {

					case <-done:
						return
					default:
						mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTDOWN, point)
						time.Sleep(time.Millisecond * 25)
						mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTUP, point)
					}
				}
			}(ch, point)
			// Start turbo clicking
			//mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTDOWN, point)
			lmbTurboDown = true
		}
		if lmbTurboDown == true && (xinput.BUTTON_Y&xinput.Button(state.Gamepad.Buttons) != xinput.BUTTON_Y) {
			ch <- 1
			// Stop turbo clicking
			//mouse.MouseEvent(mouse.MOUSEEVENTF_LEFTDOWN, point)
			lmbTurboDown = false
		}

		if xinput.BACK&xinput.Button(state.Gamepad.Buttons) == xinput.BACK {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		mouse.SetCursorPos(point.X, point.Y)
		time.Sleep(time.Second / 60)
	}

}
func handleScrolling(state *xinput.State) {
	point := &mouse.POINT{}
	mouse.GetCursorPos(point)
	speedY := int32(state.Gamepad.ThumbRY) / 800

	mouse.Scroll(mouse.DWORD(speedY))
}
func handleLStick(state *xinput.State) *mouse.POINT {
	point := &mouse.POINT{}
	mouse.GetCursorPos(point)
	fmt.Println(state.Gamepad)
	speedX := int32(state.Gamepad.ThumbLX) / 1600
	speedY := int32(state.Gamepad.ThumbLY) / 1600

	point.X = point.X + speedX
	point.Y = point.Y - speedY
	return point
}

func handleDpad(state *xinput.State) *mouse.POINT {
	const speed = 10
	point := &mouse.POINT{}
	mouse.GetCursorPos(point)
	if xinput.DPAD_UP&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_UP {
		point.Y -= speed
	}
	if xinput.DPAD_DOWN&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_DOWN {
		point.Y += speed
	}
	if xinput.DPAD_LEFT&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_LEFT {
		point.X -= speed
	}
	if xinput.DPAD_RIGHT&xinput.Button(state.Gamepad.Buttons) == xinput.DPAD_RIGHT {
		point.X += speed
	}
	return point
}
