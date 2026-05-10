package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"fmt"
)

type TaskService interface {
	FindAll(ctx context.Context) ([]models.TaskResponse, error)
	FindById(ctx context.Context, id int32) (*models.TaskResponse, error)
	Save(ctx context.Context, req models.CreateTaskRequest) (models.TaskResponse, error)
	Update(ctx context.Context, id int32, req models.CreateTaskRequest) (models.TaskResponse, error)
	DeleteById(ctx context.Context, id int32) error
	StartTask(ctx context.Context, id int32) error
	CompleteTask(ctx context.Context, id int32) error
	ReopenTask(ctx context.Context, id int32) error
	CancelTask(ctx context.Context, id int32) error
	GetTaskDetailsById(ctx context.Context, taskId int32) (*models.TaskDetailsResponse, error)
}

type taskService struct {
	taskRepo          repository.TaskRepository
	blockerService    BlockerService
	subTaskService    SubTaskService
	noteService       NoteService
	taskPolicyService TaskPolicyService
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	blockerService BlockerService,
	subTaskService SubTaskService,
	noteService NoteService,
	taskPolicyService TaskPolicyService,
) TaskService {
	return &taskService{
		taskRepo:          taskRepo,
		blockerService:    blockerService,
		subTaskService:    subTaskService,
		noteService:       noteService,
		taskPolicyService: taskPolicyService,
	}
}

func (s *taskService) FindAll(ctx context.Context) ([]models.TaskResponse, error) {
	tasks, err := s.taskRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.TaskResponse, len(tasks))
	for i, t := range tasks {
		res[i] = toTaskResponse(t)
	}
	return res, nil
}

func (s *taskService) FindById(ctx context.Context, id int32) (*models.TaskResponse, error) {
	task, err := s.taskRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toTaskResponse(task)
	return &resp, nil
}

func (s *taskService) Save(ctx context.Context, req models.CreateTaskRequest) (models.TaskResponse, error) {
	task, err := s.taskRepo.Create(ctx, req.Title, string(models.TaskStatusTodo))
	if err != nil {
		return models.TaskResponse{}, err
	}
	return toTaskResponse(task), nil
}

func (s *taskService) Update(ctx context.Context, id int32, req models.CreateTaskRequest) (models.TaskResponse, error) {
	existing, err := s.getTaskFromRepository(ctx, id)
	if err != nil {
		return models.TaskResponse{}, err
	}

	task, err := s.taskRepo.Update(ctx, id, req.Title, existing.Status)
	if err != nil {
		return models.TaskResponse{}, err
	}
	return toTaskResponse(task), nil
}

func (s *taskService) DeleteById(ctx context.Context, id int32) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *taskService) StartTask(ctx context.Context, id int32) error {
	existing, err := s.getTaskFromRepository(ctx, id)
	if err != nil {
		return err
	}

	if err := s.taskPolicyService.EnsureCanTransitionStatus(ctx, existing, models.TaskStatusInProgress); err != nil {
		return err
	}

	_, err = s.taskRepo.Update(ctx, id, existing.Title, string(models.TaskStatusInProgress))
	return err
}

func (s *taskService) CompleteTask(ctx context.Context, id int32) error {
	existing, err := s.getTaskFromRepository(ctx, id)
	if err != nil {
		return err
	}

	if err := s.taskPolicyService.EnsureCanTransitionStatus(ctx, existing, models.TaskStatusDone); err != nil {
		return err
	}

	_, err = s.taskRepo.Update(ctx, id, existing.Title, string(models.TaskStatusDone))
	return err
}

func (s *taskService) ReopenTask(ctx context.Context, id int32) error {
	existing, err := s.getTaskFromRepository(ctx, id)
	if err != nil {
		return err
	}

	if err := s.taskPolicyService.EnsureCanTransitionStatus(ctx, existing, models.TaskStatusTodo); err != nil {
		return err
	}

	_, err = s.taskRepo.Update(ctx, id, existing.Title, string(models.TaskStatusTodo))
	return err
}

func (s *taskService) CancelTask(ctx context.Context, id int32) error {
	existing, err := s.getTaskFromRepository(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.taskRepo.Update(ctx, id, existing.Title, string(models.TaskStatusCanceled))
	return err
}

func (s *taskService) GetTaskDetailsById(ctx context.Context, taskId int32) (*models.TaskDetailsResponse, error) {
	task, err := s.getTaskFromRepository(ctx, taskId)
	if err != nil {
		return nil, err
	}

	subTasks, err := s.subTaskService.FindAllByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}

	blockers, err := s.blockerService.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}

	notes, err := s.noteService.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}

	return &models.TaskDetailsResponse{
		ID:        task.ID,
		Title:     task.Title,
		Status:    task.Status,
		CreatedAt: task.CreatedAt.String(),
		UpdatedAt: task.UpdatedAt.String(),
		Notes:     notes,
		SubTasks:  subTasks,
		Blockers:  blockers,
	}, nil
}

func (s *taskService) getTaskFromRepository(ctx context.Context, taskId int32) (db.Task, error) {
	task, err := s.taskRepo.FindById(ctx, taskId)
	if err != nil {
		return db.Task{}, fmt.Errorf("task not found: %w", err)
	}
	return task, nil
}

func toTaskResponse(t db.Task) models.TaskResponse {
	return models.TaskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.String(),
		UpdatedAt: t.UpdatedAt.String(),
	}
}
