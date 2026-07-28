package core_http_response

import (
	"fmt"
	"net/http"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
)

func SetCookie(rw http.ResponseWriter, name, value string, maxAge int, sameSite http.SameSite) error {
	cfg, err := core_config.NewConfig()
	if err != nil {
		return fmt.Errorf("get core config: %w", err)
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		SameSite: sameSite,
	}

	if cfg.EnvMode == core_config.EnvModeProd {
		cookie.Secure = true
	}

	http.SetCookie(rw, cookie)

	return nil
}

func DeleteCookie(rw http.ResponseWriter, name string) error {
	return SetCookie(rw, name, "", -1, http.SameSiteNoneMode)
}
