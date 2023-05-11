package douBot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	currentTime := time.Now()
	customFormat := "2006-Jan"
	currentMonth := currentTime.Format(customFormat)
	previousMonth := currentTime.AddDate(0, -1, 0).Format(customFormat)
	doublePreviousMonth := currentTime.AddDate(0, -2, 0).Format(customFormat)

	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(currentMonth, fmt.Sprintf("/view_expenses %s", currentMonth)),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(previousMonth, fmt.Sprintf("/view_expenses %s", previousMonth)),
	)
	row3 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(doublePreviousMonth, fmt.Sprintf("/view_expenses %s", doublePreviousMonth)),
	)

	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, row1, row2, row3)

	return inlineKeyboard
}