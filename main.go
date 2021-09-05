package main

import (
	"RTC-8564/RTC-8564"
	"fmt"
	"machine"
	"time"
)

func RtcDataRecv(ch chan<- int32 ,rtc RTC_8564.Device) {

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{BaudRate: 9600})

	rtc.RtcInit()

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

func blinked(ch chan<- int32) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		ch <- 1
		led.High()
		time.Sleep(time.Millisecond * 500)
		led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
func main() {

	machine.I2C0.Configure(machine.I2CConfig{})
	rtc := RTC_8564.New(machine.I2C0)

	rtcCh2 := make(chan int32, 1)
	go RtcDataRecv(rtcCh2,rtc)

	chan2 := make(chan int32, 1)
	go blinked(chan2)

	for {
		select {
		/* RTC処理 */
		case <-rtcCh2:
			break

		/* LED点灯処理 */
		case <-chan2:
			break
		}

	}
}
