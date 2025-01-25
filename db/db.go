package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" 
)

const (
	// Вводим свои данные 
	host     = "  "
	port     = 5432
	user     = "  "
	password = "  "
	dbname   = "  "
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	var exists bool
	err = db.QueryRow(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", dbname)).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows{
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
			if err != nil && !strings.Contains(err.Error(), "already exists") {
				log.Fatal(err)
			}
			fmt.Println("Database created successfully!")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Database already exists.")
	}


	connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")
}

func Create_table() error {
	query := `
	CREATE TABLE IF NOT EXISTS library (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		published_year INTEGER NOT NULL,
		genre VARCHAR(255) NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created successfully!")
	return nil
}

func Add_book(title, author string, published_year int64, genre string) (int64, error) {
	query := `
		INSERT INTO library (title, author, published_year, genre) VALUES ($1, $2, $3, $4) RETURNING id
	`
	var id int64
	err := db.QueryRow(query, title, author, published_year, genre).Scan(&id)
	if err != nil {
		return 0, err
	}
	fmt.Println("Book added successfully!")
	return id, nil
}

func Update_book(id int64, title, author *string, published_year *int64, genre *string) error {
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *title)
		argIndex++
	}
	if author != nil {
		setClauses = append(setClauses, fmt.Sprintf("author = $%d", argIndex))
		args = append(args, *author)
		argIndex++
	}
	if published_year != nil {
		setClauses = append(setClauses, fmt.Sprintf("published_year = $%d", argIndex))
		args = append(args, *published_year)
		argIndex++
	}
	if genre != nil {
		setClauses = append(setClauses, fmt.Sprintf("genre = $%d", argIndex))
		args = append(args, *genre)
		argIndex++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE library
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), argIndex)
	args = append(args, id)

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	fmt.Println("Book updated successfully!\n\n\n")
	return nil
}

func Delete_book(id int64) error {
	query := `
		DELETE FROM library
		WHERE id = $1
	`

	_, err := db.Exec(query, id)
	if err != nil{
		return err
	}

	fmt.Println("Book deleted successfully!\n\n\n")
	return nil
}

func View_books() {
	query := `
		SELECT id, title, author, published_year, genre FROM library
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Books in the library:")
	for rows.Next() {
		var id int64
		var title, author, genre string
		var published_year int
		err := rows.Scan(&id, &title, &author, &published_year, &genre)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println( id, title, author, published_year, genre)
	}
	err = rows.Err()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("\n\n\n")
}

func BookExists(id int64) (bool, error){
	query := `
		SELECT 1 FROM library WHERE id = $1
	`

	var exists bool
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil{
		if err == sql.ErrNoRows{
			return false, nil
		}
		return false, err
	}
	return true, nil
}