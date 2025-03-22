package controller

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
    "estimate-backend/internal/usecase"
    "estimate-backend/internal/domain"
)

// EstimateController handles HTTP requests for estimate management
type EstimateController struct {
    estimateUseCase *usecase.EstimateUseCase
}

// NewEstimateController creates a new EstimateController
func NewEstimateController(eu *usecase.EstimateUseCase) *EstimateController {
    return &EstimateController{
        estimateUseCase: eu,
    }
}

// RegisterRoutes registers the routes for estimate management
func (ec *EstimateController) RegisterRoutes(e *echo.Echo) {
    e.POST("/api/estimates", ec.CreateEstimate)
    e.GET("/api/estimates/:id", ec.GetEstimate)
    e.PUT("/api/estimates/:id", ec.UpdateEstimate)
    e.GET("/api/estimates/:id/detailed", ec.GetDetailedEstimate)
    e.GET("/api/projects/:projectId/estimates", ec.GetProjectEstimates)
    e.POST("/api/estimates/compare", ec.CompareEstimates)
}

// CreateEstimateRequest represents the request body for creating an estimate
type CreateEstimateRequest struct {
    ProjectID     string                `json:"projectId"`
    ProjectName   string                `json:"projectName"`
    Tasks         []usecase.TaskInput   `json:"tasks"`
    GlobalFactors []string              `json:"globalFactors"`
    COCOMOData    *usecase.COCOMOInput  `json:"cocomoData,omitempty"`
    CreatedBy     string                `json:"createdBy"`
    Notes         string                `json:"notes"`
}

// CreateEstimate handles POST /api/estimates
func (ec *EstimateController) CreateEstimate(c echo.Context) error {
    var req CreateEstimateRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    input := usecase.CreateEstimateInput{
        ProjectID:     req.ProjectID,
        ProjectName:   req.ProjectName,
        Tasks:         req.Tasks,
        GlobalFactors: req.GlobalFactors,
        COCOMOData:    req.COCOMOData,
        CreatedBy:     req.CreatedBy,
        Notes:         req.Notes,
    }

    estimate, err := ec.estimateUseCase.CreateEstimate(input)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusCreated, estimate)
}

// GetEstimate handles GET /api/estimates/:id
func (ec *EstimateController) GetEstimate(c echo.Context) error {
    id := c.Param("id")
    estimate, err := ec.estimateUseCase.GetEstimate(id)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Estimate not found")
    }
    return c.JSON(http.StatusOK, estimate)
}

// UpdateEstimateRequest represents the request body for updating an estimate
type UpdateEstimateRequest struct {
    Tasks         []usecase.TaskInput   `json:"tasks"`
    GlobalFactors []string              `json:"globalFactors"`
    COCOMOData    *usecase.COCOMOInput  `json:"cocomoData,omitempty"`
    Notes         string                `json:"notes"`
}

// UpdateEstimate handles PUT /api/estimates/:id
func (ec *EstimateController) UpdateEstimate(c echo.Context) error {
    id := c.Param("id")
    var req UpdateEstimateRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    input := usecase.UpdateEstimateInput{
        ID:            id,
        Tasks:         req.Tasks,
        GlobalFactors: req.GlobalFactors,
        COCOMOData:    req.COCOMOData,
        Notes:         req.Notes,
    }

    estimate, err := ec.estimateUseCase.UpdateEstimate(input)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, estimate)
}

// GetDetailedEstimate handles GET /api/estimates/:id/detailed
func (ec *EstimateController) GetDetailedEstimate(c echo.Context) error {
    id := c.Param("id")
    hourlyRate, _ := strconv.ParseFloat(c.QueryParam("hourlyRate"), 64)

    estimate, cocomoResult, err := ec.estimateUseCase.GetDetailedEstimateResult(id, hourlyRate)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Estimate not found")
    }

    response := struct {
        *domain.Estimate
        COCOMODetails *domain.COCOMODetailedResult `json:"cocomoDetails,omitempty"`
    }{
        Estimate:      estimate,
        COCOMODetails: cocomoResult,
    }

    return c.JSON(http.StatusOK, response)
}

// GetProjectEstimates handles GET /api/projects/:projectId/estimates
func (ec *EstimateController) GetProjectEstimates(c echo.Context) error {
    projectID := c.Param("projectId")
    estimates, err := ec.estimateUseCase.GetProjectEstimates(projectID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, estimates)
}

// CompareEstimatesRequest represents the request body for comparing estimates
type CompareEstimatesRequest struct {
    EstimateID1 string `json:"estimateId1"`
    EstimateID2 string `json:"estimateId2"`
}

// CompareEstimates handles POST /api/estimates/compare
func (ec *EstimateController) CompareEstimates(c echo.Context) error {
    var req CompareEstimatesRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    comparison, err := ec.estimateUseCase.CompareEstimates(req.EstimateID1, req.EstimateID2)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, comparison)
}