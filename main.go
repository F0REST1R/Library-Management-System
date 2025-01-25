package main

import (
	"BD_LIBRARY/db"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func StartMenu() int64 {
	fmt.Println("Welcome to the Library Management System!\n1. Add a book\n2. Update a book\n3. Delete a book\n4. View all books\n5. Exit")
	fmt.Print("Enter your choice: ")
	var choose int
	fmt.Scan(&choose)
	return int64(choose)
}

func AddData() (string, string, int64, string) {
	fmt.Println("!!ATTENTION!!")
	fmt.Println("Write all the data in one line, separating them with ' / ', instead of a space, use '_' . If there is no data, just enter 0")
	fmt.Println("For example:")
	fmt.Println("Title_title/author/published year/genre")
	var data string
	fmt.Scan(&data)
	fmt.Println(data)

	parts := strings.Split(data, "/")
	if len(parts) != 4 {
		log.Fatal("Invalid input format. Please follow the example.")
	}

	title := strings.TrimSpace(parts[0])
	author := strings.TrimSpace(parts[1])
	published_year := strings.TrimSpace(parts[2])
	genre := strings.TrimSpace(parts[3])

	var publishedYearInt int64
	_, err := fmt.Sscanf(published_year, "%d", &publishedYearInt)
	if err != nil {
		fmt.Println("Invalid published year. Please enter a valid number.")
		return "", "", 0, ""
	}

	return title, author, publishedYearInt, genre
}


func main() {
	db.InitDB()
	db.Create_table()
	for {
		choose := StartMenu()
		switch choose {
		case 1:
			title, author, published_year, genre := AddData()
			id, err := db.Add_book(title, author, published_year, genre)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The book was added successfully %d\n\n\n\n", id)
		case 2:
			fmt.Println("Select the ID of the book you want to edit")
			var idstr string
			fmt.Scan(&idstr)

			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil{
				fmt.Println("Please enter a valid number\n")
				continue
			}

			exists, err := db.BookExists(id)
			if err != nil {
				log.Fatal(err)
			}

			if !exists {
				fmt.Println("Book with the given ID does not exist\n")
				continue
			}

			fmt.Println("Enter the number of what you want to do.\n 1. Change title.\n 2. Change author.\n 3. Change published year.\n 4. Change genre.\n 5. Change all")
			var choose_do int
			fmt.Scan(&choose_do)
			switch choose_do {
			case 1:
				fmt.Println("Instead of a space, use '_'\nEnter the title:")
				var title string
				fmt.Scan(&title)
				db.Update_book(id, &title, nil, nil, nil)
			case 2:
				fmt.Println("Instead of a space, use '_'\nEnter the author:")
				var author string
				fmt.Scan(&author)
				db.Update_book(id, nil, &author, nil, nil)
			case 3:
				fmt.Println("Enter the published_year:")
				var published_year_str string
				fmt.Scan(&published_year_str)
				published_year, err := strconv.ParseInt(published_year_str, 10, 64)
				if err != nil {
					fmt.Println("Please enter a valid number\n")
					continue
				}
				db.Update_book(id, nil, nil, &published_year, nil)
			case 4:
				fmt.Println("Instead of a space, use '_'\nEnter the genre:")
				var genre string
				fmt.Scan(&genre)
				db.Update_book(id, nil, nil, nil, &genre)
			case 5:
				title, author, published_year, genre := AddData()
				db.Update_book(id, &title, &author, &published_year, &genre)
			default:
				fmt.Println("Please enter a number from 1 to 5")
			}
		case 3:
			fmt.Println("Enter the ID of the book you want to delete")
			var idstr string
			fmt.Scan(&idstr)

			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil{
				fmt.Println("Please enter a valid number\n")
				continue
			}

			exists, err := db.BookExists(id)
			if err != nil {
				log.Fatal(err)
			}

			if !exists {
				fmt.Println("Book with the given ID does not exist\n")
				continue
			}

			err = db.Delete_book(id)
			if err != nil {
				log.Fatal(err)
			}	
		case 4:
			db.View_books()
		case 5:
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("\n\n\nPlease enter a number from 1 to 5")
		}
	}
}
