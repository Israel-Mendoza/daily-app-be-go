package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"database/sql"
	"fmt"
)

type BlockerService interface {
	FindAll(ctx context.Context) ([]models.BlockerResponse, error)
	FindById(ctx context.Context, id int32) (*models.BlockerResponse, error)
	FindByTaskId(ctx context.Context, taskId int32) ([]models.BlockerResponse, error)
	Save(ctx context.Context, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error)
	Replace(ctx context.Context, id int32, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error)
	Update(ctx context.Context, id int32, taskId int32, req models.UpdateBlockerRequest) (models.BlockerResponse, error)
	DeleteById(ctx context.Context, id int32) error
	Resolve(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error)
	Reopen(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error)
}

type blockerService struct {
	blockerRepo                repository.BlockerRepository
	taskOwnershipGuardService  TaskOwnershipGuardService
	taskSynchronizerService    TaskSynchronizerService
	subTaskSynchronizerService SubTaskSynchronizerService
}

func NewBlockerService(
	blockerRepo repository.BlockerRepository,
	taskOwnershipGuardService TaskOwnershipGuardService,
	taskSynchronizerService TaskSynchronizerService,
	subTaskSynchronizerService SubTaskSynchronizerService,
) BlockerService {
	return &blockerService{
		blockerRepo:                blockerRepo,
		taskOwnershipGuardService:  taskOwnershipGuardService,
		taskSynchronizerService:    taskSynchronizerService,
		subTaskSynchronizerService: subTaskSynchronizerService,
	}
}

func (s *blockerService) FindAll(ctx context.Context) ([]models.BlockerResponse, error) {
	blockers, err := s.blockerRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]models.BlockerResponse, len(blockers))
	for i, b := range blockers {
		res[i] = toBlockerResponse(b)
	}
	return res, nil
}

func (s *blockerService) FindById(ctx context.Context, id int32) (*models.BlockerResponse, error) {
	b, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toBlockerResponse(b)
	return &resp, nil
}

func (s *blockerService) FindByTaskId(ctx context.Context, taskId int32) ([]models.BlockerResponse, error) {
	blockers, err := s.blockerRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}
	res := make([]models.BlockerResponse, len(blockers))
	for i, b := range blockers {
		res[i] = toBlockerResponse(b)
	}
	return res, nil
}

func (s *blockerService) Save(ctx context.Context, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error) {
	if _, err := s.taskOwnershipGuardService.EnsureTaskExists(ctx, taskId); err != nil {
		return models.BlockerResponse{}, err
	}

	if _, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID); err != nil {
		return models.BlockerResponse{}, err
	}

	b, err := s.blockerRepo.Create(ctx, taskId, req.SubTaskID, req.Reason, false)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, req.SubTaskID)

	return toBlockerResponse(b), nil
}

func (s *blockerService) Replace(ctx context.Context, id int32, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error) {
	existing, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.BlockerResponse{}, fmt.Errorf("blocker does not belong to the specified task")
	}

	if _, err := s.taskOwnershipGuardService.EnsureTaskExists(ctx, taskId); err != nil {
		return models.BlockerResponse{}, err
	}

	if _, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID); err != nil {
		return models.BlockerResponse{}, err
	}

	updated, err := s.blockerRepo.Update(ctx, id, taskId, req.SubTaskID, req.Reason, existing.IsResolved.Bool)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, req.SubTaskID)

	return toBlockerResponse(updated), nil
}

func (s *blockerService) Update(ctx context.Context, id int32, taskId int32, req models.UpdateBlockerRequest) (models.BlockerResponse, error) {
	existing, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.BlockerResponse{}, fmt.Errorf("blocker does not belong to the specified task")
	}

	reason := existing.Reason
	if req.Reason != nil {
		reason = *req.Reason
	}

	isResolved := existing.IsResolved.Bool
	if req.IsResolved != nil {
		isResolved = *req.IsResolved
	}

	previousSubTaskId := toInt32Ptr(existing.SubTaskID)
	currentSubTaskId := existing.SubTaskID

	if req.SubTaskID != nil {
		_, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID)
		if err != nil {
			return models.BlockerResponse{}, err
		}
		currentSubTaskId = sql.NullInt32{Int32: *req.SubTaskID, Valid: true}
	}

	updated, err := s.blockerRepo.Update(ctx, id, taskId, toInt32Ptr(currentSubTaskId), reason, isResolved)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, previousSubTaskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, toInt32Ptr(currentSubTaskId))

	return toBlockerResponse(updated), nil
}

func (s *blockerService) DeleteById(ctx context.Context, id int32) error {
	existing, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	taskId := existing.TaskID.Int32
	subTaskId := toInt32Ptr(existing.SubTaskID)

	if err := s.blockerRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, subTaskId)

	return nil
}

func (s *blockerService) Resolve(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error) {
	existing, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.BlockerResponse{}, fmt.Errorf("blocker does not belong to the specified task")
	}

	updated, err := s.blockerRepo.Update(ctx, id, taskId, toInt32Ptr(existing.SubTaskID), existing.Reason, true)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, toInt32Ptr(existing.SubTaskID))

	return toBlockerResponse(updated), nil
}

func (s *blockerService) Reopen(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error) {
	existing, err := s.blockerRepo.FindById(ctx, id)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.BlockerResponse{}, fmt.Errorf("blocker does not belong to the specified task")
	}

	updated, err := s.blockerRepo.Update(ctx, id, taskId, toInt32Ptr(existing.SubTaskID), existing.Reason, false)
	if err != nil {
		return models.BlockerResponse{}, err
	}

	_ = s.taskSynchronizerService.SyncTaskWithBlockers(ctx, taskId)
	_ = s.subTaskSynchronizerService.SyncSubTaskWithBlockers(ctx, taskId, toInt32Ptr(existing.SubTaskID))

	return toBlockerResponse(updated), nil
}

func toBlockerResponse(b db.Blocker) models.BlockerResponse {
	return models.BlockerResponse{
		ID:         b.ID,
		TaskID:     b.TaskID.Int32,
		SubTaskID:  toInt32Ptr(b.SubTaskID),
		Reason:     b.Reason,
		IsResolved: b.IsResolved.Bool,
		CreatedAt:  b.CreatedAt.String(),
		UpdatedAt:  b.UpdatedAt.String(),
	}
}

func toInt32Ptr(n sql.NullInt32) *int32 {
	if n.Valid {
		return &n.Int32
	}
	return nil
}
