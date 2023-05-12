package douBot

import (
	"fmt"
	"strings"
	"time"

	db "github.com/brongh/go-telegram-financial-bot/db"
)

func FormatExpenses(expenses []db.Expense, total float64, month string) string {
	if len(expenses) == 0 {
		return "No expenses found"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Expenses for %s:\n=================================", month))
	for _, expense := range expenses {
		t, err := time.Parse(time.RFC3339, expense.ExpenseDate)
		if err != nil {
			fmt.Println("Error:", err)
			return "Something went wrong. Please try again"
		}

		builder.WriteString(fmt.Sprintf("\n\nTitle: %s\nAmount: $%.2f\nDate: %s\nID: %v\n\n=================================",
		expense.Title, expense.Amount, t.Format("2006-01-02"), expense.Id))
	}
	builder.WriteString(fmt.Sprintf("\n\nTotal spending in %v: $%.2f", month ,total))

	return builder.String()
}