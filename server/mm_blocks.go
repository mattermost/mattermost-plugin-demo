package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

// mm_blocks action ids used in dialog fixtures. E2E tests that open dialogs via
// this plugin should tap these ids (or the matching button labels).
const (
	mmBlocksActionSubmit       = "demo_dialog_submit"
	mmBlocksActionCancel       = "demo_dialog_cancel"
	mmBlocksActionRefresh      = "demo_dialog_refresh"
	mmBlocksActionErrors       = "demo_dialog_errors"
	mmBlocksActionError        = "demo_dialog_error"
	mmBlocksActionGoto         = "demo_dialog_goto"
	mmBlocksActionLookup       = "demo_dialog_lookup"
	mmBlocksActionFieldRefresh = "demo_dialog_field_refresh"
	mmBlocksActionOpenDetails  = "demo_dialog_open_details"
	mmBlocksActionOpenSummary  = "demo_dialog_open_summary"

	mmBlocksResponsePrefix = "Demo mm_blocks"
	mmBlocksUpdatedMarker  = "DEMO_MM_BLOCKS_UPDATED"
)

type mmBlocksActionSpec struct {
	Type    string         `json:"type"`
	URL     string         `json:"url"`
	Context map[string]any `json:"context,omitempty"`
}

type mmBlocksDialogButton struct {
	Label  string `json:"label,omitempty"`
	Action string `json:"action,omitempty"`
}

// mmBlocksDialog is the blocks-mode dialog payload returned to Mattermost.
// Actions stay as a plaintext map; the server encrypts them before the client sees them.
type mmBlocksDialog struct {
	Title   string                        `json:"title"`
	IconURL string                        `json:"icon_url,omitempty"`
	State   string                        `json:"state,omitempty"`
	Submit  *mmBlocksDialogButton         `json:"submit,omitempty"`
	Cancel  *mmBlocksDialogButton         `json:"cancel,omitempty"`
	Blocks  []any                         `json:"blocks,omitempty"`
	Actions map[string]mmBlocksActionSpec `json:"actions,omitempty"`
}

type mmBlocksActionResponse struct {
	Type             string                     `json:"type,omitempty"`
	EphemeralText    string                     `json:"ephemeral_text,omitempty"`
	SkipSlackParsing bool                       `json:"skip_slack_parsing,omitempty"`
	Update           *model.Post                `json:"update,omitempty"`
	GotoLocation     string                     `json:"goto_location,omitempty"`
	Error            string                     `json:"error,omitempty"`
	Errors           map[string]string          `json:"errors,omitempty"`
	BlockDialog      *mmBlocksDialog            `json:"block_dialog,omitempty"`
	KeepDialogOpen   bool                       `json:"keep_dialog_open,omitempty"`
	Items            []model.DialogSelectOption `json:"items,omitempty"`
}

type mmBlocksDialogOptions struct {
	Title  string
	Marker string
	State  string
}

func pluginPath(path string) string {
	return fmt.Sprintf("/plugins/%s%s", manifest.Id, path)
}

func mmBlocksAction(path string, context map[string]any) mmBlocksActionSpec {
	return mmBlocksActionSpec{
		Type:    "external",
		URL:     pluginPath(path),
		Context: context,
	}
}

func (s mmBlocksActionSpec) toMap() map[string]any {
	m := map[string]any{
		"type": s.Type,
		"url":  s.URL,
	}
	if len(s.Context) > 0 {
		m["context"] = s.Context
	}
	return m
}

// mmBlocksPostAction is gob-safe for CreatePost props (plugin RPC cannot
// encode unregistered structs stored in map[string]any).
func mmBlocksPostAction(path string, context map[string]any) map[string]any {
	return mmBlocksAction(path, context).toMap()
}

func mmBlocksDialogActions(extras map[string]mmBlocksActionSpec) map[string]mmBlocksActionSpec {
	actions := map[string]mmBlocksActionSpec{
		mmBlocksActionSubmit: mmBlocksAction("/mm_blocks_dialog_submit", map[string]any{"form": "blocks_dialog"}),
		mmBlocksActionCancel: mmBlocksAction("/mm_blocks_dialog_cancel", map[string]any{"reason": "cancel"}),
	}
	for id, spec := range extras {
		actions[id] = spec
	}
	return actions
}

func baseBlockDialog(title, state, submitLabel, cancelLabel string, blocks []any, extras map[string]mmBlocksActionSpec) *mmBlocksDialog {
	if title == "" {
		title = "Demo Blocks Dialog"
	}
	if state == "" {
		state = "demo-mm-blocks-dialog"
	}
	if submitLabel == "" {
		submitLabel = "Submit"
	}
	if cancelLabel == "" {
		cancelLabel = "Cancel"
	}
	return &mmBlocksDialog{
		Title:   title,
		State:   state,
		Submit:  &mmBlocksDialogButton{Action: mmBlocksActionSubmit, Label: submitLabel},
		Cancel:  &mmBlocksDialogButton{Action: mmBlocksActionCancel, Label: cancelLabel},
		Blocks:  blocks,
		Actions: mmBlocksDialogActions(extras),
	}
}

func getMmBlocksDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Blocks Dialog"
	}
	state := opts.State
	if state == "" {
		state = "demo-mm-blocks-dialog"
	}

	intro := "Blocks dialog — fill fields, then Submit / Next step / Show errors."
	if opts.Marker != "" {
		intro = fmt.Sprintf("Blocks dialog for **%s**. Fill fields, then Submit / Next step / Show errors.", opts.Marker)
	}

	return baseBlockDialog(title, state, "", "", []any{
		map[string]any{"type": "text", "text": intro},
		map[string]any{"type": "divider"},
		map[string]any{
			"type":          "text_input",
			"name":          "title",
			"label":         "Title",
			"placeholder":   "Short title",
			"help_text":     "Required for a successful submit.",
			"initial_value": "Demo ticket",
			"max_length":    80,
		},
		map[string]any{
			"type":        "text_input",
			"name":        "email",
			"label":       "Email",
			"subtype":     "email",
			"placeholder": "you@example.com",
			"optional":    true,
		},
		map[string]any{
			"type":        "text_input",
			"name":        "description",
			"label":       "Description",
			"multiline":   true,
			"placeholder": "Longer text…",
			"optional":    true,
			"max_length":  500,
		},
		map[string]any{
			"type":          "bool_input",
			"name":          "enabled",
			"label":         "Enabled",
			"placeholder":   "Turn this on",
			"initial_value": true,
		},
		map[string]any{
			"type":        "select",
			"name":        "priority",
			"label":       "Priority",
			"placeholder": "Choose priority",
			"options": []any{
				map[string]any{"text": "Low", "value": "low"},
				map[string]any{"text": "Medium", "value": "medium"},
				map[string]any{"text": "High", "value": "high"},
			},
			"initial_option": "medium",
		},
		map[string]any{
			"type":  "select",
			"name":  "severity",
			"label": "Severity",
			"style": "expanded",
			"options": []any{
				map[string]any{"text": "SEV-1", "value": "sev1"},
				map[string]any{"text": "SEV-2", "value": "sev2"},
			},
			"initial_option": "sev2",
		},
		map[string]any{
			"type":               "select",
			"name":               "pick",
			"label":              "Dynamic option",
			"placeholder":        "Type to search…",
			"data_source":        "dynamic",
			"data_source_action": mmBlocksActionLookup,
			"optional":           true,
			"help_text":          "Options from lookup integration.",
		},
		map[string]any{
			"type":          "date_input",
			"name":          "due_date",
			"label":         "Due date",
			"optional":      true,
			"placeholder":   "Pick a due date",
			"initial_value": "2025-01-10",
		},
		map[string]any{
			"type":     "datetime_input",
			"name":     "meeting_at",
			"label":    "Meeting time",
			"optional": true,
		},
		map[string]any{
			"type":        "file_input",
			"name":        "attachments",
			"label":       "Attachments",
			"optional":    true,
			"placeholder": "Upload evidence",
			"help_text":   "Optional file upload.",
		},
		map[string]any{"type": "divider"},
		map[string]any{
			"type": "container",
			"flow": "horizontal",
			"gap":  "medium",
			"content": []any{
				map[string]any{
					"type":      "button",
					"text":      "Next step",
					"style":     "default",
					"subtype":   "submit",
					"action_id": mmBlocksActionRefresh,
				},
				map[string]any{
					"type":      "button",
					"text":      "Show errors",
					"style":     "danger",
					"subtype":   "submit",
					"action_id": mmBlocksActionErrors,
				},
				map[string]any{
					"type":      "button",
					"text":      "Top-level error",
					"style":     "danger",
					"action_id": mmBlocksActionError,
				},
				map[string]any{
					"type":      "button",
					"text":      "Navigate away",
					"style":     "default",
					"action_id": mmBlocksActionGoto,
				},
			},
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionRefresh: mmBlocksAction("/mm_blocks_dialog_refresh", map[string]any{"scenario": "refresh"}),
		mmBlocksActionErrors:  mmBlocksAction("/mm_blocks_dialog_errors", nil),
		mmBlocksActionError:   mmBlocksAction("/mm_blocks_dialog_error", nil),
		mmBlocksActionGoto:    mmBlocksAction("/mm_blocks_dialog_goto", nil),
		mmBlocksActionLookup:  mmBlocksAction("/mm_blocks_integration_lookup", nil),
	})
}

func getMmBlocksDialogStep2(previousTitle string) *mmBlocksDialog {
	if previousTitle == "" {
		previousTitle = "Demo ticket"
	}
	return &mmBlocksDialog{
		Title:  "Step 2",
		State:  "demo-mm-blocks-dialog-step-2",
		Submit: &mmBlocksDialogButton{Action: mmBlocksActionSubmit, Label: "Finish"},
		Cancel: &mmBlocksDialogButton{Action: mmBlocksActionCancel, Label: "Cancel"},
		Blocks: []any{
			map[string]any{
				"type": "text",
				"text": fmt.Sprintf("**Step 2** — refreshed from dialog. Previous title: `%s`", previousTitle),
			},
			map[string]any{
				"type":        "text_input",
				"name":        "notes",
				"label":       "Follow-up notes",
				"multiline":   true,
				"placeholder": "Anything else?",
			},
			map[string]any{
				"type":          "bool_input",
				"name":          "confirm",
				"label":         "Confirm",
				"placeholder":   "I confirm this step",
				"initial_value": false,
			},
		},
		Actions: map[string]mmBlocksActionSpec{
			mmBlocksActionSubmit: mmBlocksAction("/mm_blocks_dialog_submit", map[string]any{"step": "2"}),
			mmBlocksActionCancel: mmBlocksAction("/mm_blocks_dialog_cancel", map[string]any{"reason": "cancel", "step": "2"}),
		},
	}
}

func getMmBlocksSimpleDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Simple Dialog"
	}
	text := "Simple blocks dialog with no form fields."
	if opts.Marker != "" {
		text = fmt.Sprintf("Simple blocks dialog for **%s**.", opts.Marker)
	}
	return baseBlockDialog(title, "demo-simple", "", "", []any{
		map[string]any{"type": "text", "text": text},
	}, nil)
}

func getMmBlocksFullDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Full Dialog"
	}
	text := "Full field mix."
	if opts.Marker != "" {
		text = fmt.Sprintf("Full dialog **%s**", opts.Marker)
	}
	return baseBlockDialog(title, "demo-full", "", "", []any{
		map[string]any{"type": "text", "text": text},
		map[string]any{
			"type":        "text_input",
			"name":        "realname",
			"label":       "Name",
			"placeholder": "Enter your name",
			"help_text":   "Your full name.",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "someemail",
			"label":       "Email",
			"subtype":     "email",
			"placeholder": "you@example.com",
			"optional":    true,
		},
		map[string]any{
			"type":        "text_input",
			"name":        "somenumber",
			"label":       "Number",
			"subtype":     "number",
			"placeholder": "Enter a number",
			"optional":    true,
		},
		map[string]any{
			"type":        "text_input",
			"name":        "somepassword",
			"label":       "Password",
			"subtype":     "password",
			"placeholder": "Enter password",
			"optional":    true,
		},
		map[string]any{
			"type":        "text_input",
			"name":        "realnametextarea",
			"label":       "Notes",
			"multiline":   true,
			"placeholder": "Longer text…",
			"optional":    true,
		},
		map[string]any{
			"type":        "select",
			"name":        "someuserselector",
			"label":       "User",
			"data_source": "users",
			"placeholder": "Select a user…",
			"optional":    true,
		},
		map[string]any{
			"type":        "select",
			"name":        "somechannelselector",
			"label":       "Channel",
			"data_source": "channels",
			"placeholder": "Select a channel…",
			"optional":    true,
		},
		map[string]any{
			"type":        "select",
			"name":        "someoptionselector",
			"label":       "Option",
			"placeholder": "Select an option…",
			"options": []any{
				map[string]any{"text": "Option1", "value": "opt1"},
				map[string]any{"text": "Option2", "value": "opt2"},
				map[string]any{"text": "Option3", "value": "opt3"},
			},
			"optional": true,
		},
		map[string]any{
			"type":  "select",
			"name":  "someradiooptions",
			"label": "Radio Option",
			"style": "expanded",
			"options": []any{
				map[string]any{"text": "Engineering", "value": "engineering"},
				map[string]any{"text": "Sales", "value": "sales"},
			},
			"optional": true,
		},
		map[string]any{
			"type":          "bool_input",
			"name":          "boolean_input",
			"label":         "Boolean Selector",
			"placeholder":   "Was this modal helpful?",
			"help_text":     "This is the help text",
			"initial_value": true,
			"optional":      true,
		},
	}, nil)
}

func getMmBlocksBooleanDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Boolean Dialog"
	}
	return baseBlockDialog(title, "demo-boolean", "", "", []any{
		map[string]any{
			"type":          "bool_input",
			"name":          "boolean_input",
			"label":         "Boolean Selector",
			"placeholder":   "Was this modal helpful?",
			"help_text":     "This is the help text",
			"initial_value": true,
			"optional":      true,
		},
	}, nil)
}

func getMmBlocksUsersChannelsDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Users Channels"
	}
	return baseBlockDialog(title, "demo-users-channels", "", "", []any{
		map[string]any{
			"type":        "select",
			"name":        "someuserselector",
			"label":       "User Selector",
			"data_source": "users",
			"placeholder": "Select a user…",
		},
		map[string]any{
			"type":        "select",
			"name":        "somechannelselector",
			"label":       "Channel Selector",
			"data_source": "channels",
			"placeholder": "Select a channel…",
			"help_text":   "Choose a channel from the list.",
			"optional":    true,
		},
	}, nil)
}

func getMmBlocksMultiselectDialog(opts mmBlocksDialogOptions, includeDefaults bool) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Multiselect"
	}
	multiOptions := map[string]any{
		"type":        "select",
		"name":        "multiselect_options",
		"label":       "Multi Option Selector",
		"multiselect": true,
		"placeholder": "Select multiple options…",
		"help_text":   "You can select multiple options from this list.",
		"options": []any{
			map[string]any{"text": "Engineering", "value": "opt1"},
			map[string]any{"text": "Sales", "value": "opt2"},
			map[string]any{"text": "Marketing", "value": "opt3"},
			map[string]any{"text": "Support", "value": "opt4"},
			map[string]any{"text": "Product", "value": "opt5"},
		},
	}
	if includeDefaults {
		multiOptions["initial_options"] = []any{"opt1", "opt3"}
	}
	return baseBlockDialog(title, "demo-multiselect", "", "", []any{
		multiOptions,
		map[string]any{
			"type":        "select",
			"name":        "multiselect_users",
			"label":       "Multi User Selector",
			"multiselect": true,
			"data_source": "users",
			"placeholder": "Select multiple users…",
			"help_text":   "Choose multiple users from the team.",
		},
		map[string]any{
			"type":        "select",
			"name":        "single_select_options",
			"label":       "Single Option Selector",
			"placeholder": "Select one option…",
			"options": []any{
				map[string]any{"text": "Engineering", "value": "opt1"},
				map[string]any{"text": "Sales", "value": "opt2"},
				map[string]any{"text": "Marketing", "value": "opt3"},
			},
			"optional": true,
		},
	}, nil)
}

func getMmBlocksDynamicDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Dynamic Select"
	}
	return baseBlockDialog(title, "demo-dynamic", "", "", []any{
		map[string]any{
			"type":               "select",
			"name":               "dynamic_role_selector",
			"label":              "Role",
			"placeholder":        "Type to search roles…",
			"data_source":        "dynamic",
			"data_source_action": mmBlocksActionLookup,
			"help_text":          "Required dynamic select.",
		},
		map[string]any{
			"type":               "select",
			"name":               "optional_dynamic_selector",
			"label":              "Optional Role",
			"placeholder":        "Optional search…",
			"data_source":        "dynamic",
			"data_source_action": mmBlocksActionLookup,
			"optional":           true,
			"initial_option":     "opt_beta",
			"help_text":          "Optional dynamic select with default.",
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionLookup: mmBlocksAction("/mm_blocks_integration_lookup", nil),
	})
}

func getMmBlocksEmptyRequiredDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Required Fields"
	}
	return baseBlockDialog(title, "demo-required", "", "", []any{
		map[string]any{
			"type":        "text_input",
			"name":        "realname",
			"label":       "Name",
			"placeholder": "Enter your name",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "someemail",
			"label":       "Email",
			"subtype":     "email",
			"placeholder": "you@example.com",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "somenumber",
			"label":       "Number",
			"subtype":     "number",
			"placeholder": "Enter a number",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "somepassword",
			"label":       "Password",
			"subtype":     "password",
			"placeholder": "Enter password",
			"optional":    true,
		},
		map[string]any{
			"type":          "bool_input",
			"name":          "boolean_input",
			"label":         "Boolean Selector",
			"placeholder":   "Was this modal helpful?",
			"help_text":     "This is the help text",
			"initial_value": true,
			"optional":      true,
		},
	}, nil)
}

func getMmBlocksFileUploadDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo File Upload"
	}
	return baseBlockDialog(title, "demo-file-upload", "Submit Files", "", []any{
		map[string]any{
			"type":        "file_input",
			"name":        "single_document",
			"label":       "Upload Single Document",
			"placeholder": "Select one document…",
			"help_text":   "Upload a single document (replaces previous selection).",
		},
		map[string]any{
			"type":           "file_input",
			"name":           "multiple_files",
			"label":          "Upload Multiple Files",
			"allow_multiple": true,
			"placeholder":    "Select multiple files…",
			"help_text":      "Upload multiple files (can select and add more).",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "description",
			"label":       "Description",
			"multiline":   true,
			"placeholder": "Describe the uploaded files…",
			"optional":    true,
			"max_length":  500,
		},
	}, nil)
}

func getMmBlocksFieldRefreshDialog(projectType, projectName string) *mmBlocksDialog {
	nameInput := map[string]any{
		"type":        "text_input",
		"name":        "project_name",
		"label":       "Project Name",
		"placeholder": "Enter project name",
	}
	if projectName != "" {
		nameInput["initial_value"] = projectName
	}
	typeSelect := map[string]any{
		"type":        "select",
		"name":        "project_type",
		"label":       "Project Type",
		"placeholder": "Select project type…",
		"onChange":    mmBlocksActionFieldRefresh,
		"options": []any{
			map[string]any{"text": "Web Application", "value": "web"},
			map[string]any{"text": "Mobile App", "value": "mobile"},
			map[string]any{"text": "API Service", "value": "api"},
		},
	}
	if projectType != "" {
		typeSelect["initial_option"] = projectType
	}

	blocks := []any{
		map[string]any{"type": "text", "text": "Enter project name then select type to see different fields"},
		nameInput,
		typeSelect,
	}

	switch projectType {
	case "web":
		blocks = append(blocks, map[string]any{
			"type":        "select",
			"name":        "framework",
			"label":       "Framework",
			"placeholder": "Select framework…",
			"options": []any{
				map[string]any{"text": "React", "value": "react"},
				map[string]any{"text": "Vue", "value": "vue"},
				map[string]any{"text": "Angular", "value": "angular"},
			},
			"optional": true,
		})
	case "mobile":
		blocks = append(blocks, map[string]any{
			"type":        "select",
			"name":        "platform",
			"label":       "Platform",
			"placeholder": "Select platform…",
			"options": []any{
				map[string]any{"text": "iOS", "value": "ios"},
				map[string]any{"text": "Android", "value": "android"},
				map[string]any{"text": "React Native", "value": "react-native"},
			},
			"optional": true,
		})
	case "api":
		blocks = append(blocks, map[string]any{
			"type":        "select",
			"name":        "language",
			"label":       "Language",
			"placeholder": "Select language…",
			"options": []any{
				map[string]any{"text": "Go", "value": "go"},
				map[string]any{"text": "Node.js", "value": "nodejs"},
				map[string]any{"text": "Python", "value": "python"},
			},
			"optional": true,
		})
	}

	return baseBlockDialog("Field Refresh Demo", "demo-field-refresh", "", "", blocks, map[string]mmBlocksActionSpec{
		mmBlocksActionFieldRefresh: mmBlocksAction("/mm_blocks_dialog_field_refresh", nil),
	})
}

func getMmBlocksMultistep1Dialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Step 1 - Personal Info"
	}
	return baseBlockDialog(title, "step1", "Next Step", "", []any{
		map[string]any{"type": "text", "text": "Multi-step registration - Step 1 of 3"},
		map[string]any{
			"type":        "text_input",
			"name":        "first_name",
			"label":       "First Name",
			"placeholder": "Enter your first name",
		},
		map[string]any{
			"type":        "text_input",
			"name":        "email",
			"label":       "Email",
			"subtype":     "email",
			"placeholder": "Enter your email address",
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionSubmit: mmBlocksAction("/mm_blocks_dialog_multistep", map[string]any{"step": "1"}),
	})
}

func getMmBlocksMultistep2Dialog() *mmBlocksDialog {
	return baseBlockDialog("Step 2 - Work Info", "step2", "Next Step", "", []any{
		map[string]any{"type": "text", "text": "Multi-step registration - Step 2 of 3"},
		map[string]any{
			"type":        "select",
			"name":        "department",
			"label":       "Department",
			"placeholder": "Select department…",
			"options": []any{
				map[string]any{"text": "Engineering", "value": "engineering"},
				map[string]any{"text": "Marketing", "value": "marketing"},
				map[string]any{"text": "Sales", "value": "sales"},
			},
		},
		map[string]any{
			"type":  "select",
			"name":  "experience_level",
			"label": "Experience Level",
			"style": "expanded",
			"options": []any{
				map[string]any{"text": "Junior", "value": "junior"},
				map[string]any{"text": "Mid-level", "value": "mid"},
				map[string]any{"text": "Senior", "value": "senior"},
			},
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionSubmit: mmBlocksAction("/mm_blocks_dialog_multistep", map[string]any{"step": "2"}),
	})
}

func getMmBlocksMultistep3Dialog() *mmBlocksDialog {
	return baseBlockDialog("Step 3 - Final Details", "step3", "Complete Registration", "", []any{
		map[string]any{"type": "text", "text": "Multi-step registration - Step 3 of 3"},
		map[string]any{
			"type":        "text_input",
			"name":        "comments",
			"label":       "Comments",
			"multiline":   true,
			"placeholder": "Any additional comments…",
			"optional":    true,
		},
		map[string]any{
			"type":        "bool_input",
			"name":        "terms_accepted",
			"label":       "Terms & Conditions",
			"placeholder": "I accept the terms",
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionSubmit: mmBlocksAction("/mm_blocks_dialog_submit", map[string]any{"step": "3", "form": "multistep"}),
	})
}

func getMmBlocksChildContentDialog(source string) *mmBlocksDialog {
	if source == "" {
		source = "Unknown"
	}
	title := source + " Dialog"
	if len(title) > 24 {
		title = title[:24]
	}
	return baseBlockDialog(title, "demo-child-"+source, "", "", []any{
		map[string]any{
			"type": "text",
			"text": fmt.Sprintf("This view was opened from the **%s** button (stacked modal via dialogs/open).", source),
		},
		map[string]any{
			"type":        "text_input",
			"name":        "child_input",
			"label":       "Child Input",
			"placeholder": "Enter value",
			"optional":    true,
		},
	}, nil)
}

func getMmBlocksActionParentDialog(opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo Action Buttons"
	}
	return baseBlockDialog(title, "demo-action-parent", "", "", []any{
		map[string]any{
			"type":        "text_input",
			"name":        "your_name",
			"label":       "Your Name",
			"placeholder": "Enter your name",
			"optional":    true,
		},
		map[string]any{
			"type": "container",
			"flow": "horizontal",
			"gap":  "medium",
			"content": []any{
				map[string]any{
					"type":      "button",
					"text":      "Open Details",
					"style":     "primary",
					"action_id": mmBlocksActionOpenDetails,
				},
				map[string]any{
					"type":      "button",
					"text":      "Open Summary",
					"style":     "default",
					"action_id": mmBlocksActionOpenSummary,
				},
			},
		},
	}, map[string]mmBlocksActionSpec{
		mmBlocksActionOpenDetails: mmBlocksAction("/mm_blocks_dialog_child", map[string]any{"source": "Details"}),
		mmBlocksActionOpenSummary: mmBlocksAction("/mm_blocks_dialog_child", map[string]any{"source": "Summary"}),
	})
}

func getMmBlocksDatetimeDialog(scenario string, opts mmBlocksDialogOptions) *mmBlocksDialog {
	title := opts.Title
	if title == "" {
		title = "Demo DateTime"
	}

	var blocks []any
	switch scenario {
	case "datetime_basic":
		blocks = []any{
			map[string]any{
				"type":        "date_input",
				"name":        "event_date",
				"label":       "Event Date",
				"placeholder": "Select a date",
				"help_text":   "Select the date for your event",
			},
			map[string]any{
				"type":            "datetime_input",
				"name":            "meeting_time",
				"label":           "Meeting Time",
				"placeholder":     "Select date and time",
				"help_text":       "Select the date and time for your meeting",
				"optional":        true,
				"datetime_config": map[string]any{"time_interval": 60},
			},
		}
	case "datetime_mindate":
		blocks = []any{
			map[string]any{
				"type":            "date_input",
				"name":            "future_date",
				"label":           "Future Date Only",
				"placeholder":     "Select a future date",
				"help_text":       "Must be today or later",
				"optional":        true,
				"datetime_config": map[string]any{"min_date": "today"},
			},
		}
	case "datetime_interval":
		blocks = []any{
			map[string]any{
				"type":            "datetime_input",
				"name":            "interval_time",
				"label":           "Custom Interval Time",
				"placeholder":     "Select time (30min intervals)",
				"help_text":       "Time picker with 30-minute intervals",
				"optional":        true,
				"datetime_config": map[string]any{"time_interval": 30},
			},
		}
	case "datetime_relative":
		blocks = []any{
			map[string]any{
				"type":          "date_input",
				"name":          "relative_date",
				"label":         "Relative Date Example",
				"placeholder":   "Today by default",
				"help_text":     "Defaults to today using relative date",
				"optional":      true,
				"initial_value": "today",
			},
			map[string]any{
				"type":          "datetime_input",
				"name":          "relative_datetime",
				"label":         "Relative DateTime Example",
				"placeholder":   "Tomorrow by default",
				"help_text":     "Defaults to tomorrow using relative date",
				"optional":      true,
				"initial_value": "+1d",
			},
		}
	case "datetime_timezone":
		blocks = []any{
			map[string]any{
				"type":      "datetime_input",
				"name":      "london_dropdown",
				"label":     "London Office Hours",
				"help_text": "Times shown in GMT - select from 60 min intervals",
				"optional":  true,
				"datetime_config": map[string]any{
					"location_timezone": "Europe/London",
					"time_interval":     60,
				},
			},
		}
	case "datetime_manual":
		blocks = []any{
			map[string]any{
				"type":            "datetime_input",
				"name":            "local_manual",
				"label":           "Your Local Time",
				"help_text":       "Type any time: 9am, 14:30, 3:45pm - no rounding",
				"optional":        true,
				"datetime_config": map[string]any{"manual_time_entry": true},
			},
			map[string]any{
				"type":      "datetime_input",
				"name":      "london_manual",
				"label":     "London Manual Entry",
				"help_text": "Type time in GMT: 9am, 14:30, 3:45pm - no rounding",
				"optional":  true,
				"datetime_config": map[string]any{
					"location_timezone": "Europe/London",
					"manual_time_entry": true,
				},
			},
		}
	default:
		blocks = []any{
			map[string]any{
				"type":        "date_input",
				"name":        "event_date",
				"label":       "Event Date",
				"placeholder": "Select a date",
				"optional":    true,
			},
		}
	}

	return baseBlockDialog(title, "demo-"+scenario, "", "", blocks, nil)
}

func getMmBlocksDialogByScenario(scenario string, opts mmBlocksDialogOptions) *mmBlocksDialog {
	switch scenario {
	case "simple":
		return getMmBlocksSimpleDialog(opts)
	case "full":
		return getMmBlocksFullDialog(opts)
	case "boolean":
		return getMmBlocksBooleanDialog(opts)
	case "users_channels":
		return getMmBlocksUsersChannelsDialog(opts)
	case "multiselect":
		return getMmBlocksMultiselectDialog(opts, false)
	case "multiselect_defaults":
		return getMmBlocksMultiselectDialog(opts, true)
	case "dynamic":
		return getMmBlocksDynamicDialog(opts)
	case "empty_required":
		return getMmBlocksEmptyRequiredDialog(opts)
	case "file_upload":
		return getMmBlocksFileUploadDialog(opts)
	case "field_refresh":
		return getMmBlocksFieldRefreshDialog("", "")
	case "multistep_1":
		return getMmBlocksMultistep1Dialog(opts)
	case "multistep_2":
		return getMmBlocksMultistep2Dialog()
	case "multistep_3":
		return getMmBlocksMultistep3Dialog()
	case "action_parent":
		return getMmBlocksActionParentDialog(opts)
	case "datetime_basic", "datetime_mindate", "datetime_interval", "datetime_relative", "datetime_timezone", "datetime_manual":
		return getMmBlocksDatetimeDialog(scenario, opts)
	default:
		return getMmBlocksDialog(opts)
	}
}

func getMmBlocksLookupOptions(searchText string) []model.DialogSelectOption {
	allOptions := []model.DialogSelectOption{
		{Text: "Alpha", Value: "opt_alpha"},
		{Text: "Beta", Value: "opt_beta"},
		{Text: "Gamma", Value: "opt_gamma"},
		{Text: "Mattermost", Value: "opt_mm"},
	}
	search := strings.ToLower(searchText)
	if search == "" {
		return allOptions
	}
	var filtered []model.DialogSelectOption
	for _, option := range allOptions {
		if strings.Contains(strings.ToLower(option.Text), search) || strings.Contains(strings.ToLower(option.Value), search) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func formatFormValuesSummary(formValues map[string]any) string {
	keys := make([]string, 0, len(formValues))
	for k := range formValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+formatFormValue(formValues[k]))
	}
	return strings.Join(parts, "&")
}

func formatFormValue(v any) string {
	switch t := v.(type) {
	case []any:
		items := make([]string, len(t))
		for i, item := range t {
			items[i] = fmt.Sprint(item)
		}
		return strings.Join(items, ",")
	case []string:
		return strings.Join(t, ",")
	default:
		return fmt.Sprint(v)
	}
}

func getUpstreamFormValues(ctx map[string]any) map[string]any {
	if ctx == nil {
		return map[string]any{}
	}
	fv, ok := ctx["form_values"].(map[string]any)
	if !ok || fv == nil {
		return map[string]any{}
	}
	return fv
}

func contextString(ctx map[string]any, key string) string {
	if ctx == nil {
		return ""
	}
	s, ok := ctx[key].(string)
	if !ok {
		return ""
	}
	return s
}

func contextValueAsString(ctx map[string]any, key string) string {
	if ctx == nil {
		return "null"
	}
	v, ok := ctx[key]
	if !ok || v == nil {
		return "null"
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

func mmBlocksDialogScenarios() []string {
	return []string{
		"simple",
		"full",
		"boolean",
		"users_channels",
		"multiselect",
		"multiselect_defaults",
		"dynamic",
		"empty_required",
		"file_upload",
		"field_refresh",
		"multistep_1",
		"action_parent",
		"datetime_basic",
		"datetime_mindate",
		"datetime_interval",
		"datetime_relative",
		"datetime_timezone",
		"datetime_manual",
	}
}
