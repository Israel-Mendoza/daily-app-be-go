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

func TestTaskSynchronizerService_SyncTaskWithBlockers_Blocked(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockBlockerRepo := new(MockBlockerRepository)
	service := NewTaskSynchronizerService(mockTaskRepo, mockBlockerRepo)

	task := db.Task{ID: 1, Title: "Task", Status: string(models.TaskStatusInProgress)}
	blocker := db.Blocker{ID: 1, IsResolved: sql.NullBool{Bool: false, Valid: true}}

	mockTaskRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockBlockerRepo.On("FindByTaskId", mock.Anything, int32(1)).Return([]db.Blocker{blocker}, nil)
	mockTaskRepo.On("Update", mock.Anything, int32(1), "Task", string(models.TaskStatusBlocked)).Return(db.Task{}, nil)

	err := service.SyncTaskWithBlockers(context.Background(), 1)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
	mockBlockerRepo.AssertExpectations(t)
}

func TestTaskSynchronizerService_SyncTaskWithBlockers_InProgress(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockBlockerRepo := new(MockBlockerRepository)
	service := NewTaskSynchronizerService(mockTaskRepo, mockBlockerRepo)

	task := db.Task{ID: 1, Title: "Task", Status: string(models.TaskStatusBlocked)}
	blocker := db.Blocker{ID: 1, IsResolved: sql.NullBool{Bool: true, Valid: true}}

	mockTaskRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockBlockerRepo.On("FindByTaskId", mock.Anything, int32(1)).Return([]db.Blocker{blocker}, nil)
	mockTaskRepo.On("Update", mock.Anything, int32(1), "Task", string(models.TaskStatusInProgress)).Return(db.Task{}, nil)

	err := service.SyncTaskWithBlockers(context.Background(), 1)

	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
	mockBlockerRepo.AssertExpectations(t)
}
