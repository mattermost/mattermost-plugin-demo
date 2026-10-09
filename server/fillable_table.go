package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mattermost/mattermost/server/public/model"
)

// /dialog fillable-table demonstrates the "flatten" approach to a fillable
// table: the grid is a *layout* over ordinary dialog elements rather than a
// new element type with its own submission protocol.
//
// Every cell is a real element with a synthesized name ("cell_<row>_weight"),
// grouped one collapsible per row, and those row collapsibles sit inside a
// container collapsible with SubType "grid". Because collapsible children
// already flatten into the top-level submission map, and the Errors map is
// already keyed by element name, this needs no change to the submission model
// and no change to the error model:
//
//   - typed cell values arrive as ordinary named fields
//   - a per-cell error addresses exactly one cell, because names are unique
//   - Optional, MaxLength, SubType and Default all work per cell for free
//
// The sample data is a warehouse order queue: a few orders with a weight, a
// case count and a status to fill in. It is deliberately mundane, because the
// point of the demo is the layout and the submission round trip, not the
// domain.
const (
	// fillableAddRowField is the refresh-enabled toggle. Changing it appends a
	// blank row. Refreshing is how a dialog changes without submitting: the
	// handler returns a form and the open modal swaps its contents.
	fillableAddRowField = "fillable_add_row"

	// maxOrderRows caps the grid. The State round trips through the client,
	// so the cap also bounds what a submission can make the plugin render.
	maxOrderRows = 12

	// fillableWidthField is the width selector. Like the add-row toggle it is a
	// refresh field: changing it re-renders the same grid at a different
	// Dialog.Size tier, which is what makes the width visibly a property of the
	// dialog rather than something baked in at open time.
	fillableWidthField = "fillable_dialog_width"

	// fillableDetailsColumn is the last column: a per-row action_button rather
	// than an input. It carries no submission value, which is why it has no
	// field on orderRow and is absent from applySubmissionToRows and
	// validateOrders.
	fillableDetailsColumn = "details"

	// fillableTableIconURL is the modal's title icon. The mattermost.org icon
	// the older samples use 404s today, so this one points at an image that is
	// still served.
	fillableTableIconURL = "https://raw.githubusercontent.com/mattermost/mattermost/master/webapp/channels/src/images/icon50x50.png"
)

var orderStatuses = []string{"New", "Picked", "Shipped", "Delayed"}

var fulfillmentTeams = []string{"Line 1", "Line 2", "Night shift", "Courier"}

// orderWidths are the tiers offered by the width selector, in the order they
// are shown. Small is first because it is what Dialog.Size defaults to: the
// width every dialog had before the field existed.
var orderWidths = []struct {
	Value string
	Label string
}{
	{model.DialogSizeSmall, "Small - 600px (today's width)"},
	{model.DialogSizeMedium, "Medium - 700px"},
	{model.DialogSizeLarge, "Large - 900px"},
}

// orderRow is one row of the order queue. Order and Customer are the
// pre-filled row headers; Weight, Cases and Status are the fillable cells.
type orderRow struct {
	ID       string `json:"id"`
	Order    string `json:"order"`
	Customer string `json:"customer"`
	Weight   string `json:"weight"`
	Cases    string `json:"cases"`
	Status   string `json:"status"`
}

type fillableTableState struct {
	Rows []orderRow `json:"rows"`

	// Size is the dialog's width tier. It lives in State so a refresh keeps the
	// width the user picked, the same way it keeps their half-typed cells.
	Size string `json:"size"`

	// NextID keeps synthesized field names stable and unique as rows are
	// added, so a value typed into row 3 stays in row 3 across a refresh.
	NextID int `json:"next_id"`
}

func defaultOrderRows() fillableTableState {
	return fillableTableState{
		Rows: []orderRow{
			{ID: "r1", Order: "SO-1042", Customer: "Acme Co"},
			{ID: "r2", Order: "SO-1043", Customer: "Globex"},
			{ID: "r3", Order: "SO-1051", Customer: "Initech"},
		},
		NextID: 4,

		// Start at the historical width. The grid does not fit, which is the
		// point: it scrolls sideways under a sticky first column until the
		// width is raised.
		Size: model.DialogSizeSmall,
	}
}

// validOrderWidth keeps an unknown width out of the dialog. State round
// trips through the client, so the value coming back is untrusted; anything
// unrecognised falls back to small rather than failing the whole dialog.
func validOrderWidth(size string) string {
	for _, width := range orderWidths {
		if width.Value == size {
			return size
		}
	}

	return model.DialogSizeSmall
}

// applyWidthToState folds a submitted width back into the state, so the
// selector behaves like every other field on a refresh.
func applyWidthToState(state fillableTableState, submission map[string]any) fillableTableState {
	if value, ok := submission[fillableWidthField]; ok {
		state.Size = validOrderWidth(strings.TrimSpace(interfaceToString(value)))
	}

	return state
}

// Cell field names. Names are capped at model.DialogElementNameMaxLength (300),
// so there is plenty of room for a row prefix.
func cellName(rowID, column string) string {
	return fmt.Sprintf("cell_%s_%s", rowID, column)
}

func rowSectionName(rowID string) string {
	return "row_" + rowID
}

// rowAction is what a Details click carries into its child dialog. A click on
// an action_button sends only the button's static Context — not the form
// values — so this is the row as it stood when the form was last rendered,
// and it has to name its own row explicitly.
type rowAction struct {
	RowID    string `json:"row_id"`
	Order    string `json:"order"`
	Customer string `json:"customer"`
}

func rowActionFromContext(ctx map[string]any) rowAction {
	return rowAction{
		RowID:    contextString(ctx, "row_id"),
		Order:    contextString(ctx, "order"),
		Customer: contextString(ctx, "customer"),
	}
}

// contextString reads one context value. A missing key has to come back empty
// rather than go through interfaceToString, which renders a nil as "<nil>" and
// would otherwise satisfy the row id check with a string that is not a row id.
func contextString(ctx map[string]any, key string) string {
	value, ok := ctx[key]
	if !ok || value == nil {
		return ""
	}

	return strings.TrimSpace(interfaceToString(value))
}

func encodeRowAction(action rowAction) string {
	b, err := json.Marshal(action)
	if err != nil {
		// Only strings are marshaled, so this cannot happen.
		panic(err)
	}

	return string(b)
}

func decodeRowAction(raw string) (rowAction, error) {
	var action rowAction
	if err := json.Unmarshal([]byte(raw), &action); err != nil {
		return action, fmt.Errorf("failed to decode row action state: %w", err)
	}
	if action.RowID == "" {
		return action, fmt.Errorf("row action state has no row id")
	}

	return action, nil
}

// orderLabel names a row in prose. A row added but not yet filled in has no
// order number, so it falls back to something that is still addressable.
func orderLabel(action rowAction) string {
	if action.Order != "" {
		return action.Order
	}

	return "the new order"
}

// matchFulfillmentTeam resolves a submitted team to one of the offered options.
func matchFulfillmentTeam(value string) (string, bool) {
	value = strings.TrimSpace(value)
	for _, team := range fulfillmentTeams {
		if team == value {
			return team, true
		}
	}

	return "", false
}

// orderTitle is the same row as a dialog title, where the prose fallback
// would not read as one.
func orderTitle(action rowAction) string {
	if action.Order != "" {
		return "Order " + action.Order
	}

	return "New order"
}

func encodeFillableTableState(state fillableTableState) string {
	b, err := json.Marshal(state)
	if err != nil {
		// Only strings and ints are marshaled, so this cannot happen.
		panic(err)
	}

	return string(b)
}

func decodeFillableTableState(raw string) (fillableTableState, error) {
	var state fillableTableState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		return state, fmt.Errorf("failed to decode dialog state: %w", err)
	}
	if len(state.Rows) == 0 {
		return state, fmt.Errorf("dialog state has no rows")
	}
	if len(state.Rows) > maxOrderRows {
		return state, fmt.Errorf("dialog state has %d rows, more than the limit of %d", len(state.Rows), maxOrderRows)
	}
	for _, row := range state.Rows {
		if row.ID == "" {
			return state, fmt.Errorf("dialog state has a row with no id")
		}
	}
	state.Size = validOrderWidth(state.Size)

	return state, nil
}

// applySubmissionToRows copies the submitted cell values back onto the rows.
// Every refresh runs through this, which is what keeps half-finished input on
// screen when a row is added: the values come back as ordinary fields and go
// straight into the next form's Defaults.
func applySubmissionToRows(state fillableTableState, submission map[string]any) fillableTableState {
	for i, row := range state.Rows {
		if value, ok := submission[cellName(row.ID, "order")]; ok {
			state.Rows[i].Order = strings.TrimSpace(interfaceToString(value))
		}
		if value, ok := submission[cellName(row.ID, "customer")]; ok {
			state.Rows[i].Customer = strings.TrimSpace(interfaceToString(value))
		}
		if value, ok := submission[cellName(row.ID, "weight")]; ok {
			state.Rows[i].Weight = strings.TrimSpace(interfaceToString(value))
		}
		if value, ok := submission[cellName(row.ID, "cases")]; ok {
			state.Rows[i].Cases = strings.TrimSpace(interfaceToString(value))
		}
		if value, ok := submission[cellName(row.ID, "status")]; ok {
			state.Rows[i].Status = strings.TrimSpace(interfaceToString(value))
		}
	}

	return state
}

func addOrderRow(state fillableTableState) (fillableTableState, error) {
	if len(state.Rows) >= maxOrderRows {
		return state, fmt.Errorf("the order queue is full (%d rows)", maxOrderRows)
	}

	state.Rows = append(state.Rows, orderRow{
		ID:       fmt.Sprintf("r%d", state.NextID),
		Order:    "",
		Customer: "",
	})
	state.NextID++

	return state, nil
}

// validateOrders returns one error per offending cell, keyed by the cell's
// element name. This is the whole point of the flatten approach: per-cell
// error addressing falls out of the existing flat Errors map.
func validateOrders(state fillableTableState) map[string]string {
	errs := map[string]string{}

	for _, row := range state.Rows {
		if row.Order == "" {
			errs[cellName(row.ID, "order")] = "An order number is required"
		}

		if row.Weight == "" {
			errs[cellName(row.ID, "weight")] = "A weight is required"
		} else if weight, err := strconv.ParseFloat(row.Weight, 64); err != nil || weight < 0 {
			errs[cellName(row.ID, "weight")] = "Enter a weight of 0 or more"
		}

		if row.Cases != "" {
			if cases, err := strconv.Atoi(row.Cases); err != nil || cases < 0 {
				errs[cellName(row.ID, "cases")] = "Cases must be a whole number"
			}
		}

		if row.Status == "" {
			errs[cellName(row.ID, "status")] = "Choose a status"
		}
	}

	return errs
}

// renderOrderTable renders the order queue as a pipe table. Introduction text
// goes through the full Markdown renderer, so a pipe table displays as a real
// table.
func renderOrderTable(state fillableTableState) string {
	var b strings.Builder
	b.WriteString("| Order # | Customer | Weight (kg) | Cases | Status |\n")
	b.WriteString("|:--|:--|--:|--:|:--|\n")
	for _, row := range state.Rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			cellOrDash(row.Order),
			cellOrDash(row.Customer),
			cellOrDash(row.Weight),
			cellOrDash(row.Cases),
			cellOrDash(row.Status),
		)
	}

	return b.String()
}

func cellOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "—"
	}

	return escapeTableCell(value)
}

// rowHeading labels the row's section. DisplayName is capped at
// model.DialogElementDisplayNameMaxLength (24), which an order number and a
// customer name only just fit inside.
func rowHeading(row orderRow, index int) string {
	order := strings.TrimSpace(row.Order)
	customer := strings.TrimSpace(row.Customer)

	switch {
	case order != "" && customer != "":
		return truncateDisplayName(order + " · " + customer)
	case order != "":
		return truncateDisplayName(order)
	default:
		return fmt.Sprintf("New order %d", index+1)
	}
}

// truncateDisplayName keeps a heading inside
// model.DialogElementDisplayNameMaxLength, which the model counts in bytes.
// Cutting on a byte boundary can split a multi-byte character — the separator
// above is one, and a customer name may be — so the cut backs off to the last
// whole character that fits.
func truncateDisplayName(name string) string {
	if len(name) <= model.DialogElementDisplayNameMaxLength {
		return name
	}

	cut := name[:model.DialogElementDisplayNameMaxLength]
	for len(cut) > 0 && !utf8.ValidString(cut) {
		cut = cut[:len(cut)-1]
	}

	return strings.TrimSpace(cut)
}

// escapeTableCell keeps user input from breaking the table: an unescaped `|`
// starts a new cell and a newline ends the row.
func escapeTableCell(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.ReplaceAll(s, "|", `\|`)
}
