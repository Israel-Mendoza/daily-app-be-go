package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"fmt"
)

type SubTaskService interface {
	FindAll(ctx context.Context) ([]models.SubTaskResponse, error)
	FindById(ctx context.Context, id int32, taskId int32) (*models.SubTaskResponse, error)
	FindAllByTaskId(ctx context.Context, taskId int32) ([]models.SubTaskResponse, error)
	Save(ctx context.Context, taskId int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error)
	DeleteById(ctx context.Context, taskId int32, id int32) error
	Update(ctx context.Context, taskId int32, id int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error)
	Complete(ctx context.Context, taskId int32, id int32) (models.SubTaskResponse, error)
	Reopen(ctx context.Context, taskId int32, id int32) error
}

type subTaskService struct {
	subTaskRepo             repository.SubTaskRepository
	taskRepo                repository.TaskRepository
	taskSynchronizerService TaskSynchronizerService
	subTaskPolicyService    SubTaskPolicyService
}

func NewSubTaskService(
	subTaskRepo repository.SubTaskRepository,
	taskRepo repository.TaskRepository,
	taskSynchronizerService TaskSynchronizerService,
	subTaskPolicyService SubTaskPolicyService,
) SubTaskService {
	return &subTaskService{
		subTaskRepo:             subTaskRepo,
		taskRepo:                taskRepo,
		taskSynchronizerService: taskSynchronizerService,
		subTaskPolicyService:    subTaskPolicyService,
	}
}

func (s *subTaskService) FindAll(ctx context.Context) ([]models.SubTaskResponse, error) {
	subTasks, err := s.subTaskRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]models.SubTaskResponse, len(subTasks))
	for i, st := range subTasks {
		res[i] = toSubTaskResponse(st)
	}
	return res, nil
}

func (s *subTaskService) FindById(ctx context.Context, id int32, taskId int32) (*models.SubTaskResponse, error) {
	st, err := s.subTaskRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if st.TaskID.Int32 != taskId {
		return nil, nil
	}
	resp := toSubTaskResponse(st)
	return &resp, nil
}

func (s *subTaskService) FindAllByTaskId(ctx context.Context, taskId int32) ([]models.SubTaskResponse, error) {
	subTasks, err := s.subTaskRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}
	res := make([]models.SubTaskResponse, len(subTasks))
	for i, st := range subTasks {
		res[i] = toSubTaskResponse(st)
	}
	return res, nil
}

func (s *subTaskService) Save(ctx context.Context, taskId int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error) {
	_, err := s.taskRepo.FindById(ctx, taskId)
	if err != nil {
		return models.SubTaskResponse{}, fmt.Errorf("task not found: %w", err)
	}

	st, err := s.subTaskRepo.Create(ctx, taskId, req.Title, false)
	if err != nil {
		return models.SubTaskResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithNewSubTasks(ctx, taskId)

	return toSubTaskResponse(st), nil
}

func (s *subTaskService) DeleteById(ctx context.Context, taskId int32, id int32) error {
	st, err := s.getSubTaskAndEnsureTaskOwnership(ctx, taskId, id)
	if err != nil {
		return err
	}
	return s.subTaskRepo.Delete(ctx, st.ID)
}

func (s *subTaskService) Update(ctx context.Context, taskId int32, id int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error) {
	st, err := s.getSubTaskAndEnsureTaskOwnership(ctx, taskId, id)
	if err != nil {
		return models.SubTaskResponse{}, err
	}

	updated, err := s.subTaskRepo.Update(ctx, id, req.Title, st.IsCompleted)
	if err != nil {
		return models.SubTaskResponse{}, err
	}
	return toSubTaskResponse(updated), nil
}

func (s *subTaskService) Complete(ctx context.Context, taskId int32, id int32) (models.SubTaskResponse, error) {
	st, err := s.getSubTaskAndEnsureTaskOwnership(ctx, taskId, id)
	if err != nil {
		return models.SubTaskResponse{}, err
	}

	if err := s.subTaskPolicyService.EnsureCanCompleteSubTask(ctx, st); err != nil {
		return models.SubTaskResponse{}, err
	}

	updated, err := s.subTaskRepo.Update(ctx, id, st.Title, true)
	if err != nil {
		return models.SubTaskResponse{}, err
	}
	return toSubTaskResponse(updated), nil
}

func (s *subTaskService) Reopen(ctx context.Context, taskId int32, id int32) error {
	st, err := s.getSubTaskAndEnsureTaskOwnership(ctx, taskId, id)
	if err != nil {
		return err
	}

	_, err = s.subTaskRepo.Update(ctx, id, st.Title, false)
	if err != nil {
		return err
	}

	_ = s.taskSynchronizerService.SyncTaskWithNewSubTasks(ctx, taskId)
	return nil
}

func (s *subTaskService) getSubTaskAndEnsureTaskOwnership(ctx context.Context, taskId int32, subTaskId int32) (db.SubTask, error) {
	st, err := s.subTaskRepo.FindById(ctx, subTaskId)
	if err != nil {
		return db.SubTask{}, fmt.Errorf("subtask not found: %w", err)
	}
	if st.TaskID.Int32 != taskId {
		return db.SubTask{}, fmt.Errorf("subtask does not belong to the specified task")
	}
	return st, nil
}

func toSubTaskResponse(st db.SubTask) models.SubTaskResponse {
	return models.SubTaskResponse{
		ID:          st.ID,
		TaskID:      st.TaskID.Int32,
		Title:       st.Title,
		IsCompleted: st.IsCompleted,
		CreatedAt:   st.CreatedAt.String(),
		UpdatedAt:   st.UpdatedAt.String(),
	}
}
