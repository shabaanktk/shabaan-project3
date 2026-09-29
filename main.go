package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

// Struct with three fields of different types
type Student struct {
	name   string
	age    int
	active bool
}

func main() {
	var choice int

	for {
		fmt.Println("\n--- Project 1 Menu ---")
		fmt.Println("1. Read a file")
		fmt.Println("2. Count bytes and characters")
		fmt.Println("3. Work with a struct")
		fmt.Println("4. Work with a list of structs")
		fmt.Println("5. Quit")
		fmt.Print("Enter your choice: ")

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			readFile()
		case 2:
			countBytesAndCharacters()
		case 3:
			workWithStruct()
		case 4:
			workWithStructList()
		case 5:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func readFile() {
	var fileName string
	var x int

	fmt.Print("Enter the name of the file: ")
	fmt.Scanln(&fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Could not open the file.")
		return
	}
	defer file.Close()

	fmt.Print("Enter a number X: ")
	fmt.Scanln(&x)

	if x <= 0 {
		fmt.Println("X must be greater than 0.")
		return
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		if lineNumber%x == 0 {
			fmt.Println(scanner.Text())
		}
		lineNumber++
	}
}

func countBytesAndCharacters() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a string: ")
	input, _ := reader.ReadString('\n')

	// Remove Enter from the end of the input
	if len(input) > 0 && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}

	if len(input) > 0 && input[len(input)-1] == '\r' {
		input = input[:len(input)-1]
	}

	fmt.Println("Number of bytes:", len(input))
	fmt.Println("Number of characters:", utf8.RuneCountInString(input))
}

func workWithStruct() {
	student := Student{
		name:   "Shabaan",
		age:    21,
		active: true,
	}

	fmt.Println("\nOriginal struct:")
	fmt.Println(student)

	changeStudent(student)

	fmt.Println("Struct after returning to original function:")
	fmt.Println(student)
}

func changeStudent(student Student) {
	student.age = 22

	fmt.Println("Struct inside changeStudent after changing age:")
	fmt.Println(student)
}

func workWithStructList() {
	students := []Student{
		{
			name:   "Shabaan",
			age:    21,
			active: true,
		},
		{
			name:   "John",
			age:    20,
			active: true,
		},
	}

	fmt.Println("\nOriginal list of structs:")
	fmt.Println(students)

	changeStudentList(students)

	fmt.Println("List after returning to original function:")
	fmt.Println(students)
}

func changeStudentList(students []Student) {
	students[1].age = 25

	fmt.Println("List inside changeStudentList after changing second student:")
	fmt.Println(students)
}
