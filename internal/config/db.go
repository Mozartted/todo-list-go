package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type TaskStatus int

const (
	PENDING TaskStatus = iota
	DONE
)

func (t TaskStatus) String() string {
	return [...]string{"PENDING", "DONE"}[t]
}

func (t TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TaskStatus) UnmarshalJSON(data []byte) error {
	// fmt.Printf("Called UnmarshalJSON: %v", data)
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "PENDING":
		*t = PENDING
	case "DONE":
		*t = DONE
	default:
		return fmt.Errorf("invalid status type:  %v", s)
	}
	return nil
}

type TaskData struct {
	Name string `json:"id"`
	// Description string     `json:"description"`
	Status TaskStatus `json:"status"`
}

func (t TaskData) toJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DbConnect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func InitialDBMigration(db *sql.DB) error {
	if ping := db.Ping(); ping != nil {
		log.Fatalf("db connection is lost %v", ping)
	}

	resp, err := db.Prepare("create table todoList( todoSet jsonb, id integer primary key autoincrement )")
	if err != nil {
		if !strings.Contains(err.Error(), "table todoList already exists") {
			return err
		}
		return nil

	}

	if _, err := resp.Exec(); err != nil {
		if !strings.Contains(err.Error(), "table todoList already exists") {
			return err
		}
	}

	return nil
}

func SaveTodoData(db *sql.DB, task TaskData) error {
	preprocessedTask, e := task.toJSON()
	if e != nil {
		log.Fatalf("%s", e.Error())
	}

	fmt.Printf("%+v, processed Task: %s", db, preprocessedTask)

	resp, err := db.Prepare(("insert into todoList(todoSet) values (?)"))
	if err != nil {
		fmt.Printf("Other Error print %v", err.Error())
		log.Fatalf("%s", e.Error())
	}

	if _, err := resp.Exec(preprocessedTask); err != nil {
		return err
	}

	return nil
}

func RetrieveAll(dbC *sql.DB) []TaskData {
	row, err := dbC.Query("select todoSet from todoList")
	if err != nil {
		fmt.Printf("Other Error print %v", err.Error())
		log.Fatalf("%s", err.Error())
	}

	var todoTaskList []TaskData

	for row.Next() {
		var todoSet string
		if err := row.Scan(&todoSet); err != nil {
			log.Fatalf("%s", err.Error())
		}

		var currentTaskData TaskData
		if err := json.Unmarshal([]byte(todoSet), &currentTaskData); err != nil {
			log.Fatalf("%s", err.Error())
		}
		todoTaskList = append(todoTaskList, currentTaskData)

	}

	return todoTaskList
}
