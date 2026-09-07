package core_http_request

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type UserAgent struct {
	AppName    string
	AppVersion string

	Browser        string
	BrowserVersion string

	OS        string
	OSVersion string

	Model string
	Type  string
}

const (
	UserAgentTypePhone   = "phone"
	UserAgentTypeTablet  = "tablet"
	UserAgentTypeDesktop = "desktop"
	UserAgentTypeUnknown = "unknown"

	UserAgentBrowserChrome  = "chrome"
	UserAgentBrowserFirefox = "firefox"
	UserAgentBrowserSafari  = "safari"
	UserAgentBrowserEdge    = "edge"
	UserAgentBrowserOpera   = "opera"
	UserAgentBrowserOther   = "other"

	UserAgentOSWindows = "windows"
	UserAgentOSMacOS   = "macos"
	UserAgentOSLinux   = "linux"
	UserAgentOSAndroid = "android"
	UserAgentOSIOS     = "ios"
	UserAgentOSOther   = "other"

	UserAgentModelUnknown = "unknown"
)

var (
	// Cascade Pro:
	//
	// Cascade Pro/0.0.0 (Android 34; Xiaomi 11T; 1)
	//
	// Groups:
	// 1 - app name
	// 2 - app version
	// 3 - OS
	// 4 - model
	// 5 - device type
	cascadeProUserAgentRegex = regexp.MustCompile(
		`^([^/]+)/([^(]+)\s*\(([^;]+);\s*([^;]+);\s*([^)]+)\)$`,
	)

	chromeRegex = regexp.MustCompile(
		`(?:Chrome|CriOS)/([0-9.]+)`,
	)

	firefoxRegex = regexp.MustCompile(
		`(?:Firefox|FxiOS)/([0-9.]+)`,
	)

	edgeRegex = regexp.MustCompile(
		`(?:Edg|EdgiOS|EdgA)/([0-9.]+)`,
	)

	operaRegex = regexp.MustCompile(
		`(?:OPR|Opera)/([0-9.]+)`,
	)

	safariRegex = regexp.MustCompile(
		`Version/([0-9.]+).*Safari/`,
	)

	windowsRegex = regexp.MustCompile(
		`Windows NT ([0-9.]+)`,
	)

	macOSRegex = regexp.MustCompile(
		`Mac OS X ([0-9_\.]+)`,
	)

	androidRegex = regexp.MustCompile(
		`Android ([0-9\.]+)`,
	)

	iosRegex = regexp.MustCompile(
		`(?:iPhone OS|CPU OS) ([0-9_]+)`,
	)
)

// ParseUserAgent parses the HTTP User-Agent header.
//
// Supported formats:
//
// 1. Cascade Pro custom User-Agent:
//
//	Cascade Pro/0.0.0 (Android 34; Xiaomi 11T; 1)
//
// 2. Standard browser User-Agent:
//
//	Mozilla/5.0 (Linux; Android 15; Pixel 8) ...
//	Mozilla/5.0 (X11; Linux x86_64) ...
//	Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...
//	Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) ...
//	Mozilla/5.0 (iPhone; CPU iPhone OS 18_6 like Mac OS X) ...
func ParseUserAgent(r *http.Request) (*UserAgent, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil: %w", core_errors.ErrInvalidArgument)
	}

	ua := strings.TrimSpace(r.UserAgent())

	if ua == "" {
		return nil, fmt.Errorf("user-agent header is empty: %w", core_errors.ErrInvalidArgument)
	}

	// Keep compatibility with the mobile/Expo client.
	if userAgent, ok := parseCascadeProUserAgent(ua); ok {
		return userAgent, nil
	}

	return parseBrowserUserAgent(ua)
}

func parseCascadeProUserAgent(ua string) (*UserAgent, bool) {
	matches := cascadeProUserAgentRegex.FindStringSubmatch(ua)

	if matches == nil {
		return nil, false
	}

	deviceType := parseDeviceType(matches[5])

	return &UserAgent{
		AppName:    strings.TrimSpace(matches[1]),
		AppVersion: strings.TrimSpace(matches[2]),
		OS:         normalizeOSName(matches[3]),
		OSVersion:  extractOSVersion(strings.TrimSpace(matches[3])),
		Model:      strings.TrimSpace(matches[4]),
		Type:       deviceType,

		Browser:        "",
		BrowserVersion: "",
	}, true
}

func parseBrowserUserAgent(ua string) (*UserAgent, error) {
	userAgent := &UserAgent{
		AppName:    "",
		AppVersion: "",

		Browser:        UserAgentBrowserOther,
		BrowserVersion: "",

		OS:        UserAgentOSOther,
		OSVersion: "",

		Model: UserAgentModelUnknown,
		Type:  UserAgentTypeUnknown,
	}

	parseBrowser(userAgent, ua)
	parseOS(userAgent, ua)
	parseDevice(userAgent, ua)

	return userAgent, nil
}

func parseBrowser(userAgent *UserAgent, ua string) {
	// Edge must be checked before Chrome because Edge's UA
	// also contains "Chrome".
	if matches := edgeRegex.FindStringSubmatch(ua); matches != nil {
		userAgent.Browser = UserAgentBrowserEdge
		userAgent.BrowserVersion = matches[1]
		return
	}

	// Opera also contains Chrome in its UA.
	if matches := operaRegex.FindStringSubmatch(ua); matches != nil {
		userAgent.Browser = UserAgentBrowserOpera
		userAgent.BrowserVersion = matches[1]
		return
	}

	// Firefox.
	if matches := firefoxRegex.FindStringSubmatch(ua); matches != nil {
		userAgent.Browser = UserAgentBrowserFirefox
		userAgent.BrowserVersion = matches[1]
		return
	}

	// Chrome.
	if matches := chromeRegex.FindStringSubmatch(ua); matches != nil {
		userAgent.Browser = UserAgentBrowserChrome
		userAgent.BrowserVersion = matches[1]
		return
	}

	// Safari.
	//
	// Safari's UA contains Safari/ but Chrome-based browsers
	// also contain Safari/. Therefore Safari is checked last.
	if matches := safariRegex.FindStringSubmatch(ua); matches != nil {
		userAgent.Browser = UserAgentBrowserSafari
		userAgent.BrowserVersion = matches[1]
		return
	}
}

func parseOS(userAgent *UserAgent, ua string) {
	switch {
	case strings.Contains(ua, "Windows NT"):
		userAgent.OS = UserAgentOSWindows

		if matches := windowsRegex.FindStringSubmatch(ua); matches != nil {
			userAgent.OSVersion = normalizeVersion(matches[1])
		}

	case strings.Contains(ua, "Android"):
		userAgent.OS = UserAgentOSAndroid

		if matches := androidRegex.FindStringSubmatch(ua); matches != nil {
			userAgent.OSVersion = matches[1]
		}

	case strings.Contains(ua, "iPhone OS"),
		strings.Contains(ua, "iPad; CPU OS"),
		strings.Contains(ua, "iPod; CPU iPhone OS"):

		userAgent.OS = UserAgentOSIOS

		if matches := iosRegex.FindStringSubmatch(ua); matches != nil {
			userAgent.OSVersion = normalizeVersion(matches[1])
		}

	case strings.Contains(ua, "Mac OS X"):
		userAgent.OS = UserAgentOSMacOS

		if matches := macOSRegex.FindStringSubmatch(ua); matches != nil {
			userAgent.OSVersion = normalizeVersion(matches[1])
		}

	case strings.Contains(ua, "Linux"):
		userAgent.OS = UserAgentOSLinux
	}
}

func parseDevice(userAgent *UserAgent, ua string) {
	switch {
	case strings.Contains(ua, "iPhone"):
		userAgent.Type = UserAgentTypePhone
		userAgent.Model = "iPhone"

	case strings.Contains(ua, "iPad"):
		userAgent.Type = UserAgentTypeTablet
		userAgent.Model = "iPad"

	case strings.Contains(ua, "iPod"):
		userAgent.Type = UserAgentTypePhone
		userAgent.Model = "iPod"

	case strings.Contains(ua, "Android"):
		userAgent.Type = parseAndroidDeviceType(ua)
		userAgent.Model = parseAndroidModel(ua)

	case strings.Contains(ua, "Windows"),
		strings.Contains(ua, "Macintosh"),
		strings.Contains(ua, "X11"),
		strings.Contains(ua, "Linux"):
		userAgent.Type = UserAgentTypeDesktop
		userAgent.Model = UserAgentModelUnknown
	}
}

func parseAndroidDeviceType(ua string) string {
	// Android tablets commonly don't contain "Mobile".
	if strings.Contains(ua, "Mobile") {
		return UserAgentTypePhone
	}

	return UserAgentTypeTablet
}

func parseAndroidModel(ua string) string {
	// Android browser UA example:
	//
	// Mozilla/5.0 (Linux; Android 15; Pixel 8 Build/...)
	//
	// We extract the token between "Android <version>;" and
	// the next "Build/" marker.

	const androidPrefix = "Android "

	start := strings.Index(ua, androidPrefix)
	if start == -1 {
		return UserAgentModelUnknown
	}

	start = strings.Index(ua[start:], ";")
	if start == -1 {
		return UserAgentModelUnknown
	}

	start += strings.Index(ua, androidPrefix) + 1

	end := strings.Index(ua[start:], "Build/")
	if end == -1 {
		return UserAgentModelUnknown
	}

	model := strings.TrimSpace(ua[start : start+end])

	model = strings.TrimSuffix(model, ";")
	model = strings.TrimSpace(model)

	if model == "" || model == "wv" {
		return UserAgentModelUnknown
	}

	return model
}

func parseDeviceType(value string) string {
	switch strings.TrimSpace(value) {
	case "1":
		return UserAgentTypePhone
	case "2":
		return UserAgentTypeTablet
	case "3":
		return UserAgentTypeDesktop
	default:
		return UserAgentTypeUnknown
	}
}

func normalizeOSName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))

	switch {
	case strings.HasPrefix(value, "android"):
		return UserAgentOSAndroid

	case strings.HasPrefix(value, "ios"):
		return UserAgentOSIOS

	case strings.HasPrefix(value, "iphone"):
		return UserAgentOSIOS

	case strings.HasPrefix(value, "windows"):
		return UserAgentOSWindows

	case strings.HasPrefix(value, "mac"):
		return UserAgentOSMacOS

	case strings.HasPrefix(value, "linux"):
		return UserAgentOSLinux

	default:
		return UserAgentOSOther
	}
}

func extractOSVersion(value string) string {
	parts := strings.Fields(value)

	if len(parts) < 2 {
		return ""
	}

	return normalizeVersion(parts[1])
}

func normalizeVersion(version string) string {
	return strings.ReplaceAll(version, "_", ".")
}
