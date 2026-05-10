package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_Save(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo, nil, nil, nil, nil)

	req := models.CreateTaskRequest{Title: "New Task"}
	task := db.Task{ID: 1, Title: "New Task", Status: string(models.TaskStatusTodo)}

	mockRepo.On("Create", mock.Anything, "New Task", string(models.TaskStatusTodo)).Return(task, nil)

	res, err := service.Save(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "New Task", res.Title)
	assert.Equal(t, string(models.TaskStatusTodo), res.Status)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_StartTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockPolicy := new(MockTaskPolicyService)
	service := NewTaskService(mockRepo, nil, nil, nil, mockPolicy)

	task := db.Task{ID: 1, Title: "Task", Status: string(models.TaskStatusTodo)}
	mockRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockPolicy.On("EnsureCanTransitionStatus", mock.Anything, task, models.TaskStatusInProgress).Return(nil)
	mockRepo.On("Update", mock.Anything, int32(1), "Task", string(models.TaskStatusInProgress)).Return(db.Task{}, nil)

	err := service.StartTask(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPolicy.AssertExpectations(t)
}

func TestTaskService_CompleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockPolicy := new(MockTaskPolicyService)
	service := NewTaskService(mockRepo, nil, nil, nil, mockPolicy)

	task := db.Task{ID: 1, Title: "Task", Status: string(models.TaskStatusInProgress)}
	mockRepo.On("FindById", mock.Anything, int32(1)).Return(task, nil)
	mockPolicy.On("EnsureCanTransitionStatus", mock.Anything, task, models.TaskStatusDone).Return(nil)
	mockRepo.On("Update", mock.Anything, int32(1), "Task", string(models.TaskStatusDone)).Return(db.Task{}, nil)

	err := service.CompleteTask(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPolicy.AssertExpectations(t)
}
