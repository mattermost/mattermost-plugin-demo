package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"unicode/utf8"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddOrderRow(t *testing.T) {
	state, err := addOrderRow(defaultOrderRows())
	require.NoError(t, err)
	assert.Len(t, state.Rows, len(defaultOrderRows().Rows)+1)
	assert.Equal(t, "r4", state.Rows[len(state.Rows)-1].ID, "the new row takes the next id")
	assert.Equal(t, 5, state.NextID, "NextID advances so cell names stay unique")
	assert.Empty(t, state.Rows[len(state.Rows)-1].Order, "the new row starts blank")

	for len(state.Rows) < maxOrderRows {
		state, err = addOrderRow(state)
		require.NoError(t, err)
	}
	_, err = addOrderRow(state)
	assert.Error(t, err, "rows are capped at maxOrderRows")
}

func TestFillableTableState(t *testing.T) {
	original := defaultOrderRows()

	state, err := decodeFillableTableState(encodeFillableTableState(original))
	require.NoError(t, err)
	assert.Equal(t, original, state)

	_, err = decodeFillableTableState("somestate")
	assert.Error(t, err, "state that is not JSON is rejected")

	_, err = decodeFillableTableState(`{"rows":[]}`)
	assert.Error(t, err, "state with no rows is rejected")

	_, err = decodeFillableTableState(`{"rows":[{"id":""}]}`)
	assert.Error(t, err, "a row with no id is rejected, because cell names are built from it")

	over := original
	over.Rows = make([]orderRow, maxOrderRows+1)
	for i := range over.Rows {
		over.Rows[i].ID = fmt.Sprintf("r%d", i)
	}
	_, err = decodeFillableTableState(encodeFillableTableState(over))
	assert.Error(t, err, "state over the row cap is rejected")

	// State round trips through the client, so an unrecognised width has to
	// fall back rather than fail the whole dialog.
	wide := original
	wide.Size = "enormous"
	state, err = decodeFillableTableState(encodeFillableTableState(wide))
	require.NoError(t, err)
	assert.Equal(t, model.DialogSizeSmall, state.Size)
}

func TestApplySubmissionToRows(t *testing.T) {
	state := applySubmissionToRows(defaultOrderRows(), map[string]any{
		cellName("r1", "order"):    "  SO-9999  ",
		cellName("r1", "weight"):   "12.5",
		cellName("r2", "status"):   "Shipped",
		cellName("r9", "order"):    "ignored",
		"some_unrelated_field":     "ignored",
		cellName("r3", "customer"): "",
	})

	assert.Equal(t, "SO-9999", state.Rows[0].Order, "values are trimmed")
	assert.Equal(t, "12.5", state.Rows[0].Weight)
	assert.Equal(t, "Acme Co", state.Rows[0].Customer, "a cell absent from the submission keeps its value")
	assert.Equal(t, "Shipped", state.Rows[1].Status)
	assert.Empty(t, state.Rows[2].Customer, "a submitted empty cell clears the value")
	assert.Len(t, state.Rows, 3, "a submission cannot add rows")
}

func TestApplyWidthToState(t *testing.T) {
	state := applyWidthToState(defaultOrderRows(), map[string]any{fillableWidthField: model.DialogSizeLarge})
	assert.Equal(t, model.DialogSizeLarge, state.Size)

	state = applyWidthToState(state, map[string]any{fillableWidthField: "enormous"})
	assert.Equal(t, model.DialogSizeSmall, state.Size, "an unrecognised width falls back rather than failing")

	state = applyWidthToState(defaultOrderRows(), map[string]any{})
	assert.Equal(t, defaultOrderRows().Size, state.Size, "an absent width leaves the state alone")
}

func TestValidateOrders(t *testing.T) {
	state := defaultOrderRows()

	// Everything required is still blank, so every row reports its own cells.
	errs := validateOrders(state)
	assert.Len(t, errs, 2*len(state.Rows), "a weight and a status error per row")
	assert.Contains(t, errs, cellName("r1", "weight"))
	assert.Contains(t, errs, cellName("r1", "status"))
	assert.NotContains(t, errs, cellName("r1", "order"), "the default rows already have order numbers")

	for i := range state.Rows {
		state.Rows[i].Weight = "10"
		state.Rows[i].Status = "New"
	}
	assert.Empty(t, validateOrders(state))

	// Each error addresses exactly one cell. That is the whole point of the
	// flatten approach: the flat Errors map is keyed by element name.
	state.Rows[1].Weight = "-5"
	state.Rows[1].Cases = "two"
	state.Rows[2].Order = ""

	errs = validateOrders(state)
	assert.Equal(t, map[string]string{
		cellName("r2", "weight"): "Enter a weight of 0 or more",
		cellName("r2", "cases"):  "Cases must be a whole number",
		cellName("r3", "order"):  "An order number is required",
	}, errs)

	// Cases is optional, so blank is fine but nonsense is not.
	state.Rows[1].Weight = "0"
	state.Rows[1].Cases = ""
	state.Rows[2].Order = "SO-1051"
	assert.Empty(t, validateOrders(state))
}

func TestRowAction(t *testing.T) {
	action := rowActionFromContext(map[string]any{
		"row_id":   "r1",
		"order":    " SO-1042 ",
		"customer": "Acme Co",
	})
	assert.Equal(t, rowAction{RowID: "r1", Order: "SO-1042", Customer: "Acme Co"}, action)

	// interfaceToString renders a nil as "<nil>", which would satisfy the row
	// id check with a string that is not a row id.
	assert.Empty(t, rowActionFromContext(map[string]any{"order": "SO-1042"}).RowID)
	assert.Empty(t, rowActionFromContext(map[string]any{"row_id": nil}).RowID)

	decoded, err := decodeRowAction(encodeRowAction(action))
	require.NoError(t, err)
	assert.Equal(t, action, decoded)

	_, err = decodeRowAction("notjson")
	assert.Error(t, err)

	_, err = decodeRowAction(`{"order":"SO-1042"}`)
	assert.Error(t, err, "a row action with no row id is rejected")
}

func TestOrderLabelAndTitle(t *testing.T) {
	filled := rowAction{RowID: "r1", Order: "SO-1042"}
	assert.Equal(t, "SO-1042", orderLabel(filled))
	assert.Equal(t, "Order SO-1042", orderTitle(filled))

	// A row added but not yet filled in still has to be addressable in prose.
	blank := rowAction{RowID: "r4"}
	assert.Equal(t, "the new order", orderLabel(blank))
	assert.Equal(t, "New order", orderTitle(blank))
}

func TestMatchFulfillmentTeam(t *testing.T) {
	team, ok := matchFulfillmentTeam("  Line 1  ")
	assert.True(t, ok)
	assert.Equal(t, "Line 1", team)

	// The submission comes from the client, so anything not offered is refused
	// rather than pasted into the post as typed.
	_, ok = matchFulfillmentTeam("Line 1</b><script>")
	assert.False(t, ok)

	_, ok = matchFulfillmentTeam("")
	assert.False(t, ok)
}

func TestRenderOrderTableEscapesCells(t *testing.T) {
	table := renderOrderTable(fillableTableState{Rows: []orderRow{
		{ID: "r1", Order: "SO-1|042", Customer: "Acme\nCo", Weight: "10", Status: "New"},
	}})

	assert.Contains(t, table, `| SO-1\|042 | Acme Co | 10 | — | New |`, "a blank cell renders as a dash")
}

func TestRowHeading(t *testing.T) {
	assert.Equal(t, "SO-1042 · Acme Co", rowHeading(orderRow{Order: "SO-1042", Customer: "Acme Co"}, 0))
	assert.Equal(t, "SO-1042", rowHeading(orderRow{Order: "SO-1042"}, 0))
	assert.Equal(t, "New order 4", rowHeading(orderRow{}, 3), "an unfilled row is still labelled")

	// DisplayName is capped at DialogElementDisplayNameMaxLength, which the
	// model counts in bytes. The separator and a customer name can both be
	// multi-byte, so a heading has to be cut on a character boundary or the
	// dialog carries invalid UTF-8 in a required field.
	for _, row := range []orderRow{
		{Order: "SO-123456789", Customer: "A Very Long Customer"},
		{Order: "SO-10", Customer: "Ωμέγα Μεταφορές ΑΕ"},
		{Order: "SO-1042678901234567890"},
	} {
		heading := rowHeading(row, 0)
		assert.LessOrEqual(t, len(heading), model.DialogElementDisplayNameMaxLength, heading)
		assert.True(t, utf8.ValidString(heading), "heading %q is not valid UTF-8", heading)
	}
}

func TestFillableTableDialogIsValid(t *testing.T) {
	state := defaultOrderRows()
	dialog := getDialogWithFillableTable(state)

	require.NoError(t, dialog.IsValid())
	assert.Equal(t, state.Size, dialog.Size, "the dialog carries the width tier the state holds")

	decoded, err := decodeFillableTableState(dialog.State)
	require.NoError(t, err)
	assert.Equal(t, state, decoded, "the grid travels in the dialog state")

	// The grid container is a collapsible whose children are one collapsible
	// per row; SubType asks the webapp to lay those rows out as a table.
	grid := dialog.Elements[0]
	require.Equal(t, "collapsible", grid.Type)
	assert.Equal(t, "grid", grid.SubType)
	require.NotNil(t, grid.CollapsibleConfig)
	require.Len(t, grid.CollapsibleConfig.Elements, len(state.Rows))

	// Every row holds the same columns in the same order, which is what lets
	// the webapp take its column headers from the first row.
	var columns []string
	for i, row := range grid.CollapsibleConfig.Elements {
		require.NotNil(t, row.CollapsibleConfig, row.Name)

		names := make([]string, 0, len(row.CollapsibleConfig.Elements))
		for _, cell := range row.CollapsibleConfig.Elements {
			names = append(names, cell.DisplayName)
		}
		if i == 0 {
			columns = names
			continue
		}
		assert.Equal(t, columns, names, "row %d has the same columns as the first", i)
	}
	assert.Equal(t, []string{"Order #", "Customer", "Weight (kg)", "Cases", "Status", "Details"}, columns)

	// Cell names are unique across the whole grid, which is what makes a
	// per-cell error address exactly one cell.
	seen := map[string]bool{}
	for _, row := range grid.CollapsibleConfig.Elements {
		for _, cell := range row.CollapsibleConfig.Elements {
			assert.False(t, seen[cell.Name], "duplicate cell name %q", cell.Name)
			seen[cell.Name] = true
		}
	}

	// The last column is a button rather than an input, and it has to name its
	// own row: a click carries only this Context, never the form values.
	details := grid.CollapsibleConfig.Elements[0].CollapsibleConfig.Elements[5]
	require.Equal(t, "action_button", details.Type)
	require.NotNil(t, details.ActionButton)
	assert.Equal(t, state.Rows[0].ID, details.ActionButton.Context["row_id"])
	assert.Equal(t, state.Rows[0].Order, details.ActionButton.Context["order"])
}

func TestFillableTableDialogStaysValidAtTheRowCap(t *testing.T) {
	state := defaultOrderRows()
	for len(state.Rows) < maxOrderRows {
		var err error
		state, err = addOrderRow(state)
		require.NoError(t, err)
	}

	dialog := getDialogWithFillableTable(state)
	require.NoError(t, dialog.IsValid(), "the dialog is valid at %d rows", maxOrderRows)
}

func TestFillableTableDialogAtEveryWidth(t *testing.T) {
	for _, width := range orderWidths {
		state := defaultOrderRows()
		state.Size = width.Value

		dialog := getDialogWithFillableTable(state)
		require.NoError(t, dialog.IsValid(), width.Value)
		assert.Equal(t, width.Value, dialog.Size)

		// The selector shows the current width, so a refresh for any other
		// reason does not silently reset it.
		assert.Equal(t, width.Value, dialog.Elements[2].Default)
	}
}

func TestFillableTableRowDetailsDialogIsValid(t *testing.T) {
	action := rowAction{RowID: "r1", Order: "SO-1042", Customer: "Acme Co"}
	dialog := getDialogFillableTableRowDetails(action)

	require.NoError(t, dialog.IsValid())
	assert.Equal(t, "Order SO-1042", dialog.Title)
	assert.Contains(t, dialog.IntroductionText, "SO-1042")

	// The submit arrives on a fresh request with no memory of the click, so
	// the row has to travel in State.
	decoded, err := decodeRowAction(dialog.State)
	require.NoError(t, err)
	assert.Equal(t, action, decoded)

	// A row added but not yet filled in still produces a valid dialog.
	blank := getDialogFillableTableRowDetails(rowAction{RowID: "r4"})
	require.NoError(t, blank.IsValid())
}

// The handler tests below drive the routes the way the server does: a
// SubmitDialogRequest in, a SubmitDialogResponse out. They are what shows that
// the flatten approach needs no bespoke payload — cell values go in as
// ordinary submission fields and errors come back in the ordinary Errors map.

func newFillableTablePlugin(t *testing.T, api *plugintest.API) *Plugin {
	t.Helper()

	plugin := &Plugin{botID: "bot-id"}
	plugin.API = api
	plugin.initializeAPI()

	return plugin
}

// postDialog sends a request through ServeHTTP, so the routes are exercised
// too, and decodes whatever the handler wrote back.
func postDialog(t *testing.T, plugin *Plugin, path string, request any) (int, model.SubmitDialogResponse) {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	defer result.Body.Close()

	raw, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	var response model.SubmitDialogResponse
	if len(raw) > 0 {
		require.NoError(t, json.Unmarshal(raw, &response), "response body: %s", raw)
	}

	return result.StatusCode, response
}

// submissionFor is what the client sends back: every cell of every row as an
// ordinary named field, valid unless a caller overrides one.
func submissionFor(state fillableTableState) map[string]any {
	submission := map[string]any{}
	for _, row := range state.Rows {
		submission[cellName(row.ID, "order")] = row.Order
		submission[cellName(row.ID, "customer")] = row.Customer
		submission[cellName(row.ID, "weight")] = "10"
		submission[cellName(row.ID, "cases")] = "2"
		submission[cellName(row.ID, "status")] = "New"
	}

	return submission
}

func TestHandleDialogFillableTableRefreshAddsARow(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	state := defaultOrderRows()
	submission := submissionFor(state)
	submission[cellName("r2", "weight")] = "7.5"
	submission[cellName("r2", "order")] = "SO-9999"
	submission["selected_field"] = fillableAddRowField
	submission[fillableAddRowField] = true

	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		Type:       "refresh",
		State:      encodeFillableTableState(state),
		Submission: submission,
	})

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "form", response.Type, "a refresh answers with a form, which swaps the open modal's contents")
	require.NotNil(t, response.Form)

	refreshed, err := decodeFillableTableState(response.Form.State)
	require.NoError(t, err)
	require.Len(t, refreshed.Rows, len(state.Rows)+1, "the row is appended")
	assert.Empty(t, refreshed.Rows[len(refreshed.Rows)-1].Order, "the new row starts blank")

	// Half-finished input survives, because the submitted cell values are
	// folded back into the rows and re-emitted as Defaults.
	assert.Equal(t, "SO-9999", refreshed.Rows[1].Order)
	assert.Equal(t, "7.5", refreshed.Rows[1].Weight)

	grid := response.Form.Elements[0]
	require.NotNil(t, grid.CollapsibleConfig)
	require.Len(t, grid.CollapsibleConfig.Elements, len(state.Rows)+1)
	assert.Equal(t, "7.5", grid.CollapsibleConfig.Elements[1].CollapsibleConfig.Elements[2].Default)

	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRefreshKeepsTheChosenWidth(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	state := defaultOrderRows()
	submission := submissionFor(state)
	submission["selected_field"] = fillableWidthField
	submission[fillableWidthField] = model.DialogSizeLarge

	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		Type:       "refresh",
		State:      encodeFillableTableState(state),
		Submission: submission,
	})

	require.Equal(t, http.StatusOK, status)
	require.NotNil(t, response.Form)
	assert.Equal(t, model.DialogSizeLarge, response.Form.Size, "the open modal resizes in place")
	assert.Len(t, response.Form.Elements[0].CollapsibleConfig.Elements, len(state.Rows), "a width change adds no rows")
}

func TestHandleDialogFillableTableRefreshAtTheRowCap(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	state := defaultOrderRows()
	for len(state.Rows) < maxOrderRows {
		var err error
		state, err = addOrderRow(state)
		require.NoError(t, err)
	}

	submission := submissionFor(state)
	submission["selected_field"] = fillableAddRowField

	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		Type:       "refresh",
		State:      encodeFillableTableState(state),
		Submission: submission,
	})

	require.Equal(t, http.StatusOK, status)
	assert.Nil(t, response.Form)
	assert.Contains(t, response.Errors, fillableAddRowField, "the cap reports on the control that hit it")
}

func TestHandleDialogFillableTableSubmitReportsPerCellErrors(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	state := defaultOrderRows()
	submission := submissionFor(state)
	submission[cellName("r2", "weight")] = "-5"

	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		State:      encodeFillableTableState(state),
		Submission: submission,
	})

	require.Equal(t, http.StatusOK, status)
	assert.Empty(t, response.Error, "a cell error is not a dialog-wide error")
	assert.Equal(t, map[string]string{
		cellName("r2", "weight"): "Enter a weight of 0 or more",
	}, response.Errors, "the error addresses one cell and leaves the modal open")

	// Nothing is posted while the form is still wrong.
	api.AssertNotCalled(t, "CreatePost", mock.Anything)
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableSubmitPostsTheQueue(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetUser", "user-id").Return(&model.User{Id: "user-id", Username: "demo_user"}, (*model.AppError)(nil)).Once()

	var posted *model.Post
	api.On("CreatePost", mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
		posted = args.Get(0).(*model.Post)
	}).Return(&model.Post{}, (*model.AppError)(nil)).Once()

	plugin := newFillableTablePlugin(t, api)

	state := defaultOrderRows()
	submission := submissionFor(state)
	submission[cellName("r1", "status")] = "Shipped"

	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		State:      encodeFillableTableState(state),
		Submission: submission,
		UserId:     "user-id",
		ChannelId:  "channel-id",
	})

	require.Equal(t, http.StatusOK, status)
	assert.Empty(t, response.Errors)

	require.NotNil(t, posted)
	assert.Equal(t, "bot-id", posted.UserId, "the queue is posted by the bot")
	assert.Equal(t, "channel-id", posted.ChannelId)
	assert.Contains(t, posted.Message, "@demo_user")
	assert.Contains(t, posted.Message, "| SO-1042 | Acme Co | 10 | 2 | Shipped |", "the submitted values reach the table")
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRejectsUnreadableState(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogWarn", mock.Anything, mock.Anything, mock.Anything).Once()
	plugin := newFillableTablePlugin(t, api)

	// State round trips through the client, so a grid that cannot be read is
	// refused rather than acted on.
	status, response := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		State:      `{"rows":[]}`,
		Submission: map[string]any{},
	})

	require.Equal(t, http.StatusOK, status)
	assert.NotEmpty(t, response.Error)
	api.AssertNotCalled(t, "CreatePost", mock.Anything)
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableCancelled(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	status, _ := postDialog(t, plugin, "/dialog/fillable-table", model.SubmitDialogRequest{
		Cancelled: true,
		State:     encodeFillableTableState(defaultOrderRows()),
	})

	assert.Equal(t, http.StatusOK, status)
	api.AssertNotCalled(t, "CreatePost", mock.Anything)
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRowActionOpensAChildDialog(t *testing.T) {
	siteURL := "http://localhost:8065"
	config := &model.Config{}
	config.SetDefaults()
	config.ServiceSettings.SiteURL = &siteURL

	api := &plugintest.API{}
	api.On("GetConfig").Return(config)

	var opened model.OpenDialogRequest
	api.On("OpenInteractiveDialog", mock.AnythingOfType("model.OpenDialogRequest")).Run(func(args mock.Arguments) {
		opened = args.Get(0).(model.OpenDialogRequest)
	}).Return((*model.AppError)(nil)).Once()

	plugin := newFillableTablePlugin(t, api)

	body, err := json.Marshal(model.PostActionIntegrationRequest{
		TriggerId: "trigger-id",
		Context: map[string]any{
			"row_id":   "r2",
			"order":    "SO-1043",
			"customer": "Globex",
		},
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/dialog/fillable-table/row-action", bytes.NewReader(body))
	plugin.ServeHTTP(nil, w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "trigger-id", opened.TriggerId, "the click brings its own trigger id")
	assert.Contains(t, opened.URL, "/dialog/fillable-table/row-details")
	require.NoError(t, opened.Dialog.IsValid())

	// The click carries no submission, so the row it refers to has to come
	// from the button's own context and travel on in State.
	action, err := decodeRowAction(opened.Dialog.State)
	require.NoError(t, err)
	assert.Equal(t, rowAction{RowID: "r2", Order: "SO-1043", Customer: "Globex"}, action)
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRowActionWithoutARow(t *testing.T) {
	api := &plugintest.API{}
	api.On("LogWarn", mock.Anything).Once()
	plugin := newFillableTablePlugin(t, api)

	body, err := json.Marshal(model.PostActionIntegrationRequest{TriggerId: "trigger-id"})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/dialog/fillable-table/row-action", bytes.NewReader(body))
	plugin.ServeHTTP(nil, w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	api.AssertNotCalled(t, "OpenInteractiveDialog", mock.Anything)
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRowDetailsPosts(t *testing.T) {
	api := &plugintest.API{}
	api.On("GetUser", "user-id").Return(&model.User{Id: "user-id", Username: "demo_user"}, (*model.AppError)(nil)).Once()

	var posted *model.Post
	api.On("CreatePost", mock.AnythingOfType("*model.Post")).Run(func(args mock.Arguments) {
		posted = args.Get(0).(*model.Post)
	}).Return(&model.Post{}, (*model.AppError)(nil)).Once()

	plugin := newFillableTablePlugin(t, api)

	status, _ := postDialog(t, plugin, "/dialog/fillable-table/row-details", model.SubmitDialogRequest{
		State:     encodeRowAction(rowAction{RowID: "r2", Order: "SO-1043", Customer: "Globex"}),
		UserId:    "user-id",
		ChannelId: "channel-id",
		Submission: map[string]any{
			"team": "Night shift",
			"note": "Two pallets\nDock 4",
		},
	})

	require.Equal(t, http.StatusOK, status)
	require.NotNil(t, posted)
	assert.Contains(t, posted.Message, "**Night shift**")
	assert.Contains(t, posted.Message, "SO-1043")
	assert.Contains(t, posted.Message, "> Two pallets\n> Dock 4", "every line of the note stays in the quote")
	api.AssertExpectations(t)
}

func TestHandleDialogFillableTableRowDetailsRejectsATeamItDidNotOffer(t *testing.T) {
	api := &plugintest.API{}
	plugin := newFillableTablePlugin(t, api)

	status, response := postDialog(t, plugin, "/dialog/fillable-table/row-details", model.SubmitDialogRequest{
		State:      encodeRowAction(rowAction{RowID: "r2", Order: "SO-1043"}),
		UserId:     "user-id",
		Submission: map[string]any{"team": "Line 9"},
	})

	require.Equal(t, http.StatusOK, status)
	assert.Contains(t, response.Errors, "team")
	api.AssertNotCalled(t, "CreatePost", mock.Anything)
	api.AssertExpectations(t)
}
