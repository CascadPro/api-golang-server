package settings_service

import "github.com/CascadePro/api-golang-server/internal/core/domain"

func changedSettingsFields(
	before domain.UserSettings,
	after domain.UserSettings,
) map[string]any {
	changes := map[string]any{
		"version": after.Version,
	}

	if before.SessionExpireTerm != after.SessionExpireTerm {
		changes["session_expire_term"] = after.SessionExpireTerm
	}

	if before.Session2FAEnabled != after.Session2FAEnabled {
		changes["session_2fa_enabled"] = after.Session2FAEnabled
	}

	preferences := make(map[string]any)

	if before.Preferences.Theme != after.Preferences.Theme {
		preferences["theme"] = after.Preferences.Theme
	}

	if before.Preferences.Accent != after.Preferences.Accent {
		preferences["accent"] = after.Preferences.Accent
	}

	if before.Preferences.Animations != after.Preferences.Animations {
		preferences["animations"] = after.Preferences.Animations
	}

	if before.Preferences.Density != after.Preferences.Density {
		preferences["density"] = after.Preferences.Density
	}

	if before.Preferences.TextSize != after.Preferences.TextSize {
		preferences["text_size"] = after.Preferences.TextSize
	}

	if before.Preferences.ConfirmImportantActions != after.Preferences.ConfirmImportantActions {
		preferences["confirm_important_actions"] = after.Preferences.ConfirmImportantActions
	}

	if len(preferences) > 0 {
		changes["preferences"] = preferences
	}

	notifications := make(map[string]any)

	if before.Notifications.Requests != after.Notifications.Requests {
		notifications["requests"] = after.Notifications.Requests
	}

	if before.Notifications.Sessions != after.Notifications.Sessions {
		notifications["sessions"] = after.Notifications.Sessions
	}

	if before.Notifications.System != after.Notifications.System {
		notifications["system"] = after.Notifications.System
	}

	if before.Notifications.Push != after.Notifications.Push {
		notifications["push"] = after.Notifications.Push
	}

	if len(notifications) > 0 {
		changes["notifications"] = notifications
	}

	return changes
}
