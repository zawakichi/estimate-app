package domain

import "time"

// EstimateStatus represents the status of an estimate
type EstimateStatus string

const (
    EstimateStatusDraft     EstimateStatus = "draft"
    EstimateStatusCompleted EstimateStatus = "completed"
    EstimateStatusApproved  EstimateStatus = "approved"
)

// ProcessEstimate represents estimation details for a specific process
type ProcessEstimate struct {
    Process     *Process
    Tasks       []Task
    BaseHours   float64
    TotalHours  float64  // After applying factors
}

// Estimate represents a work effort estimation for the entire project
type Estimate struct {
    ID              string
    ProjectID       string
    ProjectName     string
    ProcessEstimates []ProcessEstimate
    GlobalFactors   []Factor        // Factors that apply to the entire project
    COCOMOEstimate  *COCOMOEstimate // COCOMO II based estimation
    TotalHours      float64
    Status          EstimateStatus
    CreatedBy       string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Notes           string
}

// CalculationMethod represents the method used for effort calculation
type CalculationMethod string

const (
    CalculationMethodActivity CalculationMethod = "activity_based"
    CalculationMethodCOCOMO  CalculationMethod = "cocomo_based"
)

// CalculationResult represents the result of effort calculation
type CalculationResult struct {
    Method          CalculationMethod
    TotalHours      float64
    PersonMonths    float64
    TeamSize        float64
    DurationMonths  float64
    Confidence      float64  // 0-1, representing estimation confidence
}

// CalculateTotalHours calculates the total estimated hours using both activity-based and COCOMO II methods
func (e *Estimate) CalculateTotalHours(processRepo ProcessRepository) error {
    // Calculate activity-based estimation
    activityResult, err := e.calculateActivityBased(processRepo)
    if err != nil {
        return err
    }

    // Calculate COCOMO II based estimation if available
    var cocomoResult *CalculationResult
    if e.COCOMOEstimate != nil {
        cocomoResult = e.calculateCOCOMOBased()
    }

    // Combine and reconcile estimates
    e.reconcileEstimates(activityResult, cocomoResult)

    return nil
}

// calculateActivityBased performs the traditional activity-based calculation
func (e *Estimate) calculateActivityBased(processRepo ProcessRepository) (*CalculationResult, error) {
    var projectTotal float64

    // Calculate hours for each process
    for i, pe := range e.ProcessEstimates {
        process, err := processRepo.FindByID(pe.Process.ID)
        if err != nil {
            return nil, err
        }

        var processTotal float64
        // Calculate base hours for each task in the process
        for _, task := range pe.Tasks {
            // Find the corresponding activity
            var activity Activity
            for _, a := range process.Activities {
                if a.ID == task.ActivityID {
                    activity = a
                    break
                }
            }
            
            baseHours := task.CalculateBaseHours(activity)
            
            // Apply task-specific factors
            for _, factor := range task.CustomFactors {
                baseHours = factor.Apply(baseHours)
            }
            
            processTotal += baseHours
        }

        // Store the base hours before applying global factors
        e.ProcessEstimates[i].BaseHours = processTotal
        
        // Apply global factors to the process total
        for _, factor := range e.GlobalFactors {
            processTotal = factor.Apply(processTotal)
        }
        
        e.ProcessEstimates[i].TotalHours = processTotal
        projectTotal += processTotal
    }

    return &CalculationResult{
        Method:         CalculationMethodActivity,
        TotalHours:    projectTotal,
        PersonMonths:   projectTotal / 160.0, // Assuming 160 working hours per month
        TeamSize:       5.0,                  // Default team size, should be adjusted based on project scale
        DurationMonths: (projectTotal / 160.0) / 5.0,
        Confidence:     0.8,                  // Default confidence level for activity-based estimation
    }, nil
}

// calculateCOCOMOBased performs the COCOMO II based calculation
func (e *Estimate) calculateCOCOMOBased() *CalculationResult {
    // Recalculate COCOMO II estimate
    e.COCOMOEstimate.CalculateEffort()

    return &CalculationResult{
        Method:         CalculationMethodCOCOMO,
        TotalHours:    e.COCOMOEstimate.EffortPM * 160.0, // Convert person-months to hours
        PersonMonths:   e.COCOMOEstimate.EffortPM,
        TeamSize:       e.COCOMOEstimate.TeamSize,
        DurationMonths: e.COCOMOEstimate.DurationTM,
        Confidence:     0.85, // Default confidence level for COCOMO II estimation
    }
}

// reconcileEstimates combines activity-based and COCOMO II estimates
func (e *Estimate) reconcileEstimates(activityResult, cocomoResult *CalculationResult) {
    if cocomoResult == nil {
        // Use only activity-based estimation
        e.TotalHours = activityResult.TotalHours
        return
    }

    // Calculate weighted average based on confidence levels
    totalConfidence := activityResult.Confidence + cocomoResult.Confidence
    activityWeight := activityResult.Confidence / totalConfidence
    cocomoWeight := cocomoResult.Confidence / totalConfidence

    // Combine estimates
    e.TotalHours = (activityResult.TotalHours * activityWeight) +
                   (cocomoResult.TotalHours * cocomoWeight)
}

// EstimateRepository defines the interface for estimate persistence
type EstimateRepository interface {
    Save(estimate *Estimate) error
    FindByID(id string) (*Estimate, error)
    FindByProjectID(projectID string) ([]*Estimate, error)
    Update(estimate *Estimate) error
    Delete(id string) error
}