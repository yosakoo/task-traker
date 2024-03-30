package models

import(
	"time"
)


type Task struct {
    ID     int
    Status string
    UserID int
    Title  string
    Text   *string
    Time   *time.Time
}