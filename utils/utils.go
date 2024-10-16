package utils

import (
	"IA04-hotel/agt/employee"
)

func CheckWorkingSchedule(emp *employee.Employee, currDay int, currHour int) bool {
	for _, d := range emp.GetSchedule() {
		if int(d) == currDay%7 {
			switch emp.GetShift() {
			case 0:
				if currHour >= 0 && currHour <= 7 {
					return true
				} else {
					return false
				}
			case 1:
				if currHour >= 8 && currHour <= 15 {
					return true
				} else {
					return false
				}
			case 2:
				if currHour >= 16 && currHour <= 23 {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}
