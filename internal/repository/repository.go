package repository

import (
	"context"
	"database/sql"
	"time"

	"daily-app-go/db/sqlc"
)

type TaskRepository interface {
	FindById(ctx context.Context, id int32) (db.Task, error)
	FindAll(ctx context.Context) ([]db.Task, error)
	Create(ctx context.Context, title string, status string) (db.Task, error)
	Update(ctx context.Context, id int32, title string, status string) (db.Task, error)
	Delete(ctx context.Context, id int32) error
}

type taskRepository struct {
	queries *db.Queries
}

func NewTaskRepository(queries *db.Queries) TaskRepository {
	return &taskRepository{queries: queries}
}

func (r *taskRepository) FindById(ctx context.Context, id int32) (db.Task, error) {
	return r.queries.FindTaskById(ctx, id)
}

func (r *taskRepository) FindAll(ctx context.Context) ([]db.Task, error) {
	return r.queries.FindAllTasks(ctx)
}

func (r *taskRepository) Create(ctx context.Context, title string, status string) (db.Task, error) {
	return r.queries.CreateTask(ctx, db.CreateTaskParams{
		Title:  title,
		Status: status,
	})
}

func (r *taskRepository) Update(ctx context.Context, id int32, title string, status string) (db.Task, error) {
	return r.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID:     id,
		Title:  title,
		Status: status,
	})
}

func (r *taskRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteTask(ctx, id)
}

type SubTaskRepository interface {
	FindById(ctx context.Context, id int32) (db.SubTask, error)
	FindAll(ctx context.Context) ([]db.SubTask, error)
	FindByTaskId(ctx context.Context, taskId int32) ([]db.SubTask, error)
	Create(ctx context.Context, taskId int32, title string, isCompleted bool) (db.SubTask, error)
	Update(ctx context.Context, id int32, title string, isCompleted bool) (db.SubTask, error)
	Delete(ctx context.Context, id int32) error
}

type subTaskRepository struct {
	queries *db.Queries
}

func NewSubTaskRepository(queries *db.Queries) SubTaskRepository {
	return &subTaskRepository{queries: queries}
}

func (r *subTaskRepository) FindById(ctx context.Context, id int32) (db.SubTask, error) {
	return r.queries.FindSubTaskById(ctx, id)
}

func (r *subTaskRepository) FindAll(ctx context.Context) ([]db.SubTask, error) {
	return r.queries.FindAllSubTasks(ctx)
}

func (r *subTaskRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.SubTask, error) {
	return r.queries.FindSubTasksByTaskId(ctx, sql.NullInt32{Int32: taskId, Valid: true})
}

func (r *subTaskRepository) Create(ctx context.Context, taskId int32, title string, isCompleted bool) (db.SubTask, error) {
	return r.queries.CreateSubTask(ctx, db.CreateSubTaskParams{
		TaskID:      sql.NullInt32{Int32: taskId, Valid: true},
		Title:       title,
		IsCompleted: isCompleted,
	})
}

func (r *subTaskRepository) Update(ctx context.Context, id int32, title string, isCompleted bool) (db.SubTask, error) {
	return r.queries.UpdateSubTask(ctx, db.UpdateSubTaskParams{
		ID:          id,
		Title:       title,
		IsCompleted: isCompleted,
	})
}

func (r *subTaskRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteSubTask(ctx, id)
}

type NoteRepository interface {
	FindById(ctx context.Context, id int32) (db.TaskNote, error)
	FindAll(ctx context.Context) ([]db.TaskNote, error)
	FindByTaskId(ctx context.Context, taskId int32) ([]db.TaskNote, error)
	Create(ctx context.Context, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error)
	Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error)
	Delete(ctx context.Context, id int32) error
}

type noteRepository struct {
	queries *db.Queries
}

func NewNoteRepository(queries *db.Queries) NoteRepository {
	return &noteRepository{queries: queries}
}

func (r *noteRepository) FindById(ctx context.Context, id int32) (db.TaskNote, error) {
	return r.queries.FindNoteById(ctx, id)
}

func (r *noteRepository) FindAll(ctx context.Context) ([]db.TaskNote, error) {
	return r.queries.FindAllNotes(ctx)
}

func (r *noteRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.TaskNote, error) {
	return r.queries.FindNotesByTaskId(ctx, sql.NullInt32{Int32: taskId, Valid: true})
}

func (r *noteRepository) Create(ctx context.Context, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error) {
	var stId sql.NullInt32
	if subTaskId != nil {
		stId = sql.NullInt32{Int32: *subTaskId, Valid: true}
	}
	return r.queries.CreateNote(ctx, db.CreateNoteParams{
		TaskID:    sql.NullInt32{Int32: taskId, Valid: true},
		SubTaskID: stId,
		Content:   content,
		Category:  category,
	})
}

func (r *noteRepository) Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error) {
	var stId sql.NullInt32
	if subTaskId != nil {
		stId = sql.NullInt32{Int32: *subTaskId, Valid: true}
	}
	return r.queries.UpdateNote(ctx, db.UpdateNoteParams{
		ID:        id,
		TaskID:    sql.NullInt32{Int32: taskId, Valid: true},
		SubTaskID: stId,
		Content:   content,
		Category:  category,
	})
}

func (r *noteRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteNote(ctx, id)
}

type BlockerRepository interface {
	FindById(ctx context.Context, id int32) (db.Blocker, error)
	FindAll(ctx context.Context) ([]db.Blocker, error)
	FindByTaskId(ctx context.Context, taskId int32) ([]db.Blocker, error)
	FindByTaskIdAndSubTaskId(ctx context.Context, taskId int32, subTaskId *int32) ([]db.Blocker, error)
	Create(ctx context.Context, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error)
	Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error)
	Delete(ctx context.Context, id int32) error
}

type blockerRepository struct {
	queries *db.Queries
}

func NewBlockerRepository(queries *db.Queries) BlockerRepository {
	return &blockerRepository{queries: queries}
}

func (r *blockerRepository) FindById(ctx context.Context, id int32) (db.Blocker, error) {
	return r.queries.FindBlockerById(ctx, id)
}

func (r *blockerRepository) FindAll(ctx context.Context) ([]db.Blocker, error) {
	return r.queries.FindAllBlockers(ctx)
}

func (r *blockerRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.Blocker, error) {
	return r.queries.FindBlockersByTaskId(ctx, sql.NullInt32{Int32: taskId, Valid: true})
}

func (r *blockerRepository) FindByTaskIdAndSubTaskId(ctx context.Context, taskId int32, subTaskId *int32) ([]db.Blocker, error) {
	var stId sql.NullInt32
	if subTaskId != nil {
		stId = sql.NullInt32{Int32: *subTaskId, Valid: true}
	}
	return r.queries.FindBlockersByTaskIdAndSubTaskId(ctx, db.FindBlockersByTaskIdAndSubTaskIdParams{
		TaskID:    sql.NullInt32{Int32: taskId, Valid: true},
		SubTaskID: stId,
	})
}

func (r *blockerRepository) Create(ctx context.Context, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error) {
	var stId sql.NullInt32
	if subTaskId != nil {
		stId = sql.NullInt32{Int32: *subTaskId, Valid: true}
	}
	return r.queries.CreateBlocker(ctx, db.CreateBlockerParams{
		TaskID:     sql.NullInt32{Int32: taskId, Valid: true},
		SubTaskID:  stId,
		Reason:     reason,
		IsResolved: sql.NullBool{Bool: isResolved, Valid: true},
	})
}

func (r *blockerRepository) Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error) {
	var stId sql.NullInt32
	if subTaskId != nil {
		stId = sql.NullInt32{Int32: *subTaskId, Valid: true}
	}
	return r.queries.UpdateBlocker(ctx, db.UpdateBlockerParams{
		ID:         id,
		TaskID:     sql.NullInt32{Int32: taskId, Valid: true},
		SubTaskID:  stId,
		Reason:     reason,
		IsResolved: sql.NullBool{Bool: isResolved, Valid: true},
	})
}

func (r *blockerRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteBlocker(ctx, id)
}

type DailySessionRepository interface {
	FindById(ctx context.Context, id int32) (db.DailySession, error)
	FindAll(ctx context.Context) ([]db.DailySession, error)
	Create(ctx context.Context, sessionDate time.Time, rawNotesBlob string, generatedScript string) (db.DailySession, error)
	Delete(ctx context.Context, id int32) error
}

type dailySessionRepository struct {
	queries *db.Queries
}

func NewDailySessionRepository(queries *db.Queries) DailySessionRepository {
	return &dailySessionRepository{queries: queries}
}

func (r *dailySessionRepository) FindById(ctx context.Context, id int32) (db.DailySession, error) {
	return r.queries.FindDailySessionById(ctx, id)
}

func (r *dailySessionRepository) FindAll(ctx context.Context) ([]db.DailySession, error) {
	return r.queries.FindAllDailySessions(ctx)
}

func (r *dailySessionRepository) Create(ctx context.Context, sessionDate time.Time, rawNotesBlob string, generatedScript string) (db.DailySession, error) {
	return r.queries.CreateDailySession(ctx, db.CreateDailySessionParams{
		SessionDate:     sessionDate,
		RawNotesBlob:    rawNotesBlob,
		GeneratedScript: generatedScript,
	})
}

func (r *dailySessionRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteDailySession(ctx, id)
}
