package douBot

import (
	"fmt"
	"strings"

	db "github.com/brongh/go-telegram-financial-bot/db"
)

func FormatExpenses(expenses []db.Expense, month string) string {
	if len(expenses) == 0 {
		return "No expenses found"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Expenses for %s:\n=================================", month))
	for _, expense := range expenses {
		builder.WriteString(fmt.Sprintf("\n\nTitle: %s\nAmount: %.2f\nDate: %s\nID: %v\n\n=================================",
		expense.Title, expense.Amount, expense.ExpenseDate, expense.Id))
	}

	return builder.String()
}