package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskSynchronizerService_SyncTaskWithNewSubTasks(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	service := NewTaskSynchronizerService(mockTaskRepo, nil)

	task := db.Task{ID: 1, Title: "Task", Status: string(models.TaskStatusDone)}

	mockTaskRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockTaskRepo.On("Update", mock.Anything, int32(1), "Task", string(models.TaskStatusInProgress)).Return(db.Task{}, nil)

	err := service.SyncTaskWithNewSubTasks(context.Background(), 1)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
}

func TestSubTaskSynchronizerService_SyncSubTaskWithBlockers_Reopen(t *testing.T) {
	mockSubTaskRepo := new(MockSubTaskRepository)
	mockBlockerRepo := new(MockBlockerRepository)
	service := NewSubTaskSynchronizerService(mockSubTaskRepo, mockBlockerRepo)

	subTaskID := int32(2)
	subTask := db.SubTask{ID: subTaskID, TaskID: sql.NullInt32{Int32: 1, Valid: true}, Title: "Sub 1", IsCompleted: true}
	blocker := db.Blocker{ID: 3, IsResolved: sql.NullBool{Bool: false, Valid: true}}

	mockSubTaskRepo.On("FindById", mock.Anything, subTaskID).Return(subTask, nil)
	mockBlockerRepo.On("FindByTaskIdAndSubTaskId", mock.Anything, int32(1), &subTaskID).Return([]db.Blocker{blocker}, nil)
	mockSubTaskRepo.On("Update", mock.Anything, subTaskID, "Sub 1", false).Return(db.SubTask{}, nil)

	err := service.SyncSubTaskWithBlockers(context.Background(), 1, &subTaskID)

	assert.NoError(t, err)
	mockSubTaskRepo.AssertExpectations(t)
	mockBlockerRepo.AssertExpectations(t)
}
