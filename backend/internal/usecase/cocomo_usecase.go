package usecase

import (
    "errors"
    "estimate-backend/internal/domain"
)

// COCOMOUseCase handles the business logic for COCOMO II estimations
type COCOMOUseCase struct {
    cocomoRepo domain.COCOMORepository
}

// NewCOCOMOUseCase creates a new COCOMOUseCase
func NewCOCOMOUseCase(cocomoRepo domain.COCOMORepository) *COCOMOUseCase {
    return &COCOMOUseCase{
        cocomoRepo: cocomoRepo,
    }
}

// InitializeDefaultModel initializes the default COCOMO II model
func (uc *COCOMOUseCase) InitializeDefaultModel() error {
    // Initialize Early Design model
    earlyDesign := &domain.COCOMOModel{
        Name:        "Early Design",
        Description: "COCOMO II Early Design model for early project estimation",
        A:           2.94,  // Calibrated value for Early Design
        B:           0.91,  // Initial exponent
    }

    // Initialize Post-Architecture model
    postArchitecture := &domain.COCOMOModel{
        Name:        "Post-Architecture",
        Description: "COCOMO II Post-Architecture model for detailed estimation",
        A:           2.45,  // Calibrated value for Post-Architecture
        B:           0.91,  // Initial exponent
    }

    if err := uc.cocomoRepo.SaveModel(earlyDesign); err != nil {
        return err
    }
    if err := uc.cocomoRepo.SaveModel(postArchitecture); err != nil {
        return err
    }

    return nil
}

// InitializeScaleFactors initializes the default scale factors
func (uc *COCOMOUseCase) InitializeScaleFactors() error {
    scaleFactors := []domain.ScaleFactor{
        {
            Type:        domain.ScaleFactorPREC,
            Name:        "先例性",
            Description: "類似プロジェクトの経験度",
            Weight:      4.05,
        },
        {
            Type:        domain.ScaleFactorFLEX,
            Name:        "開発の柔軟性",
            Description: "開発プロセスの柔軟性",
            Weight:      3.04,
        },
        {
            Type:        domain.ScaleFactorRESL,
            Name:        "アーキテクチャ/リスク対応",
            Description: "リスク管理とアーキテクチャ対応の程度",
            Weight:      4.24,
        },
        {
            Type:        domain.ScaleFactorTEAM,
            Name:        "チーム凝集性",
            Description: "チームの協力度と一貫性",
            Weight:      3.29,
        },
        {
            Type:        domain.ScaleFactorPMAT,
            Name:        "プロセス成熟度",
            Description: "組織のプロセス成熟度",
            Weight:      4.68,
        },
    }

    for _, sf := range scaleFactors {
        if err := uc.cocomoRepo.SaveScaleFactor(&sf); err != nil {
            return err
        }
    }

    return nil
}

// InitializeCostDrivers initializes the default cost drivers
func (uc *COCOMOUseCase) InitializeCostDrivers() error {
    costDrivers := []domain.CostDriver{
        // Product Factors
        {
            Type:        domain.CostDriverRELY,
            Name:        "要求される信頼性",
            Description: "システム障害による影響の大きさ",
            Value:       1.0, // Nominal value
        },
        {
            Type:        domain.CostDriverDATA,
            Name:        "データベース規模",
            Description: "テストデータベースサイズ/プログラムサイズの比",
            Value:       1.0,
        },
        {
            Type:        domain.CostDriverCPLX,
            Name:        "製品の複雑さ",
            Description: "制御操作、演算処理、デバイス処理、データ管理、UI管理の複雑さ",
            Value:       1.0,
        },
        // Platform Factors
        {
            Type:        domain.CostDriverTIME,
            Name:        "実行時間制約",
            Description: "使用可能な実行時間の制約",
            Value:       1.0,
        },
        {
            Type:        domain.CostDriverSTOR,
            Name:        "主記憶制約",
            Description: "主記憶の制約",
            Value:       1.0,
        },
        // Personnel Factors
        {
            Type:        domain.CostDriverACAP,
            Name:        "アナリスト能力",
            Description: "分析担当者の能力と経験",
            Value:       1.0,
        },
        {
            Type:        domain.CostDriverPCAP,
            Name:        "プログラマ能力",
            Description: "プログラマの能力と経験",
            Value:       1.0,
        },
        {
            Type:        domain.CostDriverPCON,
            Name:        "要員の継続性",
            Description: "プロジェクト期間中の要員の交代率",
            Value:       1.0,
        },
        // Project Factors
        {
            Type:        domain.CostDriverTOOL,
            Name:        "ツール使用",
            Description: "使用するツールの成熟度と機能",
            Value:       1.0,
        },
        {
            Type:        domain.CostDriverSITE,
            Name:        "開発拠点の分散",
            Description: "開発チームの地理的分散と通信手段",
            Value:       1.0,
        },
    }

    for _, cd := range costDrivers {
        if err := uc.cocomoRepo.SaveCostDriver(&cd); err != nil {
            return err
        }
    }

    return nil
}

// CreateEstimateInput represents input for creating a COCOMO II estimate
type CreateEstimateInput struct {
    ModelID       string
    ProjectSize   float64              // KSLOC or Function Points
    ScaleFactors map[string]float64    // Factor ID -> Rating
    CostDrivers  map[string]float64    // Driver ID -> Rating
}

// CreateEstimate creates a new COCOMO II estimate
func (uc *COCOMOUseCase) CreateEstimate(input CreateEstimateInput) (*domain.COCOMOEstimate, error) {
    // Validate input
    if input.ProjectSize <= 0 {
        return nil, errors.New("project size must be greater than 0")
    }

    // Get model
    model, err := uc.cocomoRepo.FindModelByID(input.ModelID)
    if err != nil {
        return nil, err
    }

    // Process scale factors
    var scaleFactors []domain.ScaleFactor
    for id, rating := range input.ScaleFactors {
        sf, err := uc.cocomoRepo.FindScaleFactorByID(id)
        if err != nil {
            return nil, err
        }
        sf.Rating = rating
        scaleFactors = append(scaleFactors, *sf)
    }

    // Process cost drivers
    var costDrivers []domain.CostDriver
    for id, rating := range input.CostDrivers {
        cd, err := uc.cocomoRepo.FindCostDriverByID(id)
        if err != nil {
            return nil, err
        }
        cd.Rating = rating
        costDrivers = append(costDrivers, *cd)
    }

    // Create estimate
    estimate := &domain.COCOMOEstimate{
        ProjectSize:  input.ProjectSize,
        Model:        model,
        ScaleFactors: scaleFactors,
        CostDrivers:  costDrivers,
    }

    // Calculate effort and other metrics
    estimate.CalculateEffort()

    // Save estimate
    if err := uc.cocomoRepo.SaveEstimate(estimate); err != nil {
        return nil, err
    }

    return estimate, nil
}

// GetEstimate retrieves a COCOMO II estimate by ID
func (uc *COCOMOUseCase) GetEstimate(id string) (*domain.COCOMOEstimate, error) {
    return uc.cocomoRepo.FindEstimateByID(id)
}

// UpdateRatingsInput represents input for updating scale factor and cost driver ratings
type UpdateRatingsInput struct {
    EstimateID    string
    ScaleFactors  map[string]float64    // Factor ID -> Rating
    CostDrivers   map[string]float64    // Driver ID -> Rating
}

// UpdateRatings updates the ratings of scale factors and cost drivers
func (uc *COCOMOUseCase) UpdateRatings(input UpdateRatingsInput) (*domain.COCOMOEstimate, error) {
    estimate, err := uc.cocomoRepo.FindEstimateByID(input.EstimateID)
    if err != nil {
        return nil, err
    }

    // Update scale factor ratings
    for id, rating := range input.ScaleFactors {
        for i, sf := range estimate.ScaleFactors {
            if sf.ID == id {
                estimate.ScaleFactors[i].Rating = rating
                break
            }
        }
    }

    // Update cost driver ratings
    for id, rating := range input.CostDrivers {
        for i, cd := range estimate.CostDrivers {
            if cd.ID == id {
                estimate.CostDrivers[i].Rating = rating
                break
            }
        }
    }

    // Recalculate effort and other metrics
    estimate.CalculateEffort()

    // Save updated estimate
    if err := uc.cocomoRepo.SaveEstimate(estimate); err != nil {
        return nil, err
    }

    return estimate, nil
}