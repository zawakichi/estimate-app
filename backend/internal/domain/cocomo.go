package domain

// COCOMOModel represents the COCOMO II estimation model configuration
type COCOMOModel struct {
    ID          string
    Name        string
    Description string
    // Base coefficients for effort equation: PM = A * Size^B * EM
    A           float64 // Multiplicative constant
    B           float64 // Scale factor
}

// ScaleFactorType represents different types of COCOMO II scale factors
type ScaleFactorType string

const (
    // COCOMO II Scale Factors
    ScaleFactorPREC ScaleFactorType = "precedentedness"      // 先例性
    ScaleFactorFLEX ScaleFactorType = "development_flexibility" // 開発の柔軟性
    ScaleFactorRESL ScaleFactorType = "architecture_risk"    // アーキテクチャ/リスク解決
    ScaleFactorTEAM ScaleFactorType = "team_cohesion"        // チーム凝集性
    ScaleFactorPMAT ScaleFactorType = "process_maturity"     // プロセス成熟度
)

// ScaleFactor represents a COCOMO II scale factor
type ScaleFactor struct {
    ID          string
    Type        ScaleFactorType
    Name        string
    Description string
    Rating      float64 // Very Low (0) to Extra High (5)
    Weight      float64 // Impact on the exponential scale factor
}

// CostDriverType represents different types of COCOMO II cost drivers
type CostDriverType string

const (
    // Product Factors
    CostDriverRELY CostDriverType = "required_reliability"    // 要求される信頼性
    CostDriverDATA CostDriverType = "database_size"          // データベース規模
    CostDriverCPLX CostDriverType = "product_complexity"     // 製品の複雑さ
    CostDriverREUS CostDriverType = "required_reusability"   // 要求される再利用性
    CostDriverDOCU CostDriverType = "documentation"          // ドキュメント化

    // Platform Factors
    CostDriverTIME CostDriverType = "execution_time"         // 実行時間制約
    CostDriverSTOR CostDriverType = "storage_constraint"     // 主記憶制約
    CostDriverPVOL CostDriverType = "platform_volatility"    // プラットフォーム揮発性

    // Personnel Factors
    CostDriverACAP CostDriverType = "analyst_capability"     // アナリスト能力
    CostDriverPCAP CostDriverType = "programmer_capability"  // プログラマ能力
    CostDriverPCON CostDriverType = "personnel_continuity"   // 要員の継続性
    CostDriverAPEX CostDriverType = "application_experience" // アプリケーション経験
    CostDriverPLEX CostDriverType = "platform_experience"    // プラットフォーム経験
    CostDriverLTEX CostDriverType = "language_experience"    // 言語・ツール経験

    // Project Factors
    CostDriverTOOL CostDriverType = "tool_use"              // ツール使用
    CostDriverSITE CostDriverType = "multisite_development" // 開発拠点の分散
    CostDriverSCED CostDriverType = "schedule_constraint"    // 要求される開発工期
)

// CostDriver represents a COCOMO II cost driver
type CostDriver struct {
    ID          string
    Type        CostDriverType
    Name        string
    Description string
    Rating      float64 // Very Low (0) to Extra High (5)
    Value       float64 // Effort multiplier value
}

// COCOMOEstimate represents a COCOMO II based estimation
type COCOMOEstimate struct {
    ID           string
    ProjectSize  float64       // Size in KSLOC or Function Points
    Model        *COCOMOModel
    ScaleFactors []ScaleFactor
    CostDrivers  []CostDriver
    // Calculated values
    ExponentB    float64  // Calculated from scale factors
    EffortPM     float64  // Person-Months
    DurationTM   float64  // Time-Months
    TeamSize     float64  // Average team size
}

// CalculateEffort calculates the effort in person-months using COCOMO II
func (e *COCOMOEstimate) CalculateEffort() {
    // Calculate the exponential scale factor (B)
    e.ExponentB = e.Model.B
    for _, sf := range e.ScaleFactors {
        e.ExponentB += sf.Weight * sf.Rating
    }

    // Calculate the effort multiplier (EM)
    em := 1.0
    for _, cd := range e.CostDrivers {
        em *= cd.Value
    }

    // Calculate effort: PM = A * Size^B * EM
    e.EffortPM = e.Model.A * pow(e.ProjectSize, e.ExponentB) * em

    // Calculate duration: TDEV = C * (PM)^D
    // where C and D are empirically derived constants
    c := 3.67
    d := 0.28 + 0.2 * (e.ExponentB - 1.01)
    e.DurationTM = c * pow(e.EffortPM, d)

    // Calculate average team size
    e.TeamSize = e.EffortPM / e.DurationTM
}

// Helper function for power calculation
func pow(base, exp float64) float64 {
    result := 1.0
    for i := 0; i < int(exp); i++ {
        result *= base
    }
    return result
}

// COCOMORepository defines the interface for COCOMO II model persistence
type COCOMORepository interface {
    SaveModel(model *COCOMOModel) error
    FindModelByID(id string) (*COCOMOModel, error)
    SaveEstimate(estimate *COCOMOEstimate) error
    FindEstimateByID(id string) (*COCOMOEstimate, error)
    SaveScaleFactor(factor *ScaleFactor) error
    FindScaleFactorByID(id string) (*ScaleFactor, error)
    SaveCostDriver(driver *CostDriver) error
    FindCostDriverByID(id string) (*CostDriver, error)
}