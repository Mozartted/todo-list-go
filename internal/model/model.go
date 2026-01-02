package model

import (
	"encoding/json"
	"fmt"
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
	Id   uint   `json:"key"`
	Name string `json:"id"`
	// Description string     `json:"description"`
	Status TaskStatus `json:"status"`
}

func (t TaskData) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
