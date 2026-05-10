package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/repository"
	"fmt"
)

type NoteService interface {
	FindByTaskId(ctx context.Context, taskId int32) ([]models.NoteResponse, error)
	FindById(ctx context.Context, id int32, taskId int32) (*models.NoteResponse, error)
	Save(ctx context.Context, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error)
	Replace(ctx context.Context, id int32, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error)
	Update(ctx context.Context, id int32, taskId int32, req models.UpdateNoteRequest) (models.NoteResponse, error)
	DeleteById(ctx context.Context, taskId int32, id int32) error
}

type noteService struct {
	noteRepo                  repository.NoteRepository
	taskOwnershipGuardService TaskOwnershipGuardService
}

func NewNoteService(noteRepo repository.NoteRepository, guard TaskOwnershipGuardService) NoteService {
	return &noteService{
		noteRepo:                  noteRepo,
		taskOwnershipGuardService: guard,
	}
}

func (s *noteService) FindByTaskId(ctx context.Context, taskId int32) ([]models.NoteResponse, error) {
	notes, err := s.noteRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err
	}
	res := make([]models.NoteResponse, len(notes))
	for i, n := range notes {
		res[i] = toNoteResponse(n)
	}
	return res, nil
}

func (s *noteService) FindById(ctx context.Context, id int32, taskId int32) (*models.NoteResponse, error) {
	n, err := s.noteRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if n.TaskID.Int32 != taskId {
		return nil, nil
	}
	resp := toNoteResponse(n)
	return &resp, nil
}

func (s *noteService) Save(ctx context.Context, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error) {
	if _, err := s.taskOwnershipGuardService.EnsureTaskExists(ctx, taskId); err != nil {
		return models.NoteResponse{}, err
	}

	if _, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID); err != nil {
		return models.NoteResponse{}, err
	}

	n, err := s.noteRepo.Create(ctx, taskId, req.SubTaskID, req.Content, req.Category)
	if err != nil {
		return models.NoteResponse{}, err
	}
	return toNoteResponse(n), nil
}

func (s *noteService) Replace(ctx context.Context, id int32, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error) {
	existing, err := s.noteRepo.FindById(ctx, id)
	if err != nil {
		return models.NoteResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.NoteResponse{}, fmt.Errorf("note does not belong to the specified task")
	}

	if _, err := s.taskOwnershipGuardService.EnsureTaskExists(ctx, taskId); err != nil {
		return models.NoteResponse{}, err
	}

	if _, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID); err != nil {
		return models.NoteResponse{}, err
	}

	updated, err := s.noteRepo.Update(ctx, id, taskId, req.SubTaskID, req.Content, req.Category)
	if err != nil {
		return models.NoteResponse{}, err
	}
	return toNoteResponse(updated), nil
}

func (s *noteService) Update(ctx context.Context, id int32, taskId int32, req models.UpdateNoteRequest) (models.NoteResponse, error) {
	existing, err := s.noteRepo.FindById(ctx, id)
	if err != nil {
		return models.NoteResponse{}, err
	}

	if existing.TaskID.Int32 != taskId {
		return models.NoteResponse{}, fmt.Errorf("note does not belong to the specified task")
	}

	content := existing.Content
	if req.Content != nil {
		content = *req.Content
	}

	category := existing.Category
	if req.Category != nil {
		category = *req.Category
	}

	subTaskId := toInt32Ptr(existing.SubTaskID)
	if req.SubTaskID != nil {
		if _, err := s.taskOwnershipGuardService.EnsureSubTaskBelongsToTask(ctx, taskId, req.SubTaskID); err != nil {
			return models.NoteResponse{}, err
		}
		subTaskId = req.SubTaskID
	}

	updated, err := s.noteRepo.Update(ctx, id, taskId, subTaskId, content, category)
	if err != nil {
		return models.NoteResponse{}, err
	}
	return toNoteResponse(updated), nil
}

func (s *noteService) DeleteById(ctx context.Context, taskId int32, id int32) error {
	existing, err := s.noteRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if existing.TaskID.Int32 != taskId {
		return fmt.Errorf("note does not belong to the specified task")
	}
	return s.noteRepo.Delete(ctx, id)
}

func toNoteResponse(n db.TaskNote) models.NoteResponse {
	return models.NoteResponse{
		ID:        n.ID,
		TaskID:    n.TaskID.Int32,
		SubTaskID: toInt32Ptr(n.SubTaskID),
		Content:   n.Content,
		Category:  n.Category,
		CreatedAt: n.CreatedAt.String(),
		UpdatedAt: n.UpdatedAt.String(),
	}
}
