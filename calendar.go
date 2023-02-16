package main

import (
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type calendarKeyboardState int

const (
	yearKeyboard calendarKeyboardState = iota
	monthKeyboard
	dayKeyboard
)

var months = map[time.Month]string {
	time.January: "Январь",
	time.February: "Февраль",
	time.March: "Март",
	time.April: "Апрель",
	time.May: "Май",
	time.June: "Июнь",
	time.July: "Июль",
	time.August: "Август",
	time.September: "Сентябрь",
	time.October: "Октябрь",
	time.November: "Ноябрь",
	time.December: "Декабрь",
} 

type calendarInfo struct {
	state calendarKeyboardState
	year, month, day int
}

var calendarCache = make(map[int]*calendarInfo)

var defaultMonthKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(months[time.January], "1"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.February], "2"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.March], "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(months[time.April], "4"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.May], "5"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.June], "6"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(months[time.July], "7"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.August], "8"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.September], "9"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(months[time.October], "10"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.November], "11"),
		tgbotapi.NewInlineKeyboardButtonData(months[time.December], "12"),
	),
)

var header = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Пн", " "),
	tgbotapi.NewInlineKeyboardButtonData("Вт", " "),
	tgbotapi.NewInlineKeyboardButtonData("Cр", " "),
	tgbotapi.NewInlineKeyboardButtonData("Чт", " "),
	tgbotapi.NewInlineKeyboardButtonData("Пт", " "),
	tgbotapi.NewInlineKeyboardButtonData("Cб", " "),
	tgbotapi.NewInlineKeyboardButtonData("Вс", " "),
)

func newCalendarKeyboard(status calendarKeyboardState, t time.Time) tgbotapi.InlineKeyboardMarkup {
	switch status {
	case yearKeyboard:
		return newPickYearCalendarKeyboard(t)
	case monthKeyboard:
		return newPickMonthCalendarKeyboard(t)
	case dayKeyboard:
		return newPickDayCalendarKeyboard(t)
	default:
		return tgbotapi.InlineKeyboardMarkup{}
	}
}

func newPickYearCalendarKeyboard(t time.Time) tgbotapi.InlineKeyboardMarkup {
	curYear := t.Year()
	firstRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(curYear - 2), strconv.Itoa(curYear - 2)),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(curYear - 1), strconv.Itoa(curYear - 1)),
	)
	secondRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(string(highlight(strconv.Itoa(curYear))), strconv.Itoa(curYear)),
	)
	thirdRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(curYear + 1), strconv.Itoa(curYear + 1)),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(curYear + 2), strconv.Itoa(curYear + 2)),
	)
	return tgbotapi.NewInlineKeyboardMarkup(firstRow, secondRow, thirdRow)
}

func newPickMonthCalendarKeyboard(t time.Time) tgbotapi.InlineKeyboardMarkup {
	m := int(time.Now().Month())
	row := (m - 1) / 3
	col := (m - 1) % 3
	res := tgbotapi.InlineKeyboardMarkup {
		InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, len(defaultMonthKeyboard.InlineKeyboard)),
	}
	for i := range defaultMonthKeyboard.InlineKeyboard {
		res.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(defaultMonthKeyboard.InlineKeyboard[i]))
		copy(res.InlineKeyboard[i], defaultMonthKeyboard.InlineKeyboard[i])
	}
	if (t.Year() == time.Now().Year()) {
		res.InlineKeyboard[row][col].Text = highlight(res.InlineKeyboard[row][col].Text)
	}
	return res
}

func newPickDayCalendarKeyboard(t time.Time) tgbotapi.InlineKeyboardMarkup {
	cur := getFirstDayOfMonth(t)
	now := time.Now()
	curMonth := cur.Month()
	emptyButton := tgbotapi.NewInlineKeyboardButtonData(" ", " ")
	curRow := make([]tgbotapi.InlineKeyboardButton, daysInWeek)
	rowsCount := 0
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 7)
	rows = append(rows, header)
	rowsCount++
	for i := 0; i < daysInWeek; i++ {
		if i < getWeekday(cur) - 1 {
			curRow[i] = emptyButton
		} else {
			var text string
			if (cur.Year() == now.Year() && curMonth == now.Month() && cur.Day() == now.Day()) {
				text = highlight(strconv.Itoa(cur.Day()))
			} else {
				text = strconv.Itoa(cur.Day())
			}
			data := strconv.Itoa(cur.Day())
			curRow[i] = tgbotapi.NewInlineKeyboardButtonData(text, data)
			cur = cur.Add(day)
		}
	}
	rows = append(rows, curRow)
	curRow = make([]tgbotapi.InlineKeyboardButton, daysInWeek)
	rowsCount++
	count := 0
	for cur.Month() == curMonth {
		var text string
		if (cur.Year() == now.Year() && curMonth == now.Month() && cur.Day() == now.Day()) {
			text = highlight(strconv.Itoa(cur.Day()))
		} else {
			text = strconv.Itoa(cur.Day())
		}
		data := strconv.Itoa(cur.Day())
		curRow[count] = tgbotapi.NewInlineKeyboardButtonData(text, data)
		cur = cur.Add(day)
		if (count == 6) {
			rows = append(rows, curRow)
			rowsCount++
			curRow = make([]tgbotapi.InlineKeyboardButton, daysInWeek)
			count = 0
		} else {
			count++
		}
	}
	if count != 0 {
		for ; count < daysInWeek; count++ {
			curRow[count] = emptyButton
		}
		rows = append(rows, curRow)
		rowsCount++
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows[:rowsCount]...)
}

func highlight(s string) string {
	highlighter := string('❗')
	return highlighter + s + highlighter
}

func getFirstDayOfMonth(t time.Time) time.Time {
	return t.Add(-day * time.Duration(t.Day() - 1))
}