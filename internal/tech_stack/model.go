package tech_stack

import (
	"time"

	"github.com/google/uuid"
)

type TechStack struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required"`
	Category    string     `json:"category" db:"category"`
	Version     string     `json:"version" db:"version"`
	Role        string     `json:"role" db:"role"`
	IsCoreSkill bool       `json:"is_core_skill" db:"is_core_skill"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type TechStackCreate struct {
	Name        string `json:"name" validate:"required"`
	Category    string `json:"category"`
	Version     string `json:"version"`
	Role        string `json:"role"`
	IsCoreSkill bool   `json:"is_core_skill"`
}

type TechStackUpdate struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Version     string    `json:"version"`
	Role        string    `json:"role"`
	IsCoreSkill bool      `json:"is_core_skill"`
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
