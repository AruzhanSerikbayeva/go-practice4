package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// struct
type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	// .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// db
	connectionStr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// add new user
	u := User{Name: "Aruzhan", Email: "aruzhan@example.com", Balance: 1000000}
	if err := InsertUser(db, u); err != nil {
		log.Println("Insert user error:", err)
	}

	// all users
	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("Fetch all users error:", err)
	}
	fmt.Println("All users:", users)

	// by id
	u1, err := GetUserByID(db, 1)
	if err != nil {
		log.Println("Fetch user by ID error:", err)
	}
	fmt.Printf("User with ID %v: %+v\n", 1, u1)

	// транзакция
	if err := TransferBalance(db, 1, 2, 100); err != nil {
		log.Println("Transfer error:", err)
	} else {
		fmt.Println("Transfer completed successfully")
	}
}

// add user
func InsertUser(db *sqlx.DB, u User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, u)
	return err
}

// all users
func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

// by ID
func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}


func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Списываем мани
	_, err = tx.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Зачисляем мани
	_, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
