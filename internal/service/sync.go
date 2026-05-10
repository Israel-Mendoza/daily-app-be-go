package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"fmt"
	"slices"
)

type TaskSynchronizerService interface {
	SyncTaskWithBlockers(ctx context.Context, taskId int32) error
	SyncTaskWithNewSubTasks(ctx context.Context, taskId int32) error
}

type taskSynchronizerService struct {
	taskRepo    repository.TaskRepository
	blockerRepo repository.BlockerRepository
}

func NewTaskSynchronizerService(taskRepo repository.TaskRepository, blockerRepo repository.BlockerRepository) TaskSynchronizerService {
	return &taskSynchronizerService{
		taskRepo:    taskRepo,
		blockerRepo: blockerRepo,
	}
}

func (s *taskSynchronizerService) SyncTaskWithBlockers(ctx context.Context, taskId int32) error {
	task, err := s.taskRepo.FindById(ctx, taskId)
	if err != nil {
		return err
	}

	if task.Status == string(models.TaskStatusCanceled) {
		return nil
	}

	blockers, err := s.blockerRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return err
	}

	hasUnresolved := slices.ContainsFunc(blockers, func(b db.Blocker) bool {
		return !b.IsResolved.Bool
	})

	newStatus := task.Status
	if hasUnresolved {
		newStatus = string(models.TaskStatusBlocked)
	} else if task.Status == string(models.TaskStatusDone) {
		// Keep as DONE
	} else {
		newStatus = string(models.TaskStatusInProgress)
	}

	if newStatus != task.Status {
		_, err = s.taskRepo.Update(ctx, task.ID, task.Title, newStatus)
		return err
	}

	return nil
}

func (s *taskSynchronizerService) SyncTaskWithNewSubTasks(ctx context.Context, taskId int32) error {
	task, err := s.taskRepo.FindById(ctx, taskId)
	if err != nil {
		return err
	}

	if task.Status == string(models.TaskStatusDone) {
		_, err = s.taskRepo.Update(ctx, task.ID, task.Title, string(models.TaskStatusInProgress))
		return err
	}

	return nil
}

type SubTaskSynchronizerService interface {
	SyncSubTaskWithBlockers(ctx context.Context, taskId int32, subTaskId *int32) error
}

type subTaskSynchronizerService struct {
	subTaskRepo repository.SubTaskRepository
	blockerRepo repository.BlockerRepository
}

func NewSubTaskSynchronizerService(subTaskRepo repository.SubTaskRepository, blockerRepo repository.BlockerRepository) SubTaskSynchronizerService {
	return &subTaskSynchronizerService{
		subTaskRepo: subTaskRepo,
		blockerRepo: blockerRepo,
	}
}

func (s *subTaskSynchronizerService) SyncSubTaskWithBlockers(ctx context.Context, taskId int32, subTaskId *int32) error {
	if subTaskId == nil {
		return nil
	}

	subTask, err := s.subTaskRepo.FindById(ctx, *subTaskId)
	if err != nil {
		return err
	}

	if subTask.TaskID.Int32 != taskId {
		return fmt.Errorf("subtask does not belong to the specified task")
	}

	blockers, err := s.blockerRepo.FindByTaskIdAndSubTaskId(ctx, taskId, subTaskId)
	if err != nil {
		return err
	}

	hasUnresolved := slices.ContainsFunc(blockers, func(b db.Blocker) bool {
		return !b.IsResolved.Bool
	})

	if hasUnresolved && subTask.IsCompleted {
		_, err = s.subTaskRepo.Update(ctx, subTask.ID, subTask.Title, false)
		return err
	}

	return nil
}
