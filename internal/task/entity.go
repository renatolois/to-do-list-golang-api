package task

import (
	"time"
	"database/sql"
)

type Task struct {
	Id						string					`json:"id"`
	Title					string					`json:"title"`
	Completed			bool						`json:"completed"`
	CompletedAt		sql.NullTime		`json:"completed_at"`
	CreatedAt			time.Time				`json:"created_at"`
}
