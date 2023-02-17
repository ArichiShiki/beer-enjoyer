package timeHelper

import (
	"fmt"
	"log"
	"time"
)

const (
	DaysInWeek int           = 7
	Day        time.Duration = time.Hour * 24
)

var DayOfWeek = map[int]string{
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
	7: "Воскресенье",
}

var MonthsAsNominative = map[time.Month]string{
	time.January:   "Январь",
	time.February:  "Февраль",
	time.March:     "Март",
	time.April:     "Апрель",
	time.May:       "Май",
	time.June:      "Июнь",
	time.July:      "Июль",
	time.August:    "Август",
	time.September: "Сентябрь",
	time.October:   "Октябрь",
	time.November:  "Ноябрь",
	time.December:  "Декабрь",
}

var MonthAsGenitive = map[time.Month]string{
	time.January:   "января",
	time.February:  "февраля",
	time.March:     "марта",
	time.April:     "апреля",
	time.May:       "мая",
	time.June:      "июня",
	time.July:      "июля",
	time.August:    "августа",
	time.September: "сентября",
	time.October:   "октября",
	time.November:  "ноября",
	time.December:  "декабря",
}

func Weekday(now time.Time) int {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday
}

func NextMonday(now time.Time) time.Time {
	return now.Add(Day * time.Duration(8-Weekday(now)))
}

func FirstDayOfMonth(year int, month time.Month) time.Time {
	t, err := time.Parse("2006-January-02", fmt.Sprintf("%04d-%s-01", year, month))
	if err != nil {
		log.Panic(err)
	}
	return t
}

func DaysInMonth(year int, month time.Month) int {
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if IsLeapYear(year) {
			return 29
		} else {
			return 28
		}
	default:
		return 0
	}
}

// TODO rename
func IsLeapYear(year int) bool {
	return year%400 == 0 || year%100 != 0 && year%4 == 0
}
