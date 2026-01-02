package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mozartted/simple_todo_server/internal/model"
)

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

func SaveTodoData(db *sql.DB, task model.TaskData) error {
	preprocessedTask, e := task.ToJSON()
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

func RetrieveAll(dbC *sql.DB) []model.TaskData {
	row, err := dbC.Query("select todoSet, id  from todoList")
	if err != nil {
		fmt.Printf("Other Error print %v", err.Error())
		log.Fatalf("%s", err.Error())
	}

	var todoTaskList []model.TaskData

	for row.Next() {
		var todoSet string
		var todoId uint
		if err := row.Scan(&todoSet, &todoId); err != nil {
			log.Fatalf("%s", err.Error())
		}

		var currentTaskData model.TaskData
		if err := json.Unmarshal([]byte(todoSet), &currentTaskData); err != nil {
			log.Fatalf("%s", err.Error())
		}
		currentTaskData.Id = todoId
		todoTaskList = append(todoTaskList, currentTaskData)

	}

	return todoTaskList
}

func DeleteTask(dbc *sql.DB, index uint) []model.TaskData {
	row, err := dbc.Prepare("delete from todoList where id=?")
	if err != nil {
		log.Fatalf("Something went wrong: %v", err.Error())
	}

	if _, err := row.Exec(index); err != nil {
		log.Fatalf("Something went wrong: %v", err.Error())
	}

	return RetrieveAll(dbc)

}

func UpdateStatus(dbc *sql.DB, index uint) model.TaskData {
	row, err := dbc.Query(fmt.Sprintf("select todoSet, id from todoList where id=%d", index))
	if err != nil {
		log.Fatalf("Something went wrong: %v", err.Error())
	}
	var currentTaskData model.TaskData

	for row.Next() {
		var id uint
		var todoSet string
		if err := row.Scan(&todoSet, &id); err != nil {
			log.Fatalf("%s", err.Error())
		}

		if err := json.Unmarshal([]byte(todoSet), &currentTaskData); err != nil {
			log.Fatalf("%s", err.Error())
		}

		switch currentTaskData.Status {
		case model.PENDING:
			currentTaskData.Status = model.DONE
		case model.DONE:
			currentTaskData.Status = model.PENDING
		default:
			continue
		}

		SaveTodoData(dbc, currentTaskData)

	}
	return currentTaskData
}
