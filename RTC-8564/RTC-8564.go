package RTC_8564

import (
	"tinygo.org/x/drivers"
)

// RTC 機器情報構造体
type Device struct {
	bus        drivers.I2C
	rtcAddress uint8
}

// RTC日付・時刻格納構造体
type Rtc struct {
	Year  uint8
	Month uint8
	Day   uint8
	Wday  uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}

// BCDコードから10進数へ変換

// 初期化関数
func New(bus drivers.I2C) Device {
	return Device{
		bus:        bus,
		rtcAddress: rtcAddress,
	}
}

// RTCの初期化関数
func (d *Device) RtcInit(y uint8, mth uint8, dy uint8, h uint8, m uint8, s uint8) {
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl1, []byte{rtc_8564_stop})
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl2, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, ctrltm, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, ckout, []byte{0x83}) // 7:FE=1 10:FD=11 (1Hz)
	_ = d.bus.WriteRegister(d.rtcAddress, year, dectobcd(y))   // 年 (下位2桁)
	_ = d.bus.WriteRegister(d.rtcAddress, month, dectobcd(mth))
	_ = d.bus.WriteRegister(d.rtcAddress, day, dectobcd(dy))
	_ = d.bus.WriteRegister(d.rtcAddress, hour, dectobcd(h))
	_ = d.bus.WriteRegister(d.rtcAddress, min, dectobcd(m))
	_ = d.bus.WriteRegister(d.rtcAddress, sec, dectobcd(s))
	_ = d.bus.WriteRegister(d.rtcAddress, wday, getWeekday(y, m, dy)) // 0:日 1:月 2:火 3:水 4:木 5:金 6:土

	// Alarm Setting -> OFF
	_ = d.bus.WriteRegister(d.rtcAddress, min_alarm, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, hour_alarm, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, day_alarm, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, wday_alarm, []byte{0x00})

	// RTC Start
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl1, []byte{rtc_8564_start}) // RTC Start
}

func (d *Device) SetAlarm(hour uint8, min uint8) {
	_ = d.bus.WriteRegister(d.rtcAddress, min_alarm, dectobcd(min))
	_ = d.bus.WriteRegister(d.rtcAddress, hour_alarm, dectobcd(hour))
	_ = d.bus.WriteRegister(d.rtcAddress, day_alarm, []byte{0x80})
	_ = d.bus.WriteRegister(d.rtcAddress, wday_alarm, []byte{0x80})
}

func (d *Device) StartAlarm() {
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl2, []byte{0x02})
}

func (d *Device) StopAlarm() {
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl2, []byte{0x00})
}

func (d *Device) ResetAlarm() {
	data := make([]byte, 1)
	_ = d.bus.ReadRegister(d.rtcAddress, ctrl2, data)
	data[0] = data[0] & 0xf7
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl2, data)
}

// Zellerの公式による曜日の計算 0:日 1:月 2:火 3:水 4:木 5:金 6:土
func getWeekday(year uint8, month uint8, day uint8) []byte {

	var y uint16

	y = 2000 + uint16(year)

	wd := make([]byte, 1)
	if month <= 2 {
		month += 12
		year--
	}
	wd[0] = byte(((y + y/4 - y/100 + y/400 + ((13*uint16(month) + 8) / 5) + uint16(day)) % 7) & 0x00ff)

	return wd
}

// 時刻データの取得
func (d *Device) GetRtc() (Rtc, []byte) {

	rtc := Rtc{}
	data := make([]byte, 7)
	_ = d.bus.ReadRegister(d.rtcAddress, sec, data)
	data[5] = data[5] & 0x1f
	data[3] = data[3] & 0x3f
	data[4] = data[4] & 0x07
	data[2] = data[2] & 0x3f
	data[1] = data[1] & 0x7f
	data[0] = data[0] & 0x7f

	rtc.Year = bcdtodec(data[6])
	rtc.Month = bcdtodec(data[5])
	rtc.Day = bcdtodec(data[3])
	rtc.Wday = bcdtodec(data[4])
	rtc.Hour = bcdtodec(data[2])
	rtc.Min = bcdtodec(data[1])
	rtc.Sec = bcdtodec(data[0])

	return rtc, data
}

// BCDから10進数へ変換
func bcdtodec(b byte) uint8 {
	return ((b >> 4) * 10) + (b & 0x0f)
}

// 10進数をBCDに変換

func dectobcd(d uint8) []byte {
	bcd := make([]byte, 1)
	bcd[0] = (((d / 10) & 0x0f) << 4) + ((d % 10) & 0x0f)
	return bcd
}
