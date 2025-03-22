package domain

import "time"

// Task represents a development task that needs to be estimated
type Task struct {
    ID            string
    ProcessID     string           // Reference to the Process this task belongs to
    ActivityID    string           // Reference to the specific Activity within the Process
    Name          string
    Description   string
    Complexity    int             // 1-5 scale
    Scale         float64         // Size/scale multiplier for the base hours
    Dependencies  []string        // IDs of dependent tasks
    CustomFactors []Factor        // Task-specific factors
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// CalculateBaseHours calculates the base hours for this task
func (t *Task) CalculateBaseHours(activity Activity) float64 {
    // Base calculation using activity's standard hours and task's scale
    baseHours := activity.BaseHours * t.Scale
    
    // Adjust based on complexity (1-5 scale)
    // Complexity 3 is considered normal (multiplier 1.0)
    complexityMultiplier := 0.8 + (float64(t.Complexity) * 0.2) // Results in range 1.0 +/- 40%
    
    return baseHours * complexityMultiplier
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
    Save(task *Task) error
    FindByID(id string) (*Task, error)
    FindByProcessID(processID string) ([]*Task, error)
    FindAll() ([]*Task, error)
    Update(task *Task) error
    Delete(id string) error
}