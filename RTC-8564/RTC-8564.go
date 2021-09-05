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
	Year  	uint8
	Month 	uint8
	Day   	uint8
	Wday  	uint8
	Hour  	uint8
	Min   	uint8
	Sec   	uint8
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
func (d *Device) RtcInit() {
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl1, []byte{0x20})
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl2, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, ctrltm, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, ckout, []byte{0x83}) // 7:FE=1 10:FD=11 (1Hz)
	_ = d.bus.WriteRegister(d.rtcAddress, year, []byte{0x20})  // 年 (下位2桁)
	_ = d.bus.WriteRegister(d.rtcAddress, month, []byte{0x09})
	_ = d.bus.WriteRegister(d.rtcAddress, day, []byte{0x05})
	_ = d.bus.WriteRegister(d.rtcAddress, hour, []byte{0x10})
	_ = d.bus.WriteRegister(d.rtcAddress, min, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, sec, []byte{0x00})
	_ = d.bus.WriteRegister(d.rtcAddress, wday, []byte{0x00})  // 0:日 1:月 2:火 3:水 4:木 5:金 6:土
	_ = d.bus.WriteRegister(d.rtcAddress, ctrl1, []byte{0x00}) // RTC Start
}

// 時刻データの取得
func (d *Device) GetRtc() (Rtc, []byte) {

	bcdtodec := func(b byte) byte {
		return ((b >> 4) * 10) + (b & 0x0f)
	}

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
