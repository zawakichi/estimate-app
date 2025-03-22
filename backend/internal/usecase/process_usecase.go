package usecase

import (
    "errors"
    "estimate-backend/internal/domain"
)

// ProcessUseCase handles the business logic for development processes
type ProcessUseCase struct {
    processRepo domain.ProcessRepository
}

// NewProcessUseCase creates a new ProcessUseCase
func NewProcessUseCase(processRepo domain.ProcessRepository) *ProcessUseCase {
    return &ProcessUseCase{
        processRepo: processRepo,
    }
}

// InitializeDefaultProcesses creates the default set of development processes
func (uc *ProcessUseCase) InitializeDefaultProcesses() error {
    defaultProcesses := []domain.Process{
        {
            Category:    domain.ProcessRequirementDefinition,
            Name:       "要件定義",
            Description: "プロジェクトの要件を定義し、スコープを決定する工程",
            Order:      1,
            Activities: []domain.Activity{
                {
                    Name:        "ステークホルダーヒアリング",
                    Description: "関係者からの要件収集",
                    BaseHours:   16,
                    Deliverables: []string{"ヒアリング議事録", "要件一覧"},
                },
                {
                    Name:        "要件分析",
                    Description: "収集した要件の分析と整理",
                    BaseHours:   24,
                    Deliverables: []string{"要件定義書"},
                },
                {
                    Name:        "スコープ定義",
                    Description: "プロジェクトスコープの定義と合意形成",
                    BaseHours:   16,
                    Deliverables: []string{"スコープ定義書", "除外事項一覧"},
                },
            },
        },
        {
            Category:    domain.ProcessFunctionalSpec,
            Name:       "機能仕様検討",
            Description: "システムの機能仕様を検討する工程",
            Order:      2,
            Activities: []domain.Activity{
                {
                    Name:        "機能一覧作成",
                    Description: "システムの機能一覧の作成",
                    BaseHours:   24,
                    Deliverables: []string{"機能一覧表"},
                },
                {
                    Name:        "画面設計",
                    Description: "ユーザーインターフェースの設計",
                    BaseHours:   40,
                    Deliverables: []string{"画面設計書", "画面遷移図"},
                },
                {
                    Name:        "機能仕様書作成",
                    Description: "詳細な機能仕様の定義",
                    BaseHours:   40,
                    Deliverables: []string{"機能仕様書"},
                },
            },
        },
        {
            Category:    domain.ProcessBasicDesign,
            Name:       "基本設計",
            Description: "システムの基本的なアーキテクチャを設計する工程",
            Order:      3,
            Activities: []domain.Activity{
                {
                    Name:        "アーキテクチャ設計",
                    Description: "システム全体のアーキテクチャ設計",
                    BaseHours:   40,
                    Deliverables: []string{"アーキテクチャ設計書"},
                },
                {
                    Name:        "データベース設計",
                    Description: "データベースの基本設計",
                    BaseHours:   32,
                    Deliverables: []string{"ER図", "テーブル定義書"},
                },
                {
                    Name:        "セキュリティ設計",
                    Description: "セキュリティ要件の設計",
                    BaseHours:   24,
                    Deliverables: []string{"セキュリティ設計書"},
                },
            },
        },
        {
            Category:    domain.ProcessDetailedDesign,
            Name:       "詳細設計",
            Description: "システムの詳細な設計を行う工程",
            Order:      4,
            Activities: []domain.Activity{
                {
                    Name:        "モジュール設計",
                    Description: "各モジュールの詳細設計",
                    BaseHours:   48,
                    Deliverables: []string{"モジュール設計書"},
                },
                {
                    Name:        "API設計",
                    Description: "APIインターフェースの設計",
                    BaseHours:   32,
                    Deliverables: []string{"API仕様書"},
                },
                {
                    Name:        "単体テスト設計",
                    Description: "単体テストの設計",
                    BaseHours:   24,
                    Deliverables: []string{"単体テスト仕様書"},
                },
            },
        },
        {
            Category:    domain.ProcessImplementation,
            Name:       "実装",
            Description: "システムの実装を行う工程",
            Order:      5,
            Activities: []domain.Activity{
                {
                    Name:        "フロントエンド実装",
                    Description: "フロントエンドの実装",
                    BaseHours:   80,
                    Deliverables: []string{"ソースコード", "単体テスト結果"},
                },
                {
                    Name:        "バックエンド実装",
                    Description: "バックエンドの実装",
                    BaseHours:   80,
                    Deliverables: []string{"ソースコード", "単体テスト結果"},
                },
                {
                    Name:        "データベース実装",
                    Description: "データベースの実装",
                    BaseHours:   24,
                    Deliverables: []string{"DDLスクリプト", "初期データ"},
                },
            },
        },
        {
            Category:    domain.ProcessTesting,
            Name:       "テスト",
            Description: "システムのテストを行う工程",
            Order:      6,
            Activities: []domain.Activity{
                {
                    Name:        "結合テスト",
                    Description: "モジュール間の結合テスト",
                    BaseHours:   40,
                    Deliverables: []string{"結合テスト結果報告書"},
                },
                {
                    Name:        "システムテスト",
                    Description: "システム全体のテスト",
                    BaseHours:   56,
                    Deliverables: []string{"システムテスト結果報告書"},
                },
                {
                    Name:        "性能テスト",
                    Description: "性能要件の検証",
                    BaseHours:   32,
                    Deliverables: []string{"性能テスト結果報告書"},
                },
            },
        },
        {
            Category:    domain.ProcessDelivery,
            Name:       "納品",
            Description: "システムの納品を行う工程",
            Order:      7,
            Activities: []domain.Activity{
                {
                    Name:        "マニュアル作成",
                    Description: "各種マニュアルの作成",
                    BaseHours:   40,
                    Deliverables: []string{"運用マニュアル", "利用者マニュアル"},
                },
                {
                    Name:        "導入支援",
                    Description: "システムの導入支援",
                    BaseHours:   24,
                    Deliverables: []string{"導入手順書", "導入報告書"},
                },
                {
                    Name:        "検収対応",
                    Description: "検収作業の対応",
                    BaseHours:   16,
                    Deliverables: []string{"検収報告書"},
                },
            },
        },
    }

    for _, process := range defaultProcesses {
        if err := uc.processRepo.Save(&process); err != nil {
            return err
        }
    }

    return nil
}

// GetProcess retrieves a process by ID
func (uc *ProcessUseCase) GetProcess(id string) (*domain.Process, error) {
    return uc.processRepo.FindByID(id)
}

// GetProcessByCategory retrieves a process by its category
func (uc *ProcessUseCase) GetProcessByCategory(category domain.ProcessCategory) (*domain.Process, error) {
    return uc.processRepo.FindByCategory(category)
}

// GetAllProcesses retrieves all processes in order
func (uc *ProcessUseCase) GetAllProcesses() ([]*domain.Process, error) {
    return uc.processRepo.FindAll()
}

// UpdateProcess updates an existing process
func (uc *ProcessUseCase) UpdateProcess(process *domain.Process) error {
    if process.ID == "" {
        return errors.New("process ID is required")
    }
    return uc.processRepo.Update(process)
}

// UpdateActivity updates an activity within a process
func (uc *ProcessUseCase) UpdateActivity(processID string, activity domain.Activity) error {
    process, err := uc.processRepo.FindByID(processID)
    if err != nil {
        return err
    }

    // Find and update the activity
    found := false
    for i, act := range process.Activities {
        if act.ID == activity.ID {
            process.Activities[i] = activity
            found = true
            break
        }
    }

    if !found {
        return errors.New("activity not found in process")
    }

    return uc.processRepo.Update(process)
}