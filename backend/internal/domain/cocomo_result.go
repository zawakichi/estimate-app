package domain

// COCOMODetailedResult represents detailed COCOMO II estimation results
type COCOMODetailedResult struct {
    // Basic project information
    ProjectSize     float64 // KSLOC
    ModelType       string  // Early Design or Post-Architecture
    
    // Effort estimation
    BaseEffort      float64 // Person-months without adjustments
    AdjustedEffort  float64 // Person-months after applying all factors
    EffortRange     struct {
        Optimistic  float64 // -20% of nominal
        Nominal     float64 // Calculated effort
        Pessimistic float64 // +20% of nominal
    }
    
    // Schedule estimation
    Duration        float64 // Calendar months
    DurationRange   struct {
        Optimistic  float64
        Nominal     float64
        Pessimistic float64
    }
    
    // Team size estimation
    TeamSize        float64 // Average staff size
    TeamSizeRange   struct {
        Minimum     float64
        Average     float64
        Maximum     float64
    }
    
    // Cost estimation (if hourly rate is provided)
    CostEstimate    struct {
        HourlyRate  float64
        TotalCost   float64
        CostRange   struct {
            Minimum float64
            Nominal float64
            Maximum float64
        }
    }
    
    // Breakdown by phase (typical distribution for the selected process)
    PhaseDistribution []PhaseEffort
    
    // Factor analysis
    ScaleFactorAnalysis  []FactorAnalysis
    CostDriverAnalysis   []FactorAnalysis
    
    // Risk assessment
    RiskLevel       string  // Low, Medium, High
    RiskFactors     []RiskFactor
}

// PhaseEffort represents effort distribution for a development phase
type PhaseEffort struct {
    Phase           string  // Plans and Requirements, Product Design, Programming, etc.
    PercentEffort   float64 // Percentage of total effort
    Effort          float64 // Person-months for this phase
    Duration        float64 // Calendar months for this phase
    AverageStaff    float64 // Average staff size for this phase
}

// FactorAnalysis represents the impact analysis of a COCOMO II factor
type FactorAnalysis struct {
    Name            string
    Rating          float64 // Current rating value
    Impact          float64 // Multiplier or additive impact
    Sensitivity     float64 // How much the estimate changes with this factor
    Recommendation  string  // Optional recommendation for improvement
}

// RiskFactor represents a project risk identified through COCOMO II analysis
type RiskFactor struct {
    Category    string  // Technical, Cost, Schedule, or Process
    Name        string
    Level       string  // Low, Medium, High
    Impact      float64 // Estimated impact on effort/schedule
    Description string
    Mitigation  string  // Suggested mitigation strategy
}

// GenerateDetailedResult generates a detailed COCOMO II estimation result
func (e *COCOMOEstimate) GenerateDetailedResult(hourlyRate float64) *COCOMODetailedResult {
    result := &COCOMODetailedResult{
        ProjectSize: e.ProjectSize,
        ModelType:   e.Model.Name,
    }
    
    // Calculate base and adjusted effort
    result.BaseEffort = e.Model.A * pow(e.ProjectSize, e.Model.B)
    result.AdjustedEffort = e.EffortPM
    
    // Calculate effort range
    result.EffortRange.Nominal = e.EffortPM
    result.EffortRange.Optimistic = e.EffortPM * 0.8  // -20%
    result.EffortRange.Pessimistic = e.EffortPM * 1.2 // +20%
    
    // Calculate duration and range
    result.Duration = e.DurationTM
    result.DurationRange.Nominal = e.DurationTM
    result.DurationRange.Optimistic = e.DurationTM * 0.85  // -15%
    result.DurationRange.Pessimistic = e.DurationTM * 1.15 // +15%
    
    // Calculate team size ranges
    result.TeamSize = e.TeamSize
    result.TeamSizeRange.Average = e.TeamSize
    result.TeamSizeRange.Minimum = e.TeamSize * 0.7  // -30%
    result.TeamSizeRange.Maximum = e.TeamSize * 1.3  // +30%
    
    // Calculate cost if hourly rate is provided
    if hourlyRate > 0 {
        monthlyHours := 160.0 // Assuming 160 working hours per month
        totalCost := e.EffortPM * monthlyHours * hourlyRate
        
        result.CostEstimate.HourlyRate = hourlyRate
        result.CostEstimate.TotalCost = totalCost
        result.CostEstimate.CostRange.Nominal = totalCost
        result.CostEstimate.CostRange.Minimum = totalCost * 0.8  // -20%
        result.CostEstimate.CostRange.Maximum = totalCost * 1.2  // +20%
    }
    
    // Calculate phase distribution (typical distribution for software projects)
    result.PhaseDistribution = []PhaseEffort{
        {
            Phase:         "要件定義・計画",
            PercentEffort: 0.08,
            Effort:        e.EffortPM * 0.08,
            Duration:      e.DurationTM * 0.15,
            AverageStaff:  (e.EffortPM * 0.08) / (e.DurationTM * 0.15),
        },
        {
            Phase:         "システム設計",
            PercentEffort: 0.18,
            Effort:        e.EffortPM * 0.18,
            Duration:      e.DurationTM * 0.25,
            AverageStaff:  (e.EffortPM * 0.18) / (e.DurationTM * 0.25),
        },
        {
            Phase:         "詳細設計",
            PercentEffort: 0.25,
            Effort:        e.EffortPM * 0.25,
            Duration:      e.DurationTM * 0.35,
            AverageStaff:  (e.EffortPM * 0.25) / (e.DurationTM * 0.35),
        },
        {
            Phase:         "実装・単体テスト",
            PercentEffort: 0.26,
            Effort:        e.EffortPM * 0.26,
            Duration:      e.DurationTM * 0.45,
            AverageStaff:  (e.EffortPM * 0.26) / (e.DurationTM * 0.45),
        },
        {
            Phase:         "結合テスト",
            PercentEffort: 0.15,
            Effort:        e.EffortPM * 0.15,
            Duration:      e.DurationTM * 0.25,
            AverageStaff:  (e.EffortPM * 0.15) / (e.DurationTM * 0.25),
        },
        {
            Phase:         "システムテスト",
            PercentEffort: 0.08,
            Effort:        e.EffortPM * 0.08,
            Duration:      e.DurationTM * 0.15,
            AverageStaff:  (e.EffortPM * 0.08) / (e.DurationTM * 0.15),
        },
    }
    
    // Analyze scale factors
    for _, sf := range e.ScaleFactors {
        analysis := FactorAnalysis{
            Name:   sf.Name,
            Rating: sf.Rating,
            Impact: sf.Weight * sf.Rating,
        }
        
        // Calculate sensitivity
        sensitivity := (sf.Weight * 0.5) / e.EffortPM // Impact of 0.5 rating change
        analysis.Sensitivity = sensitivity
        
        // Add recommendations based on rating
        if sf.Rating > 3.5 {
            analysis.Recommendation = "この要因の改善により工数を削減できる可能性があります"
        }
        
        result.ScaleFactorAnalysis = append(result.ScaleFactorAnalysis, analysis)
    }
    
    // Analyze cost drivers
    for _, cd := range e.CostDrivers {
        analysis := FactorAnalysis{
            Name:   cd.Name,
            Rating: cd.Rating,
            Impact: cd.Value,
        }
        
        // Calculate sensitivity
        baseValue := cd.Value
        increasedValue := baseValue * 1.1 // 10% increase
        sensitivity := (increasedValue - baseValue) / baseValue
        analysis.Sensitivity = sensitivity
        
        // Add recommendations based on rating and impact
        if cd.Value > 1.2 {
            analysis.Recommendation = "この要因の最適化により工数を削減できる可能性があります"
        }
        
        result.CostDriverAnalysis = append(result.CostDriverAnalysis, analysis)
    }
    
    // Assess overall project risk
    result.RiskLevel = e.assessRiskLevel()
    result.RiskFactors = e.identifyRiskFactors()
    
    return result
}

// assessRiskLevel determines the overall project risk level
func (e *COCOMOEstimate) assessRiskLevel() string {
    // Count high-rated scale factors and cost drivers
    highRiskCount := 0
    
    for _, sf := range e.ScaleFactors {
        if sf.Rating > 4.0 {
            highRiskCount++
        }
    }
    
    for _, cd := range e.CostDrivers {
        if cd.Value > 1.3 {
            highRiskCount++
        }
    }
    
    if highRiskCount >= 3 {
        return "High"
    } else if highRiskCount >= 1 {
        return "Medium"
    }
    return "Low"
}

// identifyRiskFactors identifies specific project risk factors
func (e *COCOMOEstimate) identifyRiskFactors() []RiskFactor {
    var risks []RiskFactor
    
    // Analyze scale factors for risks
    for _, sf := range e.ScaleFactors {
        if sf.Rating > 4.0 {
            risk := RiskFactor{
                Category:    "Process",
                Name:        sf.Name,
                Level:      "High",
                Impact:     sf.Weight * sf.Rating,
                Description: "高いスケールファクター値による影響",
                Mitigation: "プロセスの改善とリスク軽減策の実施を検討",
            }
            risks = append(risks, risk)
        }
    }
    
    // Analyze cost drivers for risks
    for _, cd := range e.CostDrivers {
        if cd.Value > 1.3 {
            risk := RiskFactor{
                Category:    "Technical",
                Name:        cd.Name,
                Level:      "High",
                Impact:     cd.Value,
                Description: "高いコストドライバー値による影響",
                Mitigation: "技術的な対策と改善策の実施を検討",
            }
            risks = append(risks, risk)
        }
    }
    
    // Add size-related risks
    if e.ProjectSize > 100 { // Large project
        risks = append(risks, RiskFactor{
            Category:    "Technical",
            Name:        "大規模プロジェクト",
            Level:      "Medium",
            Impact:     1.3,
            Description: "プロジェクト規模が大きいことによる複雑性の増加",
            Mitigation: "モジュール化とインクリメンタル開発の採用を検討",
        })
    }
    
    return risks
}