package usecase

import (
    "errors"
    "estimate-backend/internal/domain"
)

// FactorUseCase handles the business logic for estimation factors
type FactorUseCase struct {
    factorRepo domain.FactorRepository
}

// NewFactorUseCase creates a new FactorUseCase
func NewFactorUseCase(factorRepo domain.FactorRepository) *FactorUseCase {
    return &FactorUseCase{
        factorRepo: factorRepo,
    }
}

// InitializeDefaultFactors creates the default set of estimation factors
func (uc *FactorUseCase) InitializeDefaultFactors() error {
    defaultFactors := []domain.Factor{
        // チーム経験関連の要因
        {
            Type:        domain.FactorTypeTeamExperience,
            Name:        "新規技術スタック",
            Description: "チームが使用する技術スタックが新しい場合の影響",
            Impact:      1.5, // 50%増
        },
        {
            Type:        domain.FactorTypeTeamExperience,
            Name:        "ドメイン知識不足",
            Description: "チームが業務ドメインに不慣れな場合の影響",
            Impact:      1.3, // 30%増
        },
        {
            Type:        domain.FactorTypeTeamExperience,
            Name:        "熟練チーム",
            Description: "チームが技術とドメインの両方に精通している場合",
            Impact:      0.8, // 20%減
        },

        // プロジェクト複雑性関連の要因
        {
            Type:        domain.FactorTypeProjectComplexity,
            Name:        "システム間連携多数",
            Description: "多数の外部システムとの連携が必要な場合",
            Impact:      1.4, // 40%増
        },
        {
            Type:        domain.FactorTypeProjectComplexity,
            Name:        "セキュリティ要件厳格",
            Description: "特に厳格なセキュリティ要件がある場合",
            Impact:      1.3, // 30%増
        },
        {
            Type:        domain.FactorTypeProjectComplexity,
            Name:        "パフォーマンス要件厳格",
            Description: "特に厳格なパフォーマンス要件がある場合",
            Impact:      1.25, // 25%増
        },

        // 技術的負債関連の要因
        {
            Type:        domain.FactorTypeTechnicalDebt,
            Name:        "レガシーシステム改修",
            Description: "古いシステムの改修や統合が必要な場合",
            Impact:      1.5, // 50%増
        },
        {
            Type:        domain.FactorTypeTechnicalDebt,
            Name:        "ドキュメント不足",
            Description: "既存システムのドキュメントが不足している場合",
            Impact:      1.2, // 20%増
        },
        {
            Type:        domain.FactorTypeTechnicalDebt,
            Name:        "テスト自動化不足",
            Description: "テスト自動化が不十分な場合",
            Impact:      1.15, // 15%増
        },

        // リスクバッファー関連の要因
        {
            Type:        domain.FactorTypeRiskBuffer,
            Name:        "要件不確実性",
            Description: "要件の変更や追加が予想される場合",
            Impact:      1.3, // 30%増
        },
        {
            Type:        domain.FactorTypeRiskBuffer,
            Name:        "スケジュール圧縮",
            Description: "タイトなスケジュールでの開発が必要な場合",
            Impact:      1.25, // 25%増
        },
        {
            Type:        domain.FactorTypeRiskBuffer,
            Name:        "チーム規模大",
            Description: "大規模なチームでの開発による調整コスト",
            Impact:      1.2, // 20%増
        },
    }

    for _, factor := range defaultFactors {
        if err := uc.factorRepo.Save(&factor); err != nil {
            return err
        }
    }

    return nil
}

// CreateFactorInput represents input data for creating a factor
type CreateFactorInput struct {
    Type        domain.FactorType
    Name        string
    Description string
    Impact      float64
}

// CreateFactor creates a new estimation factor
func (uc *FactorUseCase) CreateFactor(input CreateFactorInput) (*domain.Factor, error) {
    // Validate input
    if input.Name == "" {
        return nil, errors.New("factor name is required")
    }
    if input.Impact <= 0 {
        return nil, errors.New("impact must be greater than 0")
    }

    factor := &domain.Factor{
        Type:        input.Type,
        Name:        input.Name,
        Description: input.Description,
        Impact:      input.Impact,
    }

    if err := uc.factorRepo.Save(factor); err != nil {
        return nil, err
    }

    return factor, nil
}

// UpdateFactorInput represents input data for updating a factor
type UpdateFactorInput struct {
    ID          string
    Type        domain.FactorType
    Name        string
    Description string
    Impact      float64
}

// UpdateFactor updates an existing factor
func (uc *FactorUseCase) UpdateFactor(input UpdateFactorInput) (*domain.Factor, error) {
    factor, err := uc.factorRepo.FindByID(input.ID)
    if err != nil {
        return nil, err
    }

    factor.Type = input.Type
    factor.Name = input.Name
    factor.Description = input.Description
    factor.Impact = input.Impact

    if err := uc.factorRepo.Update(factor); err != nil {
        return nil, err
    }

    return factor, nil
}

// GetFactor retrieves a factor by ID
func (uc *FactorUseCase) GetFactor(id string) (*domain.Factor, error) {
    return uc.factorRepo.FindByID(id)
}

// GetAllFactors retrieves all factors
func (uc *FactorUseCase) GetAllFactors() ([]*domain.Factor, error) {
    return uc.factorRepo.FindAll()
}

// DeleteFactor deletes a factor by ID
func (uc *FactorUseCase) DeleteFactor(id string) error {
    return uc.factorRepo.Delete(id)
}