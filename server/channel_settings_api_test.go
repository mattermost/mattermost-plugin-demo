package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testChannelID = "channelid1234567890123456a"

func newChannelSettingsPlugin(api *plugintest.API) *Plugin {
	p := &Plugin{}
	p.SetAPI(api)
	p.client = pluginapi.NewClient(api, nil)
	p.initializeAPI()
	return p
}

func TestChannelSettingsSchemaRoundTrip(t *testing.T) {
	api := &plugintest.API{}
	defer api.AssertExpectations(t)

	store := map[string][]byte{}
	api.On("KVSetWithOptions", "cs_schema_"+testChannelID, mock.Anything, mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
		store[args.String(0)] = args.Get(1).([]byte)
	})
	api.On("KVGet", "cs_schema_"+testChannelID).Return(func(key string) []byte { return store[key] }, nil)
	api.On("GetChannel", testChannelID).Return(&model.Channel{Id: testChannelID, Type: model.ChannelTypeOpen}, nil)
	api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionManagePublicChannelProperties).Return(true)
	api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionReadChannel).Return(true)

	p := newChannelSettingsPlugin(api)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/channel_settings/"+testChannelID+"/schema", strings.NewReader(`{"postPrefixStyle":"bold"}`))
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/channel_settings/"+testChannelID+"/schema", nil)
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.JSONEq(t, `{"postPrefixStyle":"bold"}`, w.Body.String())
}

func TestChannelSettingsCustomRoundTrip(t *testing.T) {
	api := &plugintest.API{}
	defer api.AssertExpectations(t)

	store := map[string][]byte{}
	api.On("KVSetWithOptions", "cs_custom_"+testChannelID, mock.Anything, mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
		store[args.String(0)] = args.Get(1).([]byte)
	})
	api.On("KVGet", "cs_custom_"+testChannelID).Return(func(key string) []byte { return store[key] }, nil)
	api.On("GetChannel", testChannelID).Return(&model.Channel{Id: testChannelID, Type: model.ChannelTypeOpen}, nil)
	api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionManagePublicChannelProperties).Return(true)
	api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionReadChannel).Return(true)

	p := newChannelSettingsPlugin(api)

	// GET before any save returns the zero-valued custom state.
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/channel_settings/"+testChannelID+"/custom", nil)
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.JSONEq(t, `{"note":"","pinGreeting":false}`, w.Body.String())

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/channel_settings/"+testChannelID+"/custom", strings.NewReader(`{"note":"hello","pinGreeting":true}`))
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/channel_settings/"+testChannelID+"/custom", nil)
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.JSONEq(t, `{"note":"hello","pinGreeting":true}`, w.Body.String())
}

func TestChannelSettingsReadAuthorization(t *testing.T) {
	t.Run("get without auth is unauthorized", func(t *testing.T) {
		api := &plugintest.API{}
		p := newChannelSettingsPlugin(api)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/channel_settings/"+testChannelID+"/schema", nil)
		p.ServeHTTP(nil, w, r)
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("get without read permission is forbidden", func(t *testing.T) {
		api := &plugintest.API{}
		defer api.AssertExpectations(t)
		api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionReadChannel).Return(false)

		p := newChannelSettingsPlugin(api)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/channel_settings/"+testChannelID+"/custom", nil)
		r.Header.Set("Mattermost-User-ID", "user1")
		p.ServeHTTP(nil, w, r)
		require.Equal(t, http.StatusForbidden, w.Result().StatusCode)
	})
}

func TestChannelSettingsMalformedBody(t *testing.T) {
	api := &plugintest.API{}
	defer api.AssertExpectations(t)
	api.On("GetChannel", testChannelID).Return(&model.Channel{Id: testChannelID, Type: model.ChannelTypeOpen}, nil)
	api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionManagePublicChannelProperties).Return(true)

	p := newChannelSettingsPlugin(api)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/channel_settings/"+testChannelID+"/custom", strings.NewReader(`{not json`))
	r.Header.Set("Mattermost-User-ID", "user1")
	p.ServeHTTP(nil, w, r)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestChannelSettingsAuthAndPermission(t *testing.T) {
	t.Run("save without auth is unauthorized", func(t *testing.T) {
		api := &plugintest.API{}
		p := newChannelSettingsPlugin(api)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/channel_settings/"+testChannelID+"/schema", strings.NewReader(`{}`))
		p.ServeHTTP(nil, w, r)
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("save without permission is forbidden", func(t *testing.T) {
		api := &plugintest.API{}
		defer api.AssertExpectations(t)
		api.On("GetChannel", testChannelID).Return(&model.Channel{Id: testChannelID, Type: model.ChannelTypePrivate}, nil)
		api.On("HasPermissionToChannel", "user1", testChannelID, model.PermissionManagePrivateChannelProperties).Return(false)

		p := newChannelSettingsPlugin(api)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/channel_settings/"+testChannelID+"/custom", strings.NewReader(`{"note":"x"}`))
		r.Header.Set("Mattermost-User-ID", "user1")
		p.ServeHTTP(nil, w, r)
		require.Equal(t, http.StatusForbidden, w.Result().StatusCode)
	})
}
