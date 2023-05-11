package douBot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	db "github.com/brongh/go-telegram-financial-bot/db"
	utils "github.com/brongh/go-telegram-financial-bot/utils"
)

var bot *tgbotapi.BotAPI

type chatParams struct {
	messageID     int
	chatID        int64
	userID        int
	msgText       string
	username 	string
}

func StartBot() {
	BotInit()
	db.DbInit()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("[-] Failed to get updates from Telegram server")
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.CallbackQuery != nil {
			go handleCallbacks(update)
		} else {
			go onMessage(update)
		}
		
	}
}

func BotInit() {
	initConfig := utils.ReadConfig()
	
	var err error
	bot, err = tgbotapi.NewBotAPI(initConfig.TgToken)
	if err != nil {
		log.Println("[-] Login failed, please check your token")
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("[+] Authorized on account %s\n\n", bot.Self.UserName)
}



func onMessage(update tgbotapi.Update) {
	var chat chatParams

	chat.chatID = update.Message.Chat.ID
	chat.userID = update.Message.From.ID
	chat.messageID = update.Message.MessageID
	chat.msgText = strings.ToLower(update.Message.Text)
	chat.username = update.Message.From.UserName
	
	log.Print("[**] Received msg from userID: ", chat.chatID)

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "view":
			msg := tgbotapi.NewMessage(chat.chatID, "Select which month's expenses you would like to view: ")
			msg.ReplyMarkup = CreateInlineKeyboard()
			bot.Send(msg)
			// handle commands in callback
		default:
			msg := tgbotapi.NewMessage(chat.chatID, "Sorry, I didn't understand that command.")
			bot.Send(msg)
		}
	} else {
		var msg tgbotapi.MessageConfig
		inputStrings := strings.Split(chat.msgText, " ")
		actionPayload := inputStrings[1:]
		switch inputStrings[0] {
		case "add":
			if len(inputStrings) == 1 {
				msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid command. e.g. add <title> <amount w/o $> <date YYYY-MM-DD (optional)>\n\nadd food 13.3")
			} else {
				AddExpense(actionPayload, chat)
			}
		case "view":
			msg = tgbotapi.NewMessage(chat.chatID, "type /view to view expenses")
		case "delete":
			if len(inputStrings) == 1 {
				msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid command. e.g. delete <expenseID>\n\ndelete 1")
			} else {
				expenseID, err := strconv.Atoi(inputStrings[1])
			if err != nil {
				msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid command. e.g. delete <expenseID>\n\ndelete 1")
				bot.Send(msg)
				return
			}
			success, err := DeleteExpense(chat.userID, expenseID)
			if err != nil {
				msg = tgbotapi.NewMessage(chat.chatID, "Sorry, something went wrong. Please try again later.")
				bot.Send(msg)
				return
			}
			if success {
				msg = tgbotapi.NewMessage(chat.chatID, fmt.Sprintf("Successfully deleted expense with ID: %v", expenseID))
			}
			}
			
		default:
			msg = tgbotapi.NewMessage(chat.chatID, "Invalid command. These are the available commands:\n\nadd --> add an expense\n/view --> view expenses\ndelete --> delete an expense")
		}
		bot.Send(msg)
	}
}

func AddExpense(actionPayload []string, chat chatParams) {
	var expense db.Expense
	var msg tgbotapi.MessageConfig
	lengthOfPayload := len(actionPayload)
	userData, err := db.GetOrCreateUser(chat.username, chat.userID)
	if userData != nil {
		fmt.Print("New user created. userId: ", userData)
	}
	if err != nil {
		msg = tgbotapi.NewMessage(chat.chatID, "Sorry, something went wrong. Please try again later.")
		bot.Send(msg)
		return
	}
	expense.UserId = userData.Id
	expense.Title = actionPayload[0]
	switch lengthOfPayload {
	case 2:
		// date as today
		floatValue, err := strconv.ParseFloat(actionPayload[1], 64)
		if err != nil {
			msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid amount. e.g. 13, 3.33 or 4.1")
			bot.Send(msg)
			return
		}
		expense.Amount = floatValue
		expense.ExpenseDate = time.Now().Format("2006-01-02")
		err = db.CreateExpense(expense.UserId, expense.Title, expense.Amount, expense.ExpenseDate)
		if err != nil {
			msg = tgbotapi.NewMessage(chat.chatID, "Sorry, something went wrong. Please try again later.")
			bot.Send(msg)
			return
		}
		msg = tgbotapi.NewMessage(chat.chatID, fmt.Sprintf("Added expense: %s, %.2f, %s for userID:%v", expense.Title, expense.Amount, expense.ExpenseDate, expense.UserId))
	case 3:
		// date as specified
		floatValue, err := strconv.ParseFloat(actionPayload[1], 64)
		if err != nil {
			msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid amount. e.g. 13, 3.33 or 4.1")
			bot.Send(msg)
			return
		}
		expense.Amount = floatValue
		expenseDate, err := time.Parse("2006-01-02", actionPayload[2])
		if err != nil {
			msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid date. e.g. 2020-01-01 YYYY-MM-DD")
			bot.Send(msg)
			return
		}
		expense.ExpenseDate = expenseDate.Format("2006-01-02")
		err = db.CreateExpense(expense.UserId, expense.Title, expense.Amount, expense.ExpenseDate)
		if err != nil {
			msg = tgbotapi.NewMessage(chat.chatID, "Sorry, something went wrong. Please try again later.")
			bot.Send(msg)
			return
		}
		msg = tgbotapi.NewMessage(chat.chatID, fmt.Sprintf("Added expense: %s, %.2f, %s for userID:%v", expense.Title, expense.Amount, expense.ExpenseDate, expense.UserId))
	default: 
		msg = tgbotapi.NewMessage(chat.chatID, "Please enter a valid command. e.g. add <title> <amount w/o $> <date YYYY-MM-DD (optional)>\n\nadd food 13.3")
	}

	bot.Send(msg)
}

func handleCallbacks(update tgbotapi.Update) {
	callback := update.CallbackQuery
	callBackID := callback.Message.Chat.ID
	if update.CallbackQuery != nil {
		// handle callback queries
		// callback queries are triggered by inline keyboards
		var responseText string
		fmt.Println("\n=============================: ", update.CallbackQuery.Data)
		inputStrings := strings.Split(update.CallbackQuery.Data, " ")
		if len(inputStrings) != 2 && len(inputStrings) != 1 {
			msg := tgbotapi.NewMessage(callBackID, "Sorry, something went wrong. Please try again later.")
			bot.Send(msg)
			return
		}
		switch inputStrings[0] {
		case "/view_expenses":
			fmt.Println("VIEW: ", inputStrings)
			monthSplit := strings.Split(inputStrings[1], "-")
			if len(monthSplit) != 2 {
				fmt.Print("inputStrings error")
				msg := tgbotapi.NewMessage(callBackID, "Sorry, something went wrong. Please try again later.")
				bot.Send(msg)
				return
			}
			var year, month = monthSplit[0], monthSplit[1]
			yearInt, err := strconv.Atoi(year)
			if err != nil {
				fmt.Print("int conversion error: ", err)
				msg := tgbotapi.NewMessage(callBackID, "Sorry, something went wrong. Please try again later.")
				bot.Send(msg)
				return
			}
			parsedTime, err := time.Parse("Jan", month)
			if err != nil {
				fmt.Print("time conversion for month failed: ", err)
				msg := tgbotapi.NewMessage(callBackID, "Sorry, something went wrong. Please try again later.")
				bot.Send(msg)
				return
			}
			monthInt := int(parsedTime.Month())
			expense, err := db.ViewExpenses(int(callBackID), time.Month(monthInt), yearInt)
			if err != nil {
				fmt.Print("db error: ", err)
				msg := tgbotapi.NewMessage(callBackID, "Sorry, something went wrong. Please try again later.")
				bot.Send(msg)
				return
			}
			responseText = FormatExpenses(expense, month)
		default:
			responseText = "Sorry, something went wrong. Please try again later."
		}
		msg := tgbotapi.NewMessage(callBackID, responseText)
		bot.Send(msg)
		bot.AnswerCallbackQuery((tgbotapi.NewCallback(update.CallbackQuery.ID, "")))
	}
}

func DeleteExpense(tgId int, expenseId int)	( bool, error) {
	return db.RemoveExpense(tgId, expenseId)
}
