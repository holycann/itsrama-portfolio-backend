package gemini

import (
	"fmt"
	"strings"
)

// AIInteractionRules defines comprehensive guidelines for AI interactions
type AIInteractionRules struct {
	// Ethical Guidelines
	EthicalPrinciples []string

	// Content Restrictions
	ProhibitedTopics     []string
	SensitiveTopics      []string
	ContentFilteringMode ContentFilterMode

	// Language and Tone
	AllowedLanguages []string
	ToneGuidelines   []ToneGuideline

	// Response Formatting
	ResponseStructure ResponseStructureRule

	// Contextual Awareness
	ContextualAwareness ContextAwarenessRule
}

// ContentFilterMode defines the strictness of content filtering
type ContentFilterMode int

const (
	FilterModeStrict ContentFilterMode = iota
	FilterModeMedium
	FilterModePermissive
)

// ToneGuideline defines acceptable communication styles
type ToneGuideline struct {
	Name        string
	Description string
	Allowed     bool
}

// ResponseStructureRule defines how AI responses should be formatted
type ResponseStructureRule struct {
	MaxLength        int
	MinLength        int
	RequiredSections []string
	AllowedFormats   []string
}

// ContextAwarenessRule defines how the AI should handle context
type ContextAwarenessRule struct {
	RememberPreviousContext bool
	MaxContextLength        int
	ContextTypes            []string
}

// DefaultAIInteractionRules provides a comprehensive set of default rules
func DefaultAIInteractionRules() *AIInteractionRules {
	return &AIInteractionRules{
		// Ethical Guidelines
		EthicalPrinciples: []string{
			"Respect human dignity",
			"Promote cultural understanding",
			"Avoid harmful or discriminatory content",
			"Prioritize user safety",
			"Maintain transparency about AI nature",
		},

		// Content Restrictions
		ProhibitedTopics: []string{
			"Explicit violence",
			"Hate speech",
			"Extreme political ideologies",
			"Graphic sexual content",
			"Illegal activities",
			"Personal medical advice",
		},
		SensitiveTopics: []string{
			"Mental health",
			"Personal trauma",
			"Religious beliefs",
			"Political conflicts",
		},
		ContentFilteringMode: FilterModeStrict,

		// Language and Tone
		AllowedLanguages: []string{"id", "en"},
		ToneGuidelines: []ToneGuideline{
			{
				Name:        "Respectful",
				Description: "Always maintain a polite and considerate tone",
				Allowed:     true,
			},
			{
				Name:        "Empathetic",
				Description: "Show understanding and compassion",
				Allowed:     true,
			},
			{
				Name:        "Aggressive",
				Description: "Confrontational or hostile language",
				Allowed:     false,
			},
		},

		// Response Formatting
		ResponseStructure: ResponseStructureRule{
			MaxLength:        2000,
			MinLength:        10,
			RequiredSections: []string{"Introduction", "Main Content", "Conclusion"},
			AllowedFormats:   []string{"Plain Text", "Markdown", "Numbered List", "Bullet Points"},
		},

		// Contextual Awareness
		ContextualAwareness: ContextAwarenessRule{
			RememberPreviousContext: true,
			MaxContextLength:        5,
			ContextTypes:            []string{"Event Details", "User Preferences", "Previous Interactions"},
		},
	}
}

// GetFullSystemPolicy generates a comprehensive system policy for AI interactions
func GetFullSystemPolicy() string {
	rules := DefaultAIInteractionRules()

	var policyBuilder strings.Builder
	policyBuilder.WriteString("AI Interaction System Policy\n\n")

	// Ethical Principles
	policyBuilder.WriteString("1. Ethical Guidelines:\n")
	for _, principle := range rules.EthicalPrinciples {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", principle))
	}

	// Content Restrictions
	policyBuilder.WriteString("\n2. Content Restrictions:\n")
	policyBuilder.WriteString("   Prohibited Topics:\n")
	for _, topic := range rules.ProhibitedTopics {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", topic))
	}
	policyBuilder.WriteString("   Sensitive Topics (require extra care):\n")
	for _, topic := range rules.SensitiveTopics {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", topic))
	}
	policyBuilder.WriteString(fmt.Sprintf("   Content Filtering Mode: %s\n", getFilterModeName(rules.ContentFilteringMode)))

	// Language and Tone
	policyBuilder.WriteString("\n3. Language and Communication:\n")
	policyBuilder.WriteString("   Allowed Languages:\n")
	for _, lang := range rules.AllowedLanguages {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", lang))
	}
	policyBuilder.WriteString("   Tone Guidelines:\n")
	for _, guideline := range rules.ToneGuidelines {
		status := "Allowed"
		if !guideline.Allowed {
			status = "Not Allowed"
		}
		policyBuilder.WriteString(fmt.Sprintf("   - %s: %s (%s)\n", guideline.Name, guideline.Description, status))
	}

	// Response Formatting
	policyBuilder.WriteString("\n4. Response Formatting:\n")
	policyBuilder.WriteString(fmt.Sprintf("   - Max Length: %d characters\n", rules.ResponseStructure.MaxLength))
	policyBuilder.WriteString(fmt.Sprintf("   - Min Length: %d characters\n", rules.ResponseStructure.MinLength))
	policyBuilder.WriteString("   Required Sections:\n")
	for _, section := range rules.ResponseStructure.RequiredSections {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", section))
	}
	policyBuilder.WriteString("   Allowed Formats:\n")
	for _, format := range rules.ResponseStructure.AllowedFormats {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", format))
	}

	// Contextual Awareness
	policyBuilder.WriteString("\n5. Contextual Awareness:\n")
	contextStatus := "Enabled"
	if !rules.ContextualAwareness.RememberPreviousContext {
		contextStatus = "Disabled"
	}
	policyBuilder.WriteString(fmt.Sprintf("   - Context Retention: %s\n", contextStatus))
	policyBuilder.WriteString(fmt.Sprintf("   - Max Context Length: %d interactions\n", rules.ContextualAwareness.MaxContextLength))
	policyBuilder.WriteString("   Context Types:\n")
	for _, contextType := range rules.ContextualAwareness.ContextTypes {
		policyBuilder.WriteString(fmt.Sprintf("   - %s\n", contextType))
	}

	return policyBuilder.String()
}

// Helper function to get filter mode name
func getFilterModeName(mode ContentFilterMode) string {
	switch mode {
	case FilterModeStrict:
		return "Strict"
	case FilterModeMedium:
		return "Medium"
	case FilterModePermissive:
		return "Permissive"
	default:
		return "Unknown"
	}
}
