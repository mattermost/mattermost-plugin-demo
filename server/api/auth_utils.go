package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (api *Handlers) RequireAuthentication(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get(headerMattermostUserID)
	if userID == "" {
		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return errors.New("Unauthenticated user")
	}

	return nil
}

func (api *Handlers) DisallowGuestUsers(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get(headerMattermostUserID)
	user, appErr := api.pluginAPI.GetUser(userID)
	if appErr != nil {
		api.pluginAPI.LogError("DisallowGuestUsers: failed to get user from plugin API", "error", appErr.Error())
		return errors.New(appErr.Error())
	}

	if isGuest := user.IsGuest(); isGuest {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return errors.New("Guest users are not permitted to access this resource")
	}

	return nil
}

func (api *Handlers) RequireSystemAdmin(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get(headerMattermostUserID)
	user, appErr := api.pluginAPI.GetUser(userID)
	if appErr != nil {
		api.pluginAPI.LogError("RequireSystemAdmin: failed to get user from plugin API", "error", appErr.Error())
		return errors.New(appErr.Error())
	}

	if isSysAdmin := user.IsSystemAdmin(); !isSysAdmin {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return errors.New("Only system admins are permitted to access this resource")
	}

	return nil
}
