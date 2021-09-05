package main

import (
	"RTC-8564/RTC-8564"
	"fmt"
	"machine"
	"time"
)

func RtcDataRecv(ch chan<- int32, rtc RTC_8564.Device) {

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{BaudRate: 9600})

	for {
		ch <- 1
		rtcData, dt := rtc.GetRtc()
		uart.Write(dt)
		s := fmt.Sprintf("%d %d %d %d %d %d %d", rtcData.Year,
			rtcData.Month,
			rtcData.Day,
			rtcData.Wday,
			rtcData.Hour,
			rtcData.Min,
			rtcData.Sec)

		fmt.Println(s)
		time.Sleep(time.Millisecond * 1000)
	}
}

//func blinked(ch chan<- int32) {
//	led := machine.LED
//	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
//	for {
//		ch <- 1
//		led.High()
//		time.Sleep(time.Millisecond * 500)
//		led.Low()
//		time.Sleep(time.Millisecond * 500)
//	}
//}
func main() {

	init := machine.D2
	init.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.I2C0.Configure(machine.I2CConfig{})
	rtc := RTC_8564.New(machine.I2C0)

	rtc.RtcInit(21, 9, 5, 12, 55, 0)
	rtc.SetAlarm(12, 56)
	rtc.StartAlarm()

	rtcCh2 := make(chan int32, 1)
	go RtcDataRecv(rtcCh2, rtc)

	//chan2 := make(chan int32, 1)
	//go blinked(chan2)

	for {
		select {
		/* RTC処理 */
		case <-rtcCh2:
			break

			///* LED点灯処理 */
			//case <-chan2:
			//	break
		}

		if init.Get() == false {
			led.Low()
			time.Sleep(time.Millisecond * 100)
			rtc.ResetAlarm()

		} else {
			led.High()
		}

	}
}
