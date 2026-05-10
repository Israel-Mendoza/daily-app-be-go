package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"fmt"
	"slices"
)

type SubTaskPolicyService interface {
	EnsureCanCreateOrAssignToTask(taskStatus string) error
	EnsureCanCompleteSubTask(ctx context.Context, subTask db.SubTask) error
}

type subTaskPolicyService struct {
	blockerRepo repository.BlockerRepository
}

func NewSubTaskPolicyService(blockerRepo repository.BlockerRepository) SubTaskPolicyService {
	return &subTaskPolicyService{
		blockerRepo: blockerRepo,
	}
}

func (s *subTaskPolicyService) EnsureCanCreateOrAssignToTask(taskStatus string) error {
	if taskStatus == string(models.TaskStatusDone) || taskStatus == string(models.TaskStatusCanceled) {
		return fmt.Errorf("cannot create or assign subtasks to a task that is DONE or CANCELED")
	}
	return nil
}

func (s *subTaskPolicyService) EnsureCanCompleteSubTask(ctx context.Context, subTask db.SubTask) error {
	blockers, err := s.blockerRepo.FindByTaskIdAndSubTaskId(ctx, subTask.TaskID.Int32, &subTask.ID)
	if err != nil {
		return err
	}

	hasUnresolved := slices.ContainsFunc(blockers, func(b db.Blocker) bool {
		return !b.IsResolved.Bool
	})

	if hasUnresolved {
		return fmt.Errorf("cannot complete subtask with unresolved blockers")
	}

	return nil
}

type TaskPolicyService interface {
	EnsureCanTransitionStatus(ctx context.Context, task db.Task, targetStatus models.TaskStatus) error
}

type taskPolicyService struct {
	blockerRepo repository.BlockerRepository
	subTaskRepo repository.SubTaskRepository
}

func NewTaskPolicyService(blockerRepo repository.BlockerRepository, subTaskRepo repository.SubTaskRepository) TaskPolicyService {
	return &taskPolicyService{
		blockerRepo: blockerRepo,
		subTaskRepo: subTaskRepo,
	}
}

func (s *taskPolicyService) EnsureCanTransitionStatus(ctx context.Context, task db.Task, targetStatus models.TaskStatus) error {
	if err := s.ensureCanUnblockTask(ctx, task, targetStatus); err != nil {
		return err
	}
	if err := s.ensureCanCompleteTask(ctx, task, targetStatus); err != nil {
		return err
	}
	return nil
}

func (s *taskPolicyService) ensureCanUnblockTask(ctx context.Context, task db.Task, targetStatus models.TaskStatus) error {
	currentStatus := models.TaskStatus(task.Status)

	if targetStatus == models.TaskStatusCanceled {
		return nil
	}

	if currentStatus == models.TaskStatusBlocked && targetStatus != models.TaskStatusBlocked {
		blockers, err := s.blockerRepo.FindByTaskId(ctx, task.ID)
		if err != nil {
			return err
		}

		hasUnresolved := slices.ContainsFunc(blockers, func(b db.Blocker) bool {
			return !b.IsResolved.Bool
		})

		if hasUnresolved {
			return fmt.Errorf("cannot change task status from BLOCKED while unresolved blockers exist")
		}
	}

	return nil
}

func (s *taskPolicyService) ensureCanCompleteTask(ctx context.Context, task db.Task, targetStatus models.TaskStatus) error {
	if targetStatus == models.TaskStatusDone {
		subTasks, err := s.subTaskRepo.FindByTaskId(ctx, task.ID)
		if err != nil {
			return err
		}

		hasIncomplete := slices.ContainsFunc(subTasks, func(st db.SubTask) bool {
			return !st.IsCompleted
		})

		if hasIncomplete {
			return fmt.Errorf("cannot change task status to DONE while incomplete subtasks exist")
		}
	}
	return nil
}
