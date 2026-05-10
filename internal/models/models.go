package models

import (
	"database/sql"
	"time"
)

// TaskStatus defines the status of a task
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusBlocked    TaskStatus = "BLOCKED"
	TaskStatusDone       TaskStatus = "DONE"
	TaskStatusCanceled   TaskStatus = "CANCELED"
)

// NoteCategory defines the category of a note
type NoteCategory string

const (
	NoteCategoryGeneral  NoteCategory = "GENERAL"
	NoteCategoryLearning NoteCategory = "LEARNING"
	NoteCategoryLinks    NoteCategory = "LINKS"
	NoteCategoryIdeas    NoteCategory = "IDEAS"
	NoteCategoryDecision NoteCategory = "DECISION"
	NoteCategoryIssues   NoteCategory = "ISSUES"
)

// Request models
type CreateTaskRequest struct {
	Title string `json:"title"`
}

type CreateAndUpdateSubTaskRequest struct {
	Title string `json:"title"`
}

type CreateBlockerRequest struct {
	SubTaskID *int32 `json:"sub_task_id"`
	Reason    string `json:"reason"`
}

type UpdateBlockerRequest struct {
	SubTaskID  *int32  `json:"sub_task_id"`
	Reason     *string `json:"reason"`
	IsResolved *bool   `json:"is_resolved"`
}

type CreateNoteRequest struct {
	SubTaskID *int32 `json:"sub_task_id"`
	Content   string `json:"content"`
	Category  string `json:"category"`
}

type UpdateNoteRequest struct {
	SubTaskID *int32  `json:"sub_task_id"`
	Content   *string `json:"content"`
	Category  *string `json:"category"`
}

// Response models
type TaskResponse struct {
	ID        int32  `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SubTaskResponse struct {
	ID          int32  `json:"id"`
	TaskID      int32  `json:"task_id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BlockerResponse struct {
	ID         int32  `json:"id"`
	TaskID     int32  `json:"task_id"`
	SubTaskID  *int32 `json:"sub_task_id"`
	Reason     string `json:"reason"`
	IsResolved bool   `json:"is_resolved"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type NoteResponse struct {
	ID        int32  `json:"id"`
	TaskID    int32  `json:"task_id"`
	SubTaskID *int32 `json:"sub_task_id"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TaskDetailsResponse struct {
	ID        int32             `json:"id"`
	Title     string            `json:"title"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Notes     []NoteResponse    `json:"notes"`
	SubTasks  []SubTaskResponse `json:"sub_tasks"`
	Blockers  []BlockerResponse `json:"blockers"`
}

// Task represents the tasks table
type Task struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SubTask represents the sub_tasks table
type SubTask struct {
	ID          int32     `json:"id"`
	TaskID      int32     `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Blocker represents the blockers table
type Blocker struct {
	ID         int32         `json:"id"`
	TaskID     int32         `json:"task_id"`
	SubTaskID  sql.NullInt32 `json:"sub_task_id"`
	Reason     string        `json:"reason"`
	IsResolved bool          `json:"is_resolved"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// Note represents the task_notes table
type Note struct {
	ID        int32         `json:"id"`
	TaskID    int32         `json:"task_id"`
	SubTaskID sql.NullInt32 `json:"sub_task_id"`
	Content   string        `json:"content"`
	Category  string        `json:"category"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// DailySession represents the daily_sessions table
type DailySession struct {
	ID              int32     `json:"id"`
	SessionDate     time.Time `json:"session_date"`
	RawNotesBlob    string    `json:"raw_notes_blob"`
	CreatedAt       time.Time `json:"created_at"`
	GeneratedScript string    `json:"generated_script"`
}
