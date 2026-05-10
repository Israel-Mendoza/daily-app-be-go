package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/repository"
	"fmt"
)

type TaskOwnershipGuardService interface {
	EnsureTaskExists(ctx context.Context, taskId int32) (db.Task, error)
	EnsureSubTaskBelongsToTask(ctx context.Context, taskId int32, subTaskId *int32) (*db.SubTask, error)
}

type taskOwnershipGuardService struct {
	taskRepo    repository.TaskRepository
	subTaskRepo repository.SubTaskRepository
}

func NewTaskOwnershipGuardService(taskRepo repository.TaskRepository, subTaskRepo repository.SubTaskRepository) TaskOwnershipGuardService {
	return &taskOwnershipGuardService{
		taskRepo:    taskRepo,
		subTaskRepo: subTaskRepo,
	}
}

func (s *taskOwnershipGuardService) EnsureTaskExists(ctx context.Context, taskId int32) (db.Task, error) {
	task, err := s.taskRepo.FindById(ctx, taskId)
	if err != nil {
		return db.Task{}, fmt.Errorf("task not found: %w", err)
	}
	return task, nil
}

func (s *taskOwnershipGuardService) EnsureSubTaskBelongsToTask(ctx context.Context, taskId int32, subTaskId *int32) (*db.SubTask, error) {
	if subTaskId == nil {
		return nil, nil
	}

	subTask, err := s.subTaskRepo.FindById(ctx, *subTaskId)
	if err != nil {
		return nil, fmt.Errorf("subtask not found: %w", err)
	}

	if subTask.TaskID.Int32 != taskId {
		return nil, fmt.Errorf("subtask does not belong to the specified task")
	}

	return &subTask, nil
}
