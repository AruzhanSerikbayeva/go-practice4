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

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	connectionStr := os.Getenv("DATABASE_URL")

	db, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Insert a user
	u := User{Name: "Aruzhan", Email: "aruzhan@example.com", Balance: 1000000}
	if err := InsertUser(db, u); err != nil {
		log.Println("Insert user error:", err)
	}

	// Fetch all users
	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("Fetch all users error:", err)
	}
	fmt.Println("Users:", users)

	// User by ID
	u1, err := GetUserByID(db, 1)
	if err != nil {
		log.Println("Fetch user by id error:", err)
	}
	fmt.Printf("User with id %v: %v", 1, u1)

	// Transaction
	if err := TransferBalance(db, 1, 4, 100); err != nil {
		log.Println("Transfer error:", err)
	} else {
		fmt.Println("Transfer completed successfully")
	}
}

func InsertUser(db *sqlx.DB, u User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, u)
	return err
}
func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}
func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}
func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	t, err := db.Beginx()
	if err != nil {
		return err
	}

	// Take money
	_, err = t.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, fromID)
	if err != nil {
		t.Rollback()
		return err
	}

	// Give money
	_, err = t.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		t.Rollback()
		return err
	}
	return t.Commit()
}
