package main

import (
	"fmt"
	"time"
)

const (
	daysInWeek int = 7
	curWeekQuestion string = "Могу играть на этой неделе в:"
	nextWeekQuestion string = "Могу играть на следующей неделе в:"
	neverOption string = "не могу в эти даты"
	day time.Duration = time.Hour * 24
)

var dayOfWeek = map[int]string{
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
	7: "Воскресенье",
}

var month = map[time.Month]string{
	1: "января",
	2: "февраля",
	3: "марта",
	4: "апреля",
	5: "мая",
	6: "июня",
	7: "июля",
	8: "августа",
	9: "сентября",
	10: "октября",
	11: "ноября",
	12: "декабря",
}

func getCurWeekPollParams() []string {
	return getOptionsFromDateTillEndOfWeek(time.Now())
}

func getNextWeekPollParams() []string {
	return getOptionsFromDateTillEndOfWeek(getNextMonday(time.Now()))
}

func getOptionsFromDateTillEndOfWeek(from time.Time) []string {
	cur := from

	weekday := getWeekday(cur)

	options := make([]string, 0, daysInWeek - weekday + 1)

	for i := weekday; i <= daysInWeek; i++ {
		s := fmt.Sprintf("%d %s, %s", cur.Day(), month[cur.Month()], dayOfWeek[i])
		options = append(options, s)
		cur = cur.Add(day)
	}

	options = append(options, neverOption)

	return options
}

func getWeekday(now time.Time) int {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday
}

func getNextMonday(now time.Time) time.Time {
	return now.Add(day * time.Duration(8 - getWeekday(now)))
}