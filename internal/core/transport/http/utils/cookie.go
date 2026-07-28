package core_http_utils

import (
	"fmt"
	"net/http"
	"time"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

func ParseCookie(r *http.Request, name string) (*http.Cookie, error) {
	cfg, err := core_config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("get core config: %w", err)
	}

	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("cookie `%s` is empty: %v: %w", name, err, core_errors.ErrInvalidArgument)
	}

	if cfg.EnvMode == core_config.EnvModeDev {
		return cookie, nil
	}

	if cfg.EnvMode == core_config.EnvModeProd && !cookie.Secure {
		return nil, fmt.Errorf("cookie `%s` must be Secure='true': %w", name, core_errors.ErrInvalidArgument)
	}

	if cookie.Expires.Before(time.Now()) || cookie.MaxAge < 0 {
		return nil, fmt.Errorf("cookie `%s` is expired: %w", name, core_errors.ErrInvalidArgument)
	}

	return cookie, nil
}
