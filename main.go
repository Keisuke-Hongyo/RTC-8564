package main

import (
	"RTC-8564/RTC-8564"
	"fmt"
	"machine"
	"time"
)

func main() {

	/* アラーム出力確認 */
	init := machine.D2
	init.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	/* 表示用LEDの設定 */
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	/* RTCの設定 */
	machine.I2C0.Configure(machine.I2CConfig{})
	rtc := RTC_8564.New(machine.I2C0)

	/* UARTの設定 */
	uart := machine.UART1
	uart.Configure(machine.UARTConfig{BaudRate: 9600})

	/*RTCの初期設定 */
	rtc.RtcInit(21, 9, 5, 12, 55, 0)

	/* アラームの設定 */
	rtc.SetAlarm(12, 56)

	/* アラームスタート*/
	rtc.StartAlarm()

	for {
		/* RTCからデータを取得 */
		rtcData, dt := rtc.GetRtc()

		/* UARTでRTC-8546のデータを出力 */
		uart.Write(dt)

		/* 標準出力でRTCのデータを出力 */
		s := fmt.Sprintf("%d %d %d %d %d %d %d", rtcData.Year,
			rtcData.Month,
			rtcData.Day,
			rtcData.Wday,
			rtcData.Hour,
			rtcData.Min,
			rtcData.Sec)

		fmt.Println(s)

		/* アラームの確認 */
		if init.Get() == false {
			/* アラームON */
			led.Low()
			time.Sleep(time.Millisecond * 100)

			/* アラームリセット */
			rtc.ResetAlarm()

		} else {
			/* アラームOFF */
			led.High()
		}

		time.Sleep(time.Millisecond * 1000)

	}
}
