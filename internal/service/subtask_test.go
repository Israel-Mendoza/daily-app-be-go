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

func TestSubTaskService_Save(t *testing.T) {
	mockRepo := new(MockSubTaskRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockTaskSync := new(MockTaskSynchronizerService)
	service := NewSubTaskService(mockRepo, mockTaskRepo, mockTaskSync, nil)

	task := db.Task{ID: 1, Title: "Task"}
	req := models.CreateAndUpdateSubTaskRequest{Title: "SubTask"}
	subTask := db.SubTask{ID: 2, Title: "SubTask"}

	mockTaskRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockRepo.On("Create", mock.Anything, int32(1), "SubTask", false).Return(subTask, nil)
	mockTaskSync.On("SyncTaskWithNewSubTasks", mock.Anything, int32(1)).Return(nil)

	res, err := service.Save(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "SubTask", res.Title)
	mockRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
	mockTaskSync.AssertExpectations(t)
}

func TestSubTaskService_Complete(t *testing.T) {
	mockRepo := new(MockSubTaskRepository)
	mockPolicy := new(MockSubTaskPolicyService)
	service := NewSubTaskService(mockRepo, nil, nil, mockPolicy)

	subTask := db.SubTask{ID: 2, TaskID: sql.NullInt32{Int32: 1, Valid: true}, Title: "SubTask", IsCompleted: false}

	mockRepo.On("FindById", mock.Anything, int32(2)).Return(subTask, nil)
	mockPolicy.On("EnsureCanCompleteSubTask", mock.Anything, subTask).Return(nil)
	mockRepo.On("Update", mock.Anything, int32(2), "SubTask", true).Return(subTask, nil)

	_, err := service.Complete(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPolicy.AssertExpectations(t)
}
