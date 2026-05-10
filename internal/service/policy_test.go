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

func TestTaskPolicyService_EnsureCanTransitionStatus_Blocked(t *testing.T) {
	mockBlockerRepo := new(MockBlockerRepository)
	mockSubTaskRepo := new(MockSubTaskRepository)
	service := NewTaskPolicyService(mockBlockerRepo, mockSubTaskRepo)

	task := db.Task{ID: 1, Status: string(models.TaskStatusBlocked)}
	blocker := db.Blocker{ID: 1, IsResolved: sql.NullBool{Bool: false, Valid: true}}

	mockBlockerRepo.On("FindByTaskId", mock.Anything, int32(1)).Return([]db.Blocker{blocker}, nil)

	err := service.EnsureCanTransitionStatus(context.Background(), task, models.TaskStatusInProgress)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot change task status from BLOCKED while unresolved blockers exist")
	mockBlockerRepo.AssertExpectations(t)
}

func TestTaskPolicyService_EnsureCanTransitionStatus_DoneWithIncompleteSubTasks(t *testing.T) {
	mockBlockerRepo := new(MockBlockerRepository)
	mockSubTaskRepo := new(MockSubTaskRepository)
	service := NewTaskPolicyService(mockBlockerRepo, mockSubTaskRepo)

	task := db.Task{ID: 1, Status: string(models.TaskStatusInProgress)}
	subTask := db.SubTask{ID: 1, IsCompleted: false}

	mockSubTaskRepo.On("FindByTaskId", mock.Anything, int32(1)).Return([]db.SubTask{subTask}, nil)

	err := service.EnsureCanTransitionStatus(context.Background(), task, models.TaskStatusDone)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot change task status to DONE while incomplete subtasks exist")
	mockSubTaskRepo.AssertExpectations(t)
}
