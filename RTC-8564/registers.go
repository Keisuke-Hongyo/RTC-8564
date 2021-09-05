package RTC_8564

const rtcAddress = 0x51
const ctrl1 	 = 0x00
const ctrl2 	 = 0x01
const sec  		 = 0x02
const min		 = 0x03
const hour		 = 0x04
const day		 = 0x05
const wday		 = 0x06
const month		 = 0x07
const year		 = 0x08
const ckout		 = 0x0d 			// CLKOUT frequency 7:FE 1:FD1 0:FD0
const ctrltm	 = 0x0e 			// Timer Control 7:TE 1:TD1 0:TD0