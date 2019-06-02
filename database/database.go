package database

import (
	"GoTask/model"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var db *sql.DB

//Init function defines a database connection
func Init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=080519 dbname=GoTaskDatabase sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// SelectAll returns a list of all users
func SelectAll() ([]model.UserInfo, error) {
	rows, err := db.Query("select * from users_info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	CurUsers := []model.UserInfo{}

	for rows.Next() {
		record := model.UserInfo{}

		err := rows.Scan(&record.ID, &record.Firstname, &record.Secondname, &record.PhoneNumber, &record.Email)
		if err != nil {
			return nil, err
		}

		CurUsers = append(CurUsers, record)
	}

	for _, record := range CurUsers {
		fmt.Println(record.ID, record.Firstname, record.Secondname, record.PhoneNumber, record.Email)
	}

	return CurUsers, nil
}

// SelectBySecondName returns the user selected by name
func SelectBySecondName(secondname string) (model.UserInfo, error) {
	record := model.UserInfo{}

	rows := db.QueryRow("select * from users_info where secondname = $1", secondname)
	err := rows.Scan(&record.ID, &record.Firstname, &record.Secondname, &record.PhoneNumber, &record.Email)

	if err != nil {
		return record, errors.New("no such record")
	}

	fmt.Println("Get by secondname:", record.ID, record.Firstname, record.Secondname,
		record.PhoneNumber, record.Email)

	return record, nil
}

// AddNewUser adds a new user
func AddNewUser(newUser model.UserInfo) error {
	_, err := db.Exec("insert into users_info (id, firstname,secondname,phone_number,email) values (default, $1, $2, $3, $4)",
		newUser.Firstname, newUser.Secondname, newUser.PhoneNumber, newUser.Email)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserByID removes the user by id
func DeleteUserByID(id int) error {
	_, err := db.Exec("delete from users_info where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserDB updates user record
func UpdateUserDB(id int, record model.UserInfo) error {
	_, err := db.Exec("update users_info set firstname = $1, secondname = $2, phone_number = $3, email = $4 where id = $5",
		record.Firstname, record.Secondname, record.PhoneNumber, record.Email, id)
	if err != nil {
		return err
	}

	return nil
}

// AddNewAccount creates a new user account
func AddNewAccount(account model.PersonalAccount) error {
	_, err := db.Exec("insert into personal_accounts (id_user, balance) values ($1, $2)",
		account.IDUser, account.Balance)
	if err != nil {
		return err
	}

	return nil
}

// GetAccount returns account by id
func GetAccount(id int) (model.PersonalAccount, error) {
	var account model.PersonalAccount

	rows := db.QueryRow("select * from personal_accounts where id = $1", id)

	err := rows.Scan(&account.ID, &account.IDUser, &account.Balance)
	if err != nil {
		return account, err
	}

	return account, err
}

// CloseAccountByID closes the user account by id
func CloseAccountByID(id int) error {
	_, err := db.Exec("delete from personal_accounts where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// ChangeBalanceDB function that allows you to change the account balance,
// conducting various transactions
//
// Transaction types:
//		0 - withdraw money
//		1 - put money
//		2 - transfer money to another account
//
func ChangeBalanceDB(id int, record model.Transaction, saveIt bool) error {
	var (
		account       model.PersonalAccount
		targetAccount model.PersonalAccount

		balance       int
		targetBalance int
		targetID      int
		value         int
	)

	rows := db.QueryRow("select * from personal_accounts where id = $1", id)

	err := rows.Scan(&account.ID, &account.IDUser, &account.Balance)
	if err != nil {
		return err
	}

	if record.TransType == "2" {
		targetID, err = strconv.Atoi(record.TargetAccID)
		if err != nil {
			return err
		}

		rows := db.QueryRow("select * from personal_accounts where id = $1", targetID)

		err = rows.Scan(&targetAccount.ID, &targetAccount.IDUser, &targetAccount.Balance)
		if err != nil {
			return err
		}

		targetBalance, err = strconv.Atoi(targetAccount.Balance)
		if err != nil {
			return err
		}
	} else {
		record.TargetAccID = "0"
	}

	balance, err = strconv.Atoi(account.Balance)
	if err != nil {
		return err
	}

	value, err = strconv.Atoi(record.Value)
	if err != nil {
		return err
	}

	/* Withdraw money/or transaction to another acc */
	if record.TransType == "0" || record.TransType == "2" {
		if balance < value {
			return errors.New("not enough money")
		}
		balance -= value
	} else {
		balance += value
	}

	if record.TransType == "2" {
		_, err = db.Exec("update personal_accounts set balance = $1 where id = $2",
			targetBalance+value, targetID)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("update personal_accounts set balance = $1 where id = $2",
		balance, id)
	if err != nil {
		return err
	}

	if saveIt {
		date, err := time.Parse("2006-01-02", record.Date)
		if err != nil {
			return err
		}

		_, err = db.Exec("insert into transactions (id_pers_acc, trans_type, date, target_acc_id, value) values ($1, $2, $3, $4, $5)",
			record.IDPersAcc, record.TransType, date, record.TargetAccID, record.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeclineTransactionDB cancels the transaction
// The function performs a reverse transaction, calling ChangeBalanceDB,
// but by injecting recipients
func DeclineTransactionDB(id int) error {
	var record model.Transaction

	rows := db.QueryRow("select * from transactions where id = $1", id)

	err := rows.Scan(&record.ID, &record.IDPersAcc, &record.TransType, &record.Date, &record.TargetAccID, &record.Value)
	if err != nil {
		return err
	}

	if record.TransType == "0" {
		record.TransType = "1"
	} else if record.TransType == "1" {
		record.TransType = "0"
	} else if record.TransType == "2" {
		targetID := record.TargetAccID
		record.TargetAccID = record.IDPersAcc
		record.IDPersAcc = targetID
	}

	persID, err := strconv.Atoi(record.IDPersAcc)
	if err != nil {
		return err
	}

	err = ChangeBalanceDB(persID, record, false)
	if err != nil {
		return err
	}

	_, err = db.Exec("delete from transactions where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionsByDate возвращает транзакции, произведенные в указанный промежуток времени
func GetTransactionsByDate(date model.RangeDate) ([]model.Transaction, error) {
	rows, err := db.Query("select * from transactions where date between $1 and $2", date.StartDate, date.FinishDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []model.Transaction{}

	for rows.Next() {
		record := model.Transaction{}

		err := rows.Scan(&record.ID, &record.IDPersAcc, &record.TransType, &record.Date, &record.TargetAccID, &record.Value)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, record)
	}

	return transactions, nil
}
