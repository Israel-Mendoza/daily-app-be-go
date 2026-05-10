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

func TestBlockerService_Save(t *testing.T) {
	mockRepo := new(MockBlockerRepository)
	mockGuard := new(MockTaskOwnershipGuardService)
	mockTaskSync := new(MockTaskSynchronizerService)
	mockSubTaskSync := new(MockSubTaskSynchronizerService)
	service := NewBlockerService(mockRepo, mockGuard, mockTaskSync, mockSubTaskSync)

	task := db.Task{ID: 1, Title: "Task"}
	subTaskID := int32(2)
	subTask := db.SubTask{ID: subTaskID, Title: "SubTask"}
	req := models.CreateBlockerRequest{Reason: "Blocker", SubTaskID: &subTaskID}
	blocker := db.Blocker{ID: 3, Reason: "Blocker"}

	mockGuard.On("EnsureTaskExists", mock.Anything, int32(1)).Return(task, nil)
	mockGuard.On("EnsureSubTaskBelongsToTask", mock.Anything, int32(1), &subTaskID).Return(&subTask, nil)
	mockRepo.On("Create", mock.Anything, int32(1), &subTaskID, "Blocker", false).Return(blocker, nil)
	mockTaskSync.On("SyncTaskWithBlockers", mock.Anything, int32(1)).Return(nil)
	mockSubTaskSync.On("SyncSubTaskWithBlockers", mock.Anything, int32(1), &subTaskID).Return(nil)

	res, err := service.Save(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Blocker", res.Reason)
	mockRepo.AssertExpectations(t)
	mockGuard.AssertExpectations(t)
	mockTaskSync.AssertExpectations(t)
	mockSubTaskSync.AssertExpectations(t)
}

func TestBlockerService_Resolve(t *testing.T) {
	mockRepo := new(MockBlockerRepository)
	mockTaskSync := new(MockTaskSynchronizerService)
	mockSubTaskSync := new(MockSubTaskSynchronizerService)
	service := NewBlockerService(mockRepo, nil, mockTaskSync, mockSubTaskSync)

	subTaskID := int32(2)
	blocker := db.Blocker{
		ID:        3,
		TaskID:    sql.NullInt32{Int32: 1, Valid: true},
		SubTaskID: sql.NullInt32{Int32: subTaskID, Valid: true},
		Reason:    "Blocker",
	}

	mockRepo.On("FindById", mock.Anything, int32(3)).Return(blocker, nil)
	mockRepo.On("Update", mock.Anything, int32(3), int32(1), &subTaskID, "Blocker", true).Return(blocker, nil)
	mockTaskSync.On("SyncTaskWithBlockers", mock.Anything, int32(1)).Return(nil)
	mockSubTaskSync.On("SyncSubTaskWithBlockers", mock.Anything, int32(1), &subTaskID).Return(nil)

	_, err := service.Resolve(context.Background(), 3, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockTaskSync.AssertExpectations(t)
	mockSubTaskSync.AssertExpectations(t)
}
