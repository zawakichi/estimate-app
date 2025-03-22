package domain

// FactorType represents different types of factors that can affect estimation
type FactorType string

const (
    FactorTypeTeamExperience    FactorType = "team_experience"
    FactorTypeProjectComplexity FactorType = "project_complexity"
    FactorTypeTechnicalDebt     FactorType = "technical_debt"
    FactorTypeRiskBuffer        FactorType = "risk_buffer"
)

// Factor represents a multiplier that affects the estimation
type Factor struct {
    ID          string
    Type        FactorType
    Name        string
    Description string
    Impact      float64 // Multiplier value: 1.0 means no impact, > 1.0 increases time, < 1.0 decreases time
}

// Apply applies the factor to the given hours
func (f *Factor) Apply(hours float64) float64 {
    return hours * f.Impact
}

// FactorRepository defines the interface for factor persistence
type FactorRepository interface {
    Save(factor *Factor) error
    FindByID(id string) (*Factor, error)
    FindAll() ([]*Factor, error)
    Update(factor *Factor) error
    Delete(id string) error
}