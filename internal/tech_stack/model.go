package tech_stack

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// TechStackCategory represents the predefined categories for tech stack
type TechStackCategory string

// Enum values for TechStackCategory matching the PostgreSQL enum type
const (
	CategoryBackend        TechStackCategory = "Backend"
	CategoryFrontend       TechStackCategory = "Frontend"
	CategoryFrameworks     TechStackCategory = "Frameworks"
	CategoryVersionControl TechStackCategory = "Version Control"
	CategoryDatabase       TechStackCategory = "Database"
	CategoryDevOps         TechStackCategory = "DevOps"
	CategoryTools          TechStackCategory = "Tools"
	CategoryCMSPlatforms   TechStackCategory = "CMS & Platforms"
)

// TechStack represents a technology stack entry
// @Description Technology stack information with details about skills and technologies
// @Name TechStack
type TechStack struct {
	ID          uuid.UUID         `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string            `json:"name" db:"name" validate:"required" example:"Go"`
	Category    TechStackCategory `json:"category" db:"category" example:"Backend"`
	Version     string            `json:"version" db:"version" example:"1.20"`
	Role        string            `json:"role" db:"role" example:"Backend Development"`
	IsCoreSkill bool              `json:"is_core_skill" db:"is_core_skill" example:"true"`
	ImageUrl    string            `json:"image_url" db:"image_url" example:"https://example.com/go-logo.png"`
	CreatedAt   *time.Time        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty" db:"updated_at"`
}

// TechStackCreate represents the input for creating a new tech stack
// @Name TechStackCreate
type TechStackCreate struct {
	Name        string                `json:"name" validate:"required" example:"Python"`
	Category    TechStackCategory     `json:"category" example:"Backend"`
	Version     string                `json:"version" example:"3.9"`
	Role        string                `json:"role" example:"Data Science"`
	IsCoreSkill bool                  `json:"is_core_skill" example:"true"`
	Image       *multipart.FileHeader `json:"image" swaggerignore:"true"`
}

// TechStackUpdate represents the input for updating an existing tech stack
// @Name TechStackUpdate
type TechStackUpdate struct {
	ID          uuid.UUID             `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string                `json:"name" example:"Rust"`
	Category    TechStackCategory     `json:"category" example:"Backend"`
	Version     string                `json:"version" example:"1.65"`
	Role        string                `json:"role" example:"Systems Programming"`
	IsCoreSkill bool                  `json:"is_core_skill" example:"true"`
	Image       *multipart.FileHeader `json:"image" swaggerignore:"true"`
}

// ToTechStack converts TechStackCreate to TechStack
func (tc *TechStackCreate) ToTechStack() TechStack {
	now := time.Now().UTC()
	return TechStack{
		ID:          uuid.New(),
		Name:        tc.Name,
		Category:    tc.Category,
		Version:     tc.Version,
		Role:        tc.Role,
		IsCoreSkill: tc.IsCoreSkill,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
}

// ToTechStack converts TechStackUpdate to TechStack
func (tu *TechStackUpdate) ToTechStack() TechStack {
	now := time.Now().UTC()
	return TechStack{
		ID:          tu.ID,
		Name:        tu.Name,
		Category:    tu.Category,
		Version:     tu.Version,
		Role:        tu.Role,
		IsCoreSkill: tu.IsCoreSkill,
		UpdatedAt:   &now,
	}
}
