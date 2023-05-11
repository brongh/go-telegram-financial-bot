package db

import (
	"time"
)

func CreateExpense(userId int, title string, amount float64, date string) error {
	query := `
		INSERT INTO expenses (user_id, title, amount, expense_date)
		VALUES ($1, $2, $3, $4)
	`
	_, err := DB.Exec(query, userId, title, amount, date)

	if err != nil {
		return err
	}

	return nil
}

func ViewExpenses(tgId int, month time.Month, year int) ([]Expense, error) {
	user, err := GetUserByTgId(tgId)
	if err != nil {
		return nil, err
	}
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := `
		SELECT id, title, amount, expense_date
		FROM expenses
		WHERE user_id = $1 AND expense_date >= $2 AND expense_date < $3
		ORDER BY expense_date DESC
	`
	
	rows, err := DB.Query(query, user.Id, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []Expense{}
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.ExpenseDate)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func RemoveExpense(tgId int, expenseId int) (bool, error) {
	user, err := GetUserByTgId(tgId)
	if err != nil {
		return false, err
	}
	query := `
		DELETE FROM expenses
		WHERE user_id = $1 AND id = $2
	`
	_, err = DB.Exec(query, user.Id, expenseId)
	if err != nil {
		return false, err
	}
	return true, nil
}