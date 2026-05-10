package handler

import (
	"daily-app-go/db/sqlc"
	"daily-app-go/internal/models"
	"daily-app-go/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Register In TaskHandler Register method (around lines 25-36):
func (h *TaskHandler) Register(r gin.IRouter) {
	g := r.Group("/tasks")
	g.GET("", h.GetAllTasks)
	g.GET("/:taskId", h.GetTaskById)
	g.GET("/:taskId/details", h.GetTaskDetailsById)
	g.POST("", h.CreateTask)
	g.PATCH("/:taskId", h.UpdateTask)
	g.DELETE("/:taskId", h.DeleteTask)
	g.POST("/:taskId/start", h.StartTask)
	g.POST("/:taskId/complete", h.CompleteTask)
	g.POST("/:taskId/reopen", h.ReopenTask)
	g.POST("/:taskId/cancel", h.CancelTask)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.FindAll(c.Request.Context())
	if err != nil {
		mapError(c, err)
		return
	}
	if len(tasks) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	task, err := h.taskService.FindById(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	if task == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetTaskDetailsById(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	details, err := h.taskService.GetTaskDetailsById(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	if details == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Task details not found"})
		return
	}
	c.JSON(http.StatusOK, details)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.taskService.Save(c.Request.Context(), req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.taskService.Update(c.Request.Context(), id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	if err := h.taskService.DeleteById(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) StartTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	if err := h.taskService.StartTask(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *TaskHandler) CompleteTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	if err := h.taskService.CompleteTask(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *TaskHandler) ReopenTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	if err := h.taskService.ReopenTask(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *TaskHandler) CancelTask(c *gin.Context) {
	id, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	if id <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	if err := h.taskService.CancelTask(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

// SubTaskHandler handles subtask-related requests
type SubTaskHandler struct {
	subTaskService service.SubTaskService
}

func NewSubTaskHandler(subTaskService service.SubTaskService) *SubTaskHandler {
	return &SubTaskHandler{subTaskService: subTaskService}
}

func (h *SubTaskHandler) Register(r gin.IRouter) {
	g := r.Group("/tasks/:taskId/subtasks")
	g.GET("", h.GetSubTasksByTaskId)
	g.GET("/:subTaskId", h.GetSubTaskById)
	g.POST("", h.CreateSubTask)
	g.PATCH("/:subTaskId", h.UpdateSubTask)
	g.DELETE("/:subTaskId", h.DeleteSubTask)
	g.POST("/:subTaskId/complete", h.CompleteSubTask)
	g.POST("/:subTaskId/reopen", h.ReopenSubTask)
}

func (h *SubTaskHandler) GetSubTasksByTaskId(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	subTasks, err := h.subTaskService.FindAllByTaskId(c.Request.Context(), taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	if len(subTasks) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, subTasks)
}

func (h *SubTaskHandler) GetSubTaskById(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	id, ok := getIDParam(c, "subTaskId")
	if !ok {
		return
	}
	subTask, err := h.subTaskService.FindById(c.Request.Context(), id, taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	if subTask == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "SubTask not found"})
		return
	}
	c.JSON(http.StatusOK, subTask)
}

func (h *SubTaskHandler) CreateSubTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	var req models.CreateAndUpdateSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subTask, err := h.subTaskService.Save(c.Request.Context(), taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.Header("Location", fmt.Sprintf("/tasks/%d/subtasks/%d", taskId, subTask.ID))
	c.JSON(http.StatusCreated, subTask)
}

func (h *SubTaskHandler) UpdateSubTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	id, ok := getIDParam(c, "subTaskId")
	if !ok {
		return
	}
	var req models.CreateAndUpdateSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subTask, err := h.subTaskService.Update(c.Request.Context(), taskId, id, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, subTask)
}

func (h *SubTaskHandler) DeleteSubTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	id, ok := getIDParam(c, "subTaskId")
	if !ok {
		return
	}
	if err := h.subTaskService.DeleteById(c.Request.Context(), taskId, id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SubTaskHandler) CompleteSubTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	id, ok := getIDParam(c, "subTaskId")
	if !ok {
		return
	}
	if _, err := h.subTaskService.Complete(c.Request.Context(), taskId, id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *SubTaskHandler) ReopenSubTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	id, ok := getIDParam(c, "subTaskId")
	if !ok {
		return
	}
	if err := h.subTaskService.Reopen(c.Request.Context(), taskId, id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

// NoteHandler handles note-related requests
type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) Register(r gin.IRouter) {
	g := r.Group("/tasks/:taskId/notes")
	g.GET("", h.GetNotesForTask)
	g.GET("/:noteId", h.GetNote)
	g.POST("", h.CreateNote)
	g.PUT("/:noteId", h.ReplaceNote)
	g.PATCH("/:noteId", h.UpdateNote)
	g.DELETE("/:noteId", h.DeleteNote)
}

func (h *NoteHandler) GetNotesForTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	notes, err := h.noteService.FindByTaskId(c.Request.Context(), taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	if len(notes) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	noteId, ok := getIDParam(c, "noteId")
	if !ok {
		return
	}
	note, err := h.noteService.FindById(c.Request.Context(), noteId, taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	if note == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note, err := h.noteService.Save(c.Request.Context(), taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) ReplaceNote(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	noteId, ok := getIDParam(c, "noteId")
	if !ok {
		return
	}
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note, err := h.noteService.Replace(c.Request.Context(), noteId, taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	noteId, ok := getIDParam(c, "noteId")
	if !ok {
		return
	}
	var req models.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note, err := h.noteService.Update(c.Request.Context(), noteId, taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	noteId, ok := getIDParam(c, "noteId")
	if !ok {
		return
	}
	if err := h.noteService.DeleteById(c.Request.Context(), taskId, noteId); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// BlockerHandler handles blocker-related requests
type BlockerHandler struct {
	blockerService service.BlockerService
}

func NewBlockerHandler(blockerService service.BlockerService) *BlockerHandler {
	return &BlockerHandler{blockerService: blockerService}
}

func (h *BlockerHandler) Register(r gin.IRouter) {
	g := r.Group("/tasks/:taskId/blockers")
	g.GET("", h.GetBlockersForTask)
	g.GET("/:blockerId", h.GetBlocker)
	g.POST("", h.CreateBlocker)
	g.PUT("/:blockerId", h.ReplaceBlocker)
	g.PATCH("/:blockerId", h.UpdateBlocker)
	g.DELETE("/:blockerId", h.DeleteBlocker)
	g.POST("/:blockerId/resolve", h.ResolveBlocker)
	g.POST("/:blockerId/reopen", h.ReopenBlocker)
}

func (h *BlockerHandler) GetBlockersForTask(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockers, err := h.blockerService.FindByTaskId(c.Request.Context(), taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	if len(blockers) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, blockers)
}

func (h *BlockerHandler) GetBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	blocker, err := h.blockerService.FindById(c.Request.Context(), blockerId)
	if err != nil {
		mapError(c, err)
		return
	}
	if blocker == nil || blocker.TaskID != taskId {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Blocker not found"})
		return
	}
	c.JSON(http.StatusOK, blocker)
}

func (h *BlockerHandler) CreateBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	var req models.CreateBlockerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blocker, err := h.blockerService.Save(c.Request.Context(), taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, blocker)
}

func (h *BlockerHandler) ReplaceBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	var req models.CreateBlockerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blocker, err := h.blockerService.Replace(c.Request.Context(), blockerId, taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, blocker)
}

func (h *BlockerHandler) UpdateBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	var req models.UpdateBlockerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blocker, err := h.blockerService.Update(c.Request.Context(), blockerId, taskId, req)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, blocker)
}

func (h *BlockerHandler) DeleteBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	blocker, err := h.blockerService.FindById(c.Request.Context(), blockerId)
	if err != nil {
		mapError(c, err)
		return
	}
	if blocker == nil || blocker.TaskID != taskId {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Blocker not found"})
		return
	}
	if err := h.blockerService.DeleteById(c.Request.Context(), blockerId); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BlockerHandler) ResolveBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	blocker, err := h.blockerService.Resolve(c.Request.Context(), blockerId, taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, blocker)
}

func (h *BlockerHandler) ReopenBlocker(c *gin.Context) {
	taskId, ok := getIDParam(c, "taskId")
	if !ok {
		return
	}
	blockerId, ok := getIDParam(c, "blockerId")
	if !ok {
		return
	}
	blocker, err := h.blockerService.Reopen(c.Request.Context(), blockerId, taskId)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, blocker)
}

// DailySessionHandler handles daily session-related requests
type DailySessionHandler struct {
	dailySessionService service.DailySessionService
}

func NewDailySessionHandler(dailySessionService service.DailySessionService) *DailySessionHandler {
	return &DailySessionHandler{dailySessionService: dailySessionService}
}

func (h *DailySessionHandler) Register(r gin.IRouter) {
	g := r.Group("/daily-sessions")
	g.GET("", h.GetAllSessions)
	g.GET("/:id", h.GetSessionById)
	g.POST("", h.CreateSession)
	g.DELETE("/:id", h.DeleteSession)
}

func (h *DailySessionHandler) GetAllSessions(c *gin.Context) {
	sessions, err := h.dailySessionService.FindAll(c.Request.Context())
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *DailySessionHandler) GetSessionById(c *gin.Context) {
	id, ok := getIDParam(c, "id")
	if !ok {
		return
	}
	session, err := h.dailySessionService.FindById(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	if session == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Daily session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *DailySessionHandler) CreateSession(c *gin.Context) {
	var ds db.DailySession
	if err := c.ShouldBindJSON(&ds); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := h.dailySessionService.Save(c.Request.Context(), ds)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *DailySessionHandler) DeleteSession(c *gin.Context) {
	id, ok := getIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.dailySessionService.DeleteById(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Helpers

func getIDParam(c *gin.Context, name string) (int32, bool) {
	idStr := c.Param(name)
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid ID parameter: %s", name)})
		return 0, false
	}
	return int32(id), true
}

func mapError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	msg := err.Error()
	lowMsg := strings.ToLower(msg)

	// Not found errors
	if strings.Contains(lowMsg, "not found") || strings.Contains(lowMsg, "no rows in result set") {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	// Policy / Validation / Conflict errors
	if strings.Contains(lowMsg, "cannot") || strings.Contains(lowMsg, "must") || strings.Contains(lowMsg, "belong") {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": msg})
		return
	}

	// Default to Internal Server Error
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
}
