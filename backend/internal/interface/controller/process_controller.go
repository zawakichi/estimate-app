package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "estimate-backend/internal/usecase"
    "estimate-backend/internal/domain"
)

// ProcessController handles HTTP requests for process management
type ProcessController struct {
    processUseCase *usecase.ProcessUseCase
}

// NewProcessController creates a new ProcessController
func NewProcessController(pu *usecase.ProcessUseCase) *ProcessController {
    return &ProcessController{
        processUseCase: pu,
    }
}

// RegisterRoutes registers the routes for process management
func (pc *ProcessController) RegisterRoutes(e *echo.Echo) {
    e.GET("/api/processes", pc.GetAllProcesses)
    e.GET("/api/processes/:id", pc.GetProcess)
    e.PUT("/api/processes/:id", pc.UpdateProcess)
    e.PUT("/api/processes/:id/activities/:activityId", pc.UpdateActivity)
}

// GetAllProcesses handles GET /api/processes
func (pc *ProcessController) GetAllProcesses(c echo.Context) error {
    processes, err := pc.processUseCase.GetAllProcesses()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, processes)
}

// GetProcess handles GET /api/processes/:id
func (pc *ProcessController) GetProcess(c echo.Context) error {
    id := c.Param("id")
    process, err := pc.processUseCase.GetProcess(id)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Process not found")
    }
    return c.JSON(http.StatusOK, process)
}

// UpdateProcessRequest represents the request body for updating a process
type UpdateProcessRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Activities  []domain.Activity `json:"activities"`
}

// UpdateProcess handles PUT /api/processes/:id
func (pc *ProcessController) UpdateProcess(c echo.Context) error {
    id := c.Param("id")
    var req UpdateProcessRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    process := &domain.Process{
        ID:          id,
        Name:        req.Name,
        Description: req.Description,
        Activities:  req.Activities,
    }

    if err := pc.processUseCase.UpdateProcess(process); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, process)
}

// UpdateActivity handles PUT /api/processes/:id/activities/:activityId
func (pc *ProcessController) UpdateActivity(c echo.Context) error {
    processID := c.Param("id")
    activityID := c.Param("activityId")

    var activity domain.Activity
    if err := c.Bind(&activity); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    activity.ID = activityID
    if err := pc.processUseCase.UpdateActivity(processID, activity); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, activity)
}