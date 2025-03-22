package domain

// ProcessCategory represents the main development process categories
type ProcessCategory string

const (
    ProcessRequirementDefinition ProcessCategory = "requirement_definition"
    ProcessFunctionalSpec       ProcessCategory = "functional_specification"
    ProcessBasicDesign         ProcessCategory = "basic_design"
    ProcessDetailedDesign      ProcessCategory = "detailed_design"
    ProcessImplementation      ProcessCategory = "implementation"
    ProcessTesting            ProcessCategory = "testing"
    ProcessDelivery           ProcessCategory = "delivery"
)

// Process represents a development process category and its standard activities
type Process struct {
    ID          string
    Category    ProcessCategory
    Name        string
    Description string
    Activities  []Activity
    Order       int // For maintaining the natural order of processes
}

// Activity represents a standard activity within a process
type Activity struct {
    ID          string
    Name        string
    Description string
    BaseHours   float64    // Standard base hours for this activity
    Deliverables []string  // Expected deliverables from this activity
}

// ProcessRepository defines the interface for process persistence
type ProcessRepository interface {
    Save(process *Process) error
    FindByID(id string) (*Process, error)
    FindByCategory(category ProcessCategory) (*Process, error)
    FindAll() ([]*Process, error)
    Update(process *Process) error
    Delete(id string) error
}