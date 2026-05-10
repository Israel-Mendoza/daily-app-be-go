package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskWorkflowIntegration(t *testing.T) {
	ctx := context.Background()

	// Initialize Mocks
	mockTaskRepo := new(MockTaskRepository)
	mockSubTaskRepo := new(MockSubTaskRepository)
	mockBlockerRepo := new(MockBlockerRepository)
	mockNoteRepo := new(MockNoteRepository)

	// Initialize Services (using real implementations with mocked repositories)
	taskPolicyService := NewTaskPolicyService(mockBlockerRepo, mockSubTaskRepo)
	subTaskPolicyService := NewSubTaskPolicyService(mockBlockerRepo)
	taskSync := NewTaskSynchronizerService(mockTaskRepo, mockBlockerRepo)
	subTaskSync := NewSubTaskSynchronizerService(mockSubTaskRepo, mockBlockerRepo)
	guard := NewTaskOwnershipGuardService(mockTaskRepo, mockSubTaskRepo)

	blockerService := NewBlockerService(mockBlockerRepo, guard, taskSync, subTaskSync)
	noteService := NewNoteService(mockNoteRepo, guard)
	subTaskService := NewSubTaskService(mockSubTaskRepo, mockTaskRepo, taskSync, subTaskPolicyService)
	taskService := NewTaskService(mockTaskRepo, blockerService, subTaskService, noteService, taskPolicyService)

	t.Run("should follow full task workflow with blockers", func(t *testing.T) {
		taskId := int32(1)
		task := db.Task{ID: taskId, Title: "Task 1", Status: string(models.TaskStatusTodo)}

		// 1. Create Task
		mockTaskRepo.On("Create", ctx, "Task 1", string(models.TaskStatusTodo)).Return(task, nil)
		taskResponse, err := taskService.Save(ctx, models.CreateTaskRequest{Title: "Task 1"})
		assert.NoError(t, err)
		assert.Equal(t, taskId, taskResponse.ID)

		// 2. Start Task
		mockTaskRepo.On("FindById", ctx, taskId).Return(task, nil).Once()
		// Policy check
		mockBlockerRepo.On("FindByTaskId", ctx, taskId).Return([]db.Blocker{}, nil).Once()
		mockSubTaskRepo.On("FindByTaskId", ctx, taskId).Return([]db.SubTask{}, nil).Once()
		// Update to IN_PROGRESS
		taskInProgress := task
		taskInProgress.Status = string(models.TaskStatusInProgress)
		mockTaskRepo.On("Update", ctx, taskId, "Task 1", string(models.TaskStatusInProgress)).Return(taskInProgress, nil).Once()

		err = taskService.StartTask(ctx, taskId)
		assert.NoError(t, err)

		// 3. Create Blocker -> Should move to BLOCKED
		subTaskID := (*int32)(nil)
		blocker := db.Blocker{ID: 10, TaskID: sql.NullInt32{Int32: taskId, Valid: true}, Reason: "Wait for info"}

		mockTaskRepo.On("FindById", ctx, taskId).Return(taskInProgress, nil).Once() // guard
		mockBlockerRepo.On("Create", ctx, taskId, subTaskID, "Wait for info", false).Return(blocker, nil).Once()
		// Sync
		mockTaskRepo.On("FindById", ctx, taskId).Return(taskInProgress, nil).Once() // sync
		mockBlockerRepo.On("FindByTaskId", ctx, taskId).Return([]db.Blocker{{ID: 10, IsResolved: sql.NullBool{Bool: false, Valid: true}}}, nil).Once()
		mockTaskRepo.On("Update", ctx, taskId, "Task 1", string(models.TaskStatusBlocked)).Return(db.Task{}, nil).Once()

		_, err = blockerService.Save(ctx, taskId, models.CreateBlockerRequest{Reason: "Wait for info", SubTaskID: nil})
		assert.NoError(t, err)

		// 4. Resolve Blocker -> Should move back to IN_PROGRESS
		blockerResolved := blocker
		blockerResolved.IsResolved = sql.NullBool{Bool: true, Valid: true}
		taskBlocked := taskInProgress
		taskBlocked.Status = string(models.TaskStatusBlocked)

		mockBlockerRepo.On("FindById", ctx, int32(10)).Return(blocker, nil).Once()
		mockBlockerRepo.On("Update", ctx, int32(10), taskId, subTaskID, "Wait for info", true).Return(blockerResolved, nil).Once()
		// Sync
		mockTaskRepo.On("FindById", ctx, taskId).Return(taskBlocked, nil).Once()
		mockBlockerRepo.On("FindByTaskId", ctx, taskId).Return([]db.Blocker{blockerResolved}, nil).Once()
		mockTaskRepo.On("Update", ctx, taskId, "Task 1", string(models.TaskStatusInProgress)).Return(taskInProgress, nil).Once()

		_, err = blockerService.Resolve(ctx, 10, taskId)
		assert.NoError(t, err)

		// 5. Complete Task -> Should move to DONE
		mockTaskRepo.On("FindById", ctx, taskId).Return(taskInProgress, nil).Once()
		// Policy check
		mockBlockerRepo.On("FindByTaskId", ctx, taskId).Return([]db.Blocker{blockerResolved}, nil).Once()
		mockSubTaskRepo.On("FindByTaskId", ctx, taskId).Return([]db.SubTask{}, nil).Once()
		// Update to DONE
		mockTaskRepo.On("Update", ctx, taskId, "Task 1", string(models.TaskStatusDone)).Return(db.Task{}, nil).Once()

		err = taskService.CompleteTask(ctx, taskId)
		assert.NoError(t, err)
	})

	t.Run("should reopen task and subtask when blocker is reopened", func(t *testing.T) {
		mTaskRepo := new(MockTaskRepository)
		mSubTaskRepo := new(MockSubTaskRepository)
		mBlockerRepo := new(MockBlockerRepository)

		tSync := NewTaskSynchronizerService(mTaskRepo, mBlockerRepo)
		stSync := NewSubTaskSynchronizerService(mSubTaskRepo, mBlockerRepo)
		tGuard := NewTaskOwnershipGuardService(mTaskRepo, mSubTaskRepo)

		bService := NewBlockerService(mBlockerRepo, tGuard, tSync, stSync)

		taskId := int32(2)
		subTaskId := int32(20)
		blockerId := int32(200)

		taskInProgress := db.Task{ID: taskId, Title: "Task 1", Status: string(models.TaskStatusInProgress)}
		subTaskCompleted := db.SubTask{ID: subTaskId, TaskID: sql.NullInt32{Int32: taskId, Valid: true}, Title: "Sub 1", IsCompleted: true}
		blocker := db.Blocker{ID: blockerId, TaskID: sql.NullInt32{Int32: taskId, Valid: true}, SubTaskID: sql.NullInt32{Int32: subTaskId, Valid: true}, Reason: "B1", IsResolved: sql.NullBool{Bool: false, Valid: true}}
		blockerResolved := blocker
		blockerResolved.IsResolved = sql.NullBool{Bool: true, Valid: true}

		// Reopen Blocker workflow
		mBlockerRepo.On("FindById", ctx, blockerId).Return(blockerResolved, nil).Once()
		mBlockerRepo.On("Update", ctx, blockerId, taskId, &subTaskId, "B1", false).Return(blocker, nil).Once()

		// Task Sync
		mTaskRepo.On("FindById", ctx, taskId).Return(taskInProgress, nil).Once()
		mBlockerRepo.On("FindByTaskId", ctx, taskId).Return([]db.Blocker{blocker}, nil).Once()
		mTaskRepo.On("Update", ctx, taskId, "Task 1", string(models.TaskStatusBlocked)).Return(db.Task{}, nil).Once()

		// SubTask Sync
		mSubTaskRepo.On("FindById", ctx, subTaskId).Return(subTaskCompleted, nil).Once()
		mBlockerRepo.On("FindByTaskIdAndSubTaskId", ctx, taskId, &subTaskId).Return([]db.Blocker{blocker}, nil).Once()
		mSubTaskRepo.On("Update", ctx, subTaskId, "Sub 1", false).Return(db.SubTask{}, nil).Once()

		_, err := bService.Reopen(ctx, blockerId, taskId)
		assert.NoError(t, err)

		mBlockerRepo.AssertExpectations(t)
		mTaskRepo.AssertExpectations(t)
		mSubTaskRepo.AssertExpectations(t)
	})
}
