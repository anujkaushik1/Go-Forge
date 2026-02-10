package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TaskHeaders struct {
	id         string
	user_id    string
	task_name  string
	category   string
	created_at string
	status     string
	expires_on string
}

type Users struct {
	user_id  string
	email    string
	password string
}

var fileMap = map[string]any{
	"task_manager.txt": &TaskHeaders{},
	"users.txt":        &Users{},
}

func addHeadersToFile(file *os.File, headers any) string {
	t := reflect.TypeOf(headers)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldNames := make([]string, t.NumField())
	for i := 0; i < len(fieldNames); i++ {
		fieldNames[i] = t.Field(i).Name
	}

	_, err := file.WriteString(strings.Join(fieldNames, ",") + "\n")

	if err != nil {
		fmt.Println(err)
		return "failed"
	}

	return "success"

}

func (u Users) toCsv() []string {
	return []string{u.user_id, u.email, u.password}
}

func addDataToFile(file *os.File, data *[]string) string {
	_, err := file.WriteString(strings.Join(*data, ",") + "\n")

	if err != nil {
		fmt.Println(err)
		return "failed"
	}

	return "success"
}

func initializeFile() {

	for fileName, structPtr := range fileMap {
		_, err := os.Stat(fileName)
		fileExists := !os.IsNotExist(err)

		if fileExists {
			fmt.Println("File already exists ->> " + fileName)
			continue
		}

		file, err := os.OpenFile(
			fileName,
			os.O_CREATE|os.O_RDWR|os.O_APPEND,
			0644,
		)
		println(file.Fd())

		if err != nil {
			panic(err)
		}
		defer file.Close()
		t := reflect.TypeOf(structPtr)
		if t.Kind() == reflect.Ptr {
			addHeadersToFile(file, structPtr)
		} else {
			addHeadersToFile(file, &structPtr)

		}

	}

}

func readFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(
		fileName,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		0644,
	)
	return file, err
}

func signupUser(file *os.File, email string, password string) string {

	user_id := strconv.Itoa(rand.Intn(10000) + 1)

	user := Users{user_id: user_id, email: email, password: password}

	csvData := user.toCsv()
	status := addDataToFile(file, &csvData)

	if status == "success" {
		writeSession(&user)
		return "SIGNUP_SUCEESS"
	}

	return "SIGNUP_FAILED"
}

func readSession() (*Users, error) {
	f, err := os.Open("session.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("invalid session")
	}

	row := records[1]

	return &Users{
		user_id:  row[0],
		email:    row[1],
		password: row[2],
	}, nil
}

func writeSession(user *Users) error {
	f, err := os.Create("session.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"user_id", "email", "password"})
	writer.Write([]string{
		user.user_id,
		user.email,
		user.password,
	})

	return writer.Error()
}

func loginUser(email string, password string) string {

	file, err := readFile("users.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	for i := 1; i < len(records); i++ {
		record := records[i]
		recordUserId := record[0]
		recordEmail := record[1]
		recordPassword := record[2]

		if recordEmail == email && recordPassword == password {
			writeSession(&Users{user_id: recordUserId, email: recordEmail, password: recordPassword})
			return "LOGIN_SUCCESSFULL"
		}
	}

	signupResponse := signupUser(file, email, password)

	return signupResponse

}
func (task *TaskHeaders) addTask() {
	file, err := os.OpenFile("task_manager.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	task.id = strconv.Itoa(rand.Intn(9000) + 1000)
	err = writer.Write([]string{
		task.id,
		task.user_id,
		task.task_name,
		task.category,
		task.created_at,
		task.status,
		task.expires_on,
	})

	if err != nil {
		panic(err)
	}

	writer.Flush()
}

func printTasks(tasks *[][]string) {
	for _, value := range *tasks {
		fmt.Println("Task Name : ", value[2])
		fmt.Println("Category : ", value[3])

		fmt.Println("---------")
	}
}

func listTasksByUserId(user_id string) {

	file, err := readFile("task_manager.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	recordsOfUser := make([][]string, 0, len(records))
	for i := 1; i < len(records); i++ {
		recordUserId := records[i][1]

		if recordUserId == user_id {
			recordsOfUser = append(recordsOfUser, records[i])
		}
	}

	printTasks(&recordsOfUser)

}

func main() {
	initializeFile()

	command := os.Args[1]

	if command == "login" {
		loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
		email := loginCmd.String("email", "", "user email")
		password := loginCmd.String("password", "", "user pass")
		loginCmd.Parse(os.Args[2:])
		response := loginUser(*email, *password)

		fmt.Println(response)

	}

	if command == "add" {
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		taskName := addCmd.String("task_name", "", "taskname")
		category := addCmd.String("category", "", "category")

		addCmd.Parse(os.Args[2:])
		sessionData, _ := readSession()

		task := TaskHeaders{
			user_id:   sessionData.user_id,
			task_name: *taskName,
			category:  *category,
			status:    "pending",
		}

		task.addTask()

	}

	if command == "list" {
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Parse(os.Args[2:])
		sessionData, _ := readSession()
		listTasksByUserId(sessionData.user_id)
	}

}
