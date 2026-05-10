package tests

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"daily-app-go/db/sqlc"
	"daily-app-go/internal/repository"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type IntegrationTestSuite struct {
	suite.Suite
	pgContainer *postgres.PostgresContainer
	db          *sql.DB
	queries     *db.Queries
	repo        struct {
		task         repository.TaskRepository
		subTask      repository.SubTaskRepository
		blocker      repository.BlockerRepository
		note         repository.NoteRepository
		dailySession repository.DailySessionRepository
	}
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b))))
	initScript := filepath.Join(basepath, "init.sql")

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(initScript),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	s.NoError(err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	s.NoError(err)

	dbConn, err := sql.Open("postgres", connStr)
	s.NoError(err)

	err = dbConn.Ping()
	s.NoError(err)

	s.pgContainer = pgContainer
	s.db = dbConn
	s.queries = db.New(dbConn)

	// Initialize repositories
	s.repo.task = repository.NewTaskRepository(s.queries)
	s.repo.subTask = repository.NewSubTaskRepository(s.queries)
	s.repo.blocker = repository.NewBlockerRepository(s.queries)
	s.repo.note = repository.NewNoteRepository(s.queries)
	s.repo.dailySession = repository.NewDailySessionRepository(s.queries)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.db != nil {
		s.db.Close()
	}
	if s.pgContainer != nil {
		s.pgContainer.Terminate(ctx)
	}
}

func (s *IntegrationTestSuite) SetupTest() {
	// Clean up tables before each test to ensure isolation
	ctx := context.Background()
	_, err := s.db.ExecContext(ctx, "TRUNCATE tasks, sub_tasks, blockers, task_notes, daily_sessions RESTART IDENTITY CASCADE")
	s.NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestTaskRepository() {
	ctx := context.Background()

	// Create
	task, err := s.repo.task.Create(ctx, "Test Task", "TODO")
	s.NoError(err)
	s.Equal("Test Task", task.Title)
	s.Equal("TODO", task.Status)

	// FindById
	found, err := s.repo.task.FindById(ctx, task.ID)
	s.NoError(err)
	s.Equal(task.ID, found.ID)
	s.Equal("Test Task", found.Title)

	// Update
	updated, err := s.repo.task.Update(ctx, task.ID, "Updated Task", "IN_PROGRESS")
	s.NoError(err)
	s.Equal("Updated Task", updated.Title)
	s.Equal("IN_PROGRESS", updated.Status)

	// FindAll
	tasks, err := s.repo.task.FindAll(ctx)
	s.NoError(err)
	s.Len(tasks, 1)

	// Delete
	err = s.repo.task.Delete(ctx, task.ID)
	s.NoError(err)

	_, err = s.repo.task.FindById(ctx, task.ID)
	s.Error(err)
}

func (s *IntegrationTestSuite) TestSubTaskRepository() {
	ctx := context.Background()

	task, _ := s.repo.task.Create(ctx, "Parent Task", "TODO")

	// Create
	st, err := s.repo.subTask.Create(ctx, task.ID, "SubTask 1", false)
	s.NoError(err)
	s.Equal("SubTask 1", st.Title)
	s.False(st.IsCompleted)

	// FindByTaskId
	sts, err := s.repo.subTask.FindByTaskId(ctx, task.ID)
	s.NoError(err)
	s.Len(sts, 1)
	s.Equal(st.ID, sts[0].ID)

	// Update
	updated, err := s.repo.subTask.Update(ctx, st.ID, "Updated SubTask", true)
	s.NoError(err)
	s.True(updated.IsCompleted)
}

func (s *IntegrationTestSuite) TestBlockerRepository() {
	ctx := context.Background()

	task, _ := s.repo.task.Create(ctx, "Blocked Task", "TODO")

	// Create task blocker
	blocker, err := s.repo.blocker.Create(ctx, task.ID, nil, "Missing requirements", false)
	s.NoError(err)
	s.Equal("Missing requirements", blocker.Reason)
	s.False(blocker.IsResolved.Bool)

	// FindByTaskId
	blockers, err := s.repo.blocker.FindByTaskId(ctx, task.ID)
	s.NoError(err)
	s.Len(blockers, 1)

	// Update
	updated, err := s.repo.blocker.Update(ctx, blocker.ID, task.ID, nil, "Requirements found", true)
	s.NoError(err)
	s.True(updated.IsResolved.Bool)
}

func (s *IntegrationTestSuite) TestNoteRepository() {
	ctx := context.Background()

	task, _ := s.repo.task.Create(ctx, "Note Task", "TODO")

	// Create
	note, err := s.repo.note.Create(ctx, task.ID, nil, "Important note", "GENERAL")
	s.NoError(err)
	s.Equal("Important note", note.Content)
	s.Equal("GENERAL", note.Category)

	// FindByTaskId
	notes, err := s.repo.note.FindByTaskId(ctx, task.ID)
	s.NoError(err)
	s.Len(notes, 1)

	// Update
	updated, err := s.repo.note.Update(ctx, note.ID, task.ID, nil, "Updated note", "DECISION")
	s.NoError(err)
	s.Equal("Updated note", updated.Content)
	s.Equal("DECISION", updated.Category)
}

func (s *IntegrationTestSuite) TestDailySessionRepository() {
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	// Create
	session, err := s.repo.dailySession.Create(ctx, now, "Raw notes content", "Generated script")
	s.NoError(err)
	s.Equal("Raw notes content", session.RawNotesBlob)
	s.WithinDuration(now, session.SessionDate.UTC(), time.Second)

	// FindAll
	sessions, err := s.repo.dailySession.FindAll(ctx)
	s.NoError(err)
	s.Len(sessions, 1)

	// Delete
	err = s.repo.dailySession.Delete(ctx, session.ID)
	s.NoError(err)

	sessions, _ = s.repo.dailySession.FindAll(ctx)
	s.Len(sessions, 0)
}
