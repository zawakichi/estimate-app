package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "estimate-backend/internal/usecase"
    "estimate-backend/internal/domain"
)

// COCOMOController handles HTTP requests for COCOMO II related operations
type COCOMOController struct {
    cocomoUseCase *usecase.COCOMOUseCase
}

// NewCOCOMOController creates a new COCOMOController
func NewCOCOMOController(cu *usecase.COCOMOUseCase) *COCOMOController {
    return &COCOMOController{
        cocomoUseCase: cu,
    }
}

// RegisterRoutes registers the routes for COCOMO II management
func (cc *COCOMOController) RegisterRoutes(e *echo.Echo) {
    e.GET("/api/cocomo/models", cc.GetModels)
    e.GET("/api/cocomo/scale-factors", cc.GetScaleFactors)
    e.GET("/api/cocomo/cost-drivers", cc.GetCostDrivers)
    e.POST("/api/cocomo/calculate", cc.CalculateEstimate)
}

// GetModels handles GET /api/cocomo/models
func (cc *COCOMOController) GetModels(c echo.Context) error {
    // Initialize default models if not exists
    if err := cc.cocomoUseCase.InitializeDefaultModel(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    // Return the models (Early Design and Post-Architecture)
    return c.JSON(http.StatusOK, map[string]interface{}{
        "models": []string{"Early Design", "Post-Architecture"},
    })
}

// GetScaleFactors handles GET /api/cocomo/scale-factors
func (cc *COCOMOController) GetScaleFactors(c echo.Context) error {
    // Initialize default scale factors if not exists
    if err := cc.cocomoUseCase.InitializeScaleFactors(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    // Return the scale factors with their descriptions and weight ranges
    return c.JSON(http.StatusOK, map[string]interface{}{
        "scaleFactors": []map[string]interface{}{
            {
                "type": domain.ScaleFactorPREC,
                "name": "先例性",
                "description": "類似プロジェクトの経験度",
                "ratingGuide": map[string]string{
                    "very_low":    "全く新しい開発",
                    "low":         "大部分が新規",
                    "nominal":     "類似経験あり",
                    "high":        "ほぼ同様の開発経験あり",
                    "very_high":   "ほぼ同一の開発",
                },
            },
            {
                "type": domain.ScaleFactorFLEX,
                "name": "開発の柔軟性",
                "description": "開発プロセスの柔軟性",
                "ratingGuide": map[string]string{
                    "very_low":    "厳格な制約あり",
                    "low":         "一部柔軟性あり",
                    "nominal":     "ある程度柔軟",
                    "high":        "大部分が柔軟",
                    "very_high":   "完全に柔軟",
                },
            },
            // 他のスケールファクターも同様に定義
        },
    })
}

// GetCostDrivers handles GET /api/cocomo/cost-drivers
func (cc *COCOMOController) GetCostDrivers(c echo.Context) error {
    // Initialize default cost drivers if not exists
    if err := cc.cocomoUseCase.InitializeCostDrivers(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    // Return the cost drivers with their descriptions and rating guides
    return c.JSON(http.StatusOK, map[string]interface{}{
        "costDrivers": []map[string]interface{}{
            {
                "type": domain.CostDriverRELY,
                "name": "要求される信頼性",
                "description": "システム障害による影響の大きさ",
                "ratingGuide": map[string]string{
                    "very_low":    "軽微な不便",
                    "low":         "軽度の損失",
                    "nominal":     "中程度の損失",
                    "high":        "大きな損失",
                    "very_high":   "人命に関わる",
                },
            },
            {
                "type": domain.CostDriverCPLX,
                "name": "製品の複雑さ",
                "description": "制御操作、演算処理、デバイス処理、データ管理、UI管理の複雑さ",
                "ratingGuide": map[string]string{
                    "very_low":    "単純な処理",
                    "low":         "やや複雑",
                    "nominal":     "中程度",
                    "high":        "複雑",
                    "very_high":   "非常に複雑",
                },
            },
            // 他のコストドライバーも同様に定義
        },
    })
}

// CalculateEstimateRequest represents the request body for COCOMO II calculation
type CalculateEstimateRequest struct {
    ModelID       string             `json:"modelId"`
    KSLOC        float64            `json:"ksloc"`
    ScaleFactors map[string]float64 `json:"scaleFactors"`
    CostDrivers  map[string]float64 `json:"costDrivers"`
}

// CalculateEstimate handles POST /api/cocomo/calculate
func (cc *COCOMOController) CalculateEstimate(c echo.Context) error {
    var req CalculateEstimateRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    input := usecase.CreateEstimateInput{
        ModelID:      req.ModelID,
        ProjectSize:  req.KSLOC,
        ScaleFactors: req.ScaleFactors,
        CostDrivers:  req.CostDrivers,
    }

    estimate, err := cc.cocomoUseCase.CreateEstimate(input)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    // Generate detailed result with cost calculation
    detailedResult := estimate.GenerateDetailedResult(0) // hourlyRate = 0 for now

    return c.JSON(http.StatusOK, detailedResult)
}