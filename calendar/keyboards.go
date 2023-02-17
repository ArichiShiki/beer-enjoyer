package calendar

import (
	"fmt"
	"strconv"
	"time"

	"arichi/beerEnjoyer/timeHelper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type calendarKeyboardState int

const (
	YearKeyboard calendarKeyboardState = iota
	MonthKeyboard
	DayKeyboard
	Finish
)

const (
	EmptyData string = "empty"
	PrevData string = "prev"
	NextData string = "next"
	GoToMonthData string = "goToMonth"
	GoToYearData string = "goToYear"
)
var emptyButton = tgbotapi.NewInlineKeyboardButtonData(" ", EmptyData)
var prevButton = tgbotapi.NewInlineKeyboardButtonData("<--", PrevData)
var nextButton = tgbotapi.NewInlineKeyboardButtonData("-->", NextData)

var dayKeyboardHeader = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("Пн", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Вт", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Cр", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Чт", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Пт", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Cб", EmptyData),
	tgbotapi.NewInlineKeyboardButtonData("Вс", EmptyData),
)

var defaultDayButtons []tgbotapi.InlineKeyboardButton = []tgbotapi.InlineKeyboardButton{
	tgbotapi.NewInlineKeyboardButtonData("1", "1"),
	tgbotapi.NewInlineKeyboardButtonData("2", "2"),
	tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	tgbotapi.NewInlineKeyboardButtonData("4", "4"),
	tgbotapi.NewInlineKeyboardButtonData("5", "5"),
	tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	tgbotapi.NewInlineKeyboardButtonData("7", "7"),
	tgbotapi.NewInlineKeyboardButtonData("8", "8"),
	tgbotapi.NewInlineKeyboardButtonData("9", "9"),
	tgbotapi.NewInlineKeyboardButtonData("10", "10"),
	tgbotapi.NewInlineKeyboardButtonData("11", "11"),
	tgbotapi.NewInlineKeyboardButtonData("12", "12"),
	tgbotapi.NewInlineKeyboardButtonData("13", "13"),
	tgbotapi.NewInlineKeyboardButtonData("14", "14"),
	tgbotapi.NewInlineKeyboardButtonData("15", "15"),
	tgbotapi.NewInlineKeyboardButtonData("16", "16"),
	tgbotapi.NewInlineKeyboardButtonData("17", "17"),
	tgbotapi.NewInlineKeyboardButtonData("18", "18"),
	tgbotapi.NewInlineKeyboardButtonData("19", "19"),
	tgbotapi.NewInlineKeyboardButtonData("20", "20"),
	tgbotapi.NewInlineKeyboardButtonData("21", "21"),
	tgbotapi.NewInlineKeyboardButtonData("22", "22"),
	tgbotapi.NewInlineKeyboardButtonData("23", "23"),
	tgbotapi.NewInlineKeyboardButtonData("24", "24"),
	tgbotapi.NewInlineKeyboardButtonData("25", "25"),
	tgbotapi.NewInlineKeyboardButtonData("26", "26"),
	tgbotapi.NewInlineKeyboardButtonData("27", "27"),
	tgbotapi.NewInlineKeyboardButtonData("28", "28"),
	tgbotapi.NewInlineKeyboardButtonData("29", "29"),
	tgbotapi.NewInlineKeyboardButtonData("30", "30"),
	tgbotapi.NewInlineKeyboardButtonData("31", "31"),
}

func newPickYearCalendarKeyboard(year int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year-4)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year-3)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year-2)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year-1)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year+1)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year+2)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year+3)),
		tgbotapi.NewInlineKeyboardButtonData(yearDataButton(year+4)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(prevButton, nextButton))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func yearDataButton(year int) (string, string) {
	if year == time.Now().Year() {
		return highlight(strconv.Itoa(year)), strconv.Itoa(year)
	} else {
		return strconv.Itoa(year), strconv.Itoa(year)
	}
}

func newPickMonthCalendarKeyboard(year int) tgbotapi.InlineKeyboardMarkup {
	res := newDefaultMonthKeyboard()
	now := time.Now()
	if year == now.Year() {
		m := int(now.Month())
		row := (m - 1) / 3
		col := (m - 1) % 3
		res.InlineKeyboard[row][col].Text = highlight(res.InlineKeyboard[row][col].Text)
	}
	res.InlineKeyboard = append(res.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		prevButton,
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(year), GoToYearData),
		nextButton,
	))
	return res
}

func newDefaultMonthKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.January], "1"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.February], "2"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.March], "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.April], "4"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.May], "5"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.June], "6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.July], "7"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.August], "8"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.September], "9"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.October], "10"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.November], "11"),
			tgbotapi.NewInlineKeyboardButtonData(timeHelper.MonthsAsNominative[time.December], "12"),
		),
	)
}

func newPickDayCalendarKeyboard(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	d := timeHelper.FirstDayOfMonth(year, month)

	var buttons []tgbotapi.InlineKeyboardButton
	for i := 0; i < timeHelper.Weekday(d)-1; i++ {
		buttons = append(buttons, emptyButton)
	}
	now := time.Now()
	if now.Year() == year && now.Month() == month {
		buttons = append(buttons, newDefaultMonthButtons(year, month, now.Day())...)
	} else {
		buttons = append(buttons, newDefaultMonthButtons(year, month, 0)...)
	}
	if len(buttons) % 7 != 0 {
		remain := timeHelper.DaysInWeek - len(buttons) % 7
		for i := 0; i < remain; i++ {
			buttons = append(buttons, emptyButton)
		}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, dayKeyboardHeader)
	for i := 0; i < len(buttons); i += 7 {
		rows = append(rows, buttons[i:i+7])
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		prevButton,
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v %v", timeHelper.MonthsAsNominative[month], year), GoToMonthData),
		nextButton,
	))

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func newDefaultMonthButtons(year int, month time.Month, highlightDay int) []tgbotapi.InlineKeyboardButton {
	l := timeHelper.DaysInMonth(year, month)
	buttons := make([]tgbotapi.InlineKeyboardButton, l)
	copy(buttons, defaultDayButtons)
	if 1 <= highlightDay && highlightDay <= l {
		buttons[highlightDay-1].Text = highlight(buttons[highlightDay-1].Text)
	}
	return buttons
}

func highlight(s string) string {
	highlighter := string('❗')
	return highlighter + s + highlighter
}
