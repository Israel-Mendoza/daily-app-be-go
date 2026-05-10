package service

import (
	"context"
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindById(ctx context.Context, id int32) (db.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]db.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Task), args.Error(1)
}

func (m *MockTaskRepository) Create(ctx context.Context, title string, status string) (db.Task, error) {
	args := m.Called(ctx, title, status)
	return args.Get(0).(db.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, id int32, title string, status string) (db.Task, error) {
	args := m.Called(ctx, id, title, status)
	return args.Get(0).(db.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockBlockerService struct {
	mock.Mock
}

func (m *MockBlockerService) FindAll(ctx context.Context) ([]models.BlockerResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) FindById(ctx context.Context, id int32) (*models.BlockerResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) FindByTaskId(ctx context.Context, taskId int32) ([]models.BlockerResponse, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) Save(ctx context.Context, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error) {
	args := m.Called(ctx, taskId, req)
	return args.Get(0).(models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) Replace(ctx context.Context, id int32, taskId int32, req models.CreateBlockerRequest) (models.BlockerResponse, error) {
	args := m.Called(ctx, id, taskId, req)
	return args.Get(0).(models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) Update(ctx context.Context, id int32, taskId int32, req models.UpdateBlockerRequest) (models.BlockerResponse, error) {
	args := m.Called(ctx, id, taskId, req)
	return args.Get(0).(models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) DeleteById(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBlockerService) Resolve(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error) {
	args := m.Called(ctx, id, taskId)
	return args.Get(0).(models.BlockerResponse), args.Error(1)
}

func (m *MockBlockerService) Reopen(ctx context.Context, id int32, taskId int32) (models.BlockerResponse, error) {
	args := m.Called(ctx, id, taskId)
	return args.Get(0).(models.BlockerResponse), args.Error(1)
}

type MockSubTaskService struct {
	mock.Mock
}

func (m *MockSubTaskService) FindAll(ctx context.Context) ([]models.SubTaskResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) FindById(ctx context.Context, id int32, taskId int32) (*models.SubTaskResponse, error) {
	args := m.Called(ctx, id, taskId)
	return args.Get(0).(*models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) FindAllByTaskId(ctx context.Context, taskId int32) ([]models.SubTaskResponse, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) Save(ctx context.Context, taskId int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error) {
	args := m.Called(ctx, taskId, req)
	return args.Get(0).(models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) DeleteById(ctx context.Context, taskId int32, id int32) error {
	args := m.Called(ctx, taskId, id)
	return args.Error(0)
}

func (m *MockSubTaskService) Update(ctx context.Context, taskId int32, id int32, req models.CreateAndUpdateSubTaskRequest) (models.SubTaskResponse, error) {
	args := m.Called(ctx, taskId, id, req)
	return args.Get(0).(models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) Complete(ctx context.Context, taskId int32, id int32) (models.SubTaskResponse, error) {
	args := m.Called(ctx, taskId, id)
	return args.Get(0).(models.SubTaskResponse), args.Error(1)
}

func (m *MockSubTaskService) Reopen(ctx context.Context, taskId int32, id int32) error {
	args := m.Called(ctx, taskId, id)
	return args.Error(0)
}

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) FindByTaskId(ctx context.Context, taskId int32) ([]models.NoteResponse, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]models.NoteResponse), args.Error(1)
}

func (m *MockNoteService) FindById(ctx context.Context, id int32, taskId int32) (*models.NoteResponse, error) {
	args := m.Called(ctx, id, taskId)
	return args.Get(0).(*models.NoteResponse), args.Error(1)
}

func (m *MockNoteService) Save(ctx context.Context, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error) {
	args := m.Called(ctx, taskId, req)
	return args.Get(0).(models.NoteResponse), args.Error(1)
}

func (m *MockNoteService) Replace(ctx context.Context, id int32, taskId int32, req models.CreateNoteRequest) (models.NoteResponse, error) {
	args := m.Called(ctx, id, taskId, req)
	return args.Get(0).(models.NoteResponse), args.Error(1)
}

func (m *MockNoteService) Update(ctx context.Context, id int32, taskId int32, req models.UpdateNoteRequest) (models.NoteResponse, error) {
	args := m.Called(ctx, id, taskId, req)
	return args.Get(0).(models.NoteResponse), args.Error(1)
}

func (m *MockNoteService) DeleteById(ctx context.Context, taskId int32, id int32) error {
	args := m.Called(ctx, taskId, id)
	return args.Error(0)
}

type MockTaskPolicyService struct {
	mock.Mock
}

func (m *MockTaskPolicyService) EnsureCanTransitionStatus(ctx context.Context, task db.Task, targetStatus models.TaskStatus) error {
	args := m.Called(ctx, task, targetStatus)
	return args.Error(0)
}

type MockBlockerRepository struct {
	mock.Mock
}

func (m *MockBlockerRepository) FindById(ctx context.Context, id int32) (db.Blocker, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) FindAll(ctx context.Context) ([]db.Blocker, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.Blocker, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) FindByTaskIdAndSubTaskId(ctx context.Context, taskId int32, subTaskId *int32) ([]db.Blocker, error) {
	args := m.Called(ctx, taskId, subTaskId)
	return args.Get(0).([]db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) Create(ctx context.Context, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error) {
	args := m.Called(ctx, taskId, subTaskId, reason, isResolved)
	return args.Get(0).(db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, reason string, isResolved bool) (db.Blocker, error) {
	args := m.Called(ctx, id, taskId, subTaskId, reason, isResolved)
	return args.Get(0).(db.Blocker), args.Error(1)
}

func (m *MockBlockerRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTaskOwnershipGuardService struct {
	mock.Mock
}

func (m *MockTaskOwnershipGuardService) EnsureTaskExists(ctx context.Context, taskId int32) (db.Task, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).(db.Task), args.Error(1)
}

func (m *MockTaskOwnershipGuardService) EnsureSubTaskBelongsToTask(ctx context.Context, taskId int32, subTaskId *int32) (*db.SubTask, error) {
	args := m.Called(ctx, taskId, subTaskId)
	return args.Get(0).(*db.SubTask), args.Error(1)
}

type MockTaskSynchronizerService struct {
	mock.Mock
}

func (m *MockTaskSynchronizerService) SyncTaskWithBlockers(ctx context.Context, taskId int32) error {
	args := m.Called(ctx, taskId)
	return args.Error(0)
}

func (m *MockTaskSynchronizerService) SyncTaskWithNewSubTasks(ctx context.Context, taskId int32) error {
	args := m.Called(ctx, taskId)
	return args.Error(0)
}

type MockSubTaskSynchronizerService struct {
	mock.Mock
}

func (m *MockSubTaskSynchronizerService) SyncSubTaskWithBlockers(ctx context.Context, taskId int32, subTaskId *int32) error {
	args := m.Called(ctx, taskId, subTaskId)
	return args.Error(0)
}

type MockSubTaskRepository struct {
	mock.Mock
}

func (m *MockSubTaskRepository) FindById(ctx context.Context, id int32) (db.SubTask, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.SubTask), args.Error(1)
}

func (m *MockSubTaskRepository) FindAll(ctx context.Context) ([]db.SubTask, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.SubTask), args.Error(1)
}

func (m *MockSubTaskRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.SubTask, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]db.SubTask), args.Error(1)
}

func (m *MockSubTaskRepository) Create(ctx context.Context, taskId int32, title string, isCompleted bool) (db.SubTask, error) {
	args := m.Called(ctx, taskId, title, isCompleted)
	return args.Get(0).(db.SubTask), args.Error(1)
}

func (m *MockSubTaskRepository) Update(ctx context.Context, id int32, title string, isCompleted bool) (db.SubTask, error) {
	args := m.Called(ctx, id, title, isCompleted)
	return args.Get(0).(db.SubTask), args.Error(1)
}

func (m *MockSubTaskRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockSubTaskPolicyService struct {
	mock.Mock
}

func (m *MockSubTaskPolicyService) EnsureCanCreateOrAssignToTask(taskStatus string) error {
	args := m.Called(taskStatus)
	return args.Error(0)
}

func (m *MockSubTaskPolicyService) EnsureCanCompleteSubTask(ctx context.Context, subTask db.SubTask) error {
	args := m.Called(ctx, subTask)
	return args.Error(0)
}

type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) FindById(ctx context.Context, id int32) (db.TaskNote, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.TaskNote), args.Error(1)
}

func (m *MockNoteRepository) FindAll(ctx context.Context) ([]db.TaskNote, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.TaskNote), args.Error(1)
}

func (m *MockNoteRepository) FindByTaskId(ctx context.Context, taskId int32) ([]db.TaskNote, error) {
	args := m.Called(ctx, taskId)
	return args.Get(0).([]db.TaskNote), args.Error(1)
}

func (m *MockNoteRepository) Create(ctx context.Context, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error) {
	args := m.Called(ctx, taskId, subTaskId, content, category)
	return args.Get(0).(db.TaskNote), args.Error(1)
}

func (m *MockNoteRepository) Update(ctx context.Context, id int32, taskId int32, subTaskId *int32, content string, category string) (db.TaskNote, error) {
	args := m.Called(ctx, id, taskId, subTaskId, content, category)
	return args.Get(0).(db.TaskNote), args.Error(1)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
