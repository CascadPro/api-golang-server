package domain

import (
	"fmt"
	"slices"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type PreferencesTheme string

type PreferencesAccent string

type PreferencesAnimation string

type PreferencesDensity string

type PreferencesTextSize string

const (
	ThemeSystem = PreferencesTheme("system")
	ThemeLight  = PreferencesTheme("light")
	ThemeDark   = PreferencesTheme("dark")

	AccentOrange = PreferencesAccent("orange")
	AccentBlue   = PreferencesAccent("blue")
	AccentMauve  = PreferencesAccent("mauve")
	AccentGreen  = PreferencesAccent("green")
	AccentTeal   = PreferencesAccent("teal")
	AccentRose   = PreferencesAccent("rose")

	AnimationsSystem   = PreferencesAnimation("system")
	AnimationsEnabled  = PreferencesAnimation("enabled")
	AnimationsDisabled = PreferencesAnimation("disabled")

	DensityComfortable = PreferencesDensity("comfortable")
	DensityCompact     = PreferencesDensity("compact")

	TextSizeSmall  = PreferencesTextSize("small")
	TextSizeMedium = PreferencesTextSize("medium")
	TextSizeLarge  = PreferencesTextSize("large")
)

var (
	PreferencesThemes = []PreferencesTheme{ThemeSystem, ThemeLight, ThemeDark}

	PreferencesAccents = []PreferencesAccent{AccentOrange, AccentBlue, AccentMauve, AccentGreen, AccentTeal, AccentRose}

	PreferencesAnimations = []PreferencesAnimation{AnimationsSystem, AnimationsEnabled, AnimationsDisabled}

	PreferencesDensities = []PreferencesDensity{DensityComfortable, DensityCompact}

	PreferencesTextSizes = []PreferencesTextSize{TextSizeSmall, TextSizeMedium, TextSizeLarge}
)

type UserPreferences struct {
	Theme                   PreferencesTheme
	Accent                  PreferencesAccent
	Animations              PreferencesAnimation
	Density                 PreferencesDensity
	TextSize                PreferencesTextSize
	ConfirmImportantActions bool
}

func DefaultUserPreferences() UserPreferences {
	return UserPreferences{
		Theme: "system", Accent: "orange",
		Animations: "system", Density: "comfortable",
		TextSize: "medium", ConfirmImportantActions: true,
	}
}

func (p *UserPreferences) Validate() error {
	if !slices.Contains(PreferencesThemes, p.Theme) {
		return fmt.Errorf("invalid user `Theme` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if !slices.Contains(PreferencesAccents, p.Accent) {
		return fmt.Errorf("invalid user `Accent` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if !slices.Contains(PreferencesAnimations, p.Animations) {
		return fmt.Errorf("invalid user `Animations` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if !slices.Contains(PreferencesDensities, p.Density) {
		return fmt.Errorf("invalid user `Density` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if !slices.Contains(PreferencesTextSizes, p.TextSize) {
		return fmt.Errorf("invalid user `TextSize` preferences: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type UserPreferencesPatch struct {
	Theme                   *PreferencesTheme
	Accent                  *PreferencesAccent
	Animations              *PreferencesAnimation
	Density                 *PreferencesDensity
	TextSize                *PreferencesTextSize
	ConfirmImportantActions *bool
}

func (p *UserPreferencesPatch) Validate() error {
	if p.Theme != nil && !slices.Contains(PreferencesThemes, *p.Theme) {
		return fmt.Errorf("invalid user `Theme` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if p.Accent != nil && !slices.Contains(PreferencesAccents, *p.Accent) {
		return fmt.Errorf("invalid user `Accent` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if p.Animations != nil && !slices.Contains(PreferencesAnimations, *p.Animations) {
		return fmt.Errorf("invalid user `Animations` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if p.Density != nil && !slices.Contains(PreferencesDensities, *p.Density) {
		return fmt.Errorf("invalid user `Density` preferences: %w", core_errors.ErrInvalidArgument)
	}
	if p.TextSize != nil && !slices.Contains(PreferencesTextSizes, *p.TextSize) {
		return fmt.Errorf("invalid user `TextSize` preferences: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (r *UserPreferencesPatch) Changes() map[string]any {
	changes := make(map[string]any)

	if r.Theme != nil {
		changes["0"] = *r.Theme
	}
	if r.Accent != nil {
		changes["1"] = *r.Accent
	}
	if r.Animations != nil {
		changes["2"] = *r.Animations
	}
	if r.Density != nil {
		changes["3"] = *r.Density
	}
	if r.TextSize != nil {
		changes["4"] = *r.TextSize
	}
	if r.ConfirmImportantActions != nil {
		changes["5"] = *r.ConfirmImportantActions
	}

	return changes
}

type UserNotificationPreferences struct {
	Requests bool
	Sessions bool
	System   bool
	Push     bool
}

func DefaultUserNotificationPreferences() UserNotificationPreferences {
	return UserNotificationPreferences{
		Requests: true,
		Sessions: true,
		System:   true,
		Push:     false,
	}
}

type UserNotificationPreferencesPatch struct {
	Requests *bool
	Sessions *bool
	System   *bool
	Push     *bool
}

func (p *UserNotificationPreferencesPatch) Validate() error {
	return nil
}

func (p *UserNotificationPreferencesPatch) Changes() map[string]bool {
	changes := make(map[string]bool)

	if p.Requests != nil {
		changes["0"] = *p.Requests
	}
	if p.Sessions != nil {
		changes["1"] = *p.Sessions
	}
	if p.System != nil {
		changes["2"] = *p.System
	}
	if p.Push != nil {
		changes["3"] = *p.Push
	}

	return changes
}
