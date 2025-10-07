package main

import (
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
)

func getDialogWithSampleElements() model.Dialog {
	return model.Dialog{
		CallbackId: "somecallbackid",
		Title:      "Test Title",
		IconURL:    "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		Elements: []model.DialogElement{{
			DisplayName: "Display Name",
			Name:        "realname",
			Type:        "text",
			Default:     "default text",
			Placeholder: "placeholder",
			HelpText:    "This a regular input in an interactive dialog triggered by a test integration.",
		}, {
			DisplayName: "Email",
			Name:        dialogElementNameEmail,
			Type:        "text",
			SubType:     "email",
			Placeholder: "placeholder@bladekick.com",
			HelpText:    "This a regular email input in an interactive dialog triggered by a test integration.",
		}, {
			DisplayName: "Password",
			Name:        "somepassword",
			Type:        "text",
			SubType:     "password",
			Placeholder: "Password",
			HelpText:    "This a password input in an interactive dialog triggered by a test integration.",
		}, {
			DisplayName: "Number",
			Name:        dialogElementNameNumber,
			Type:        "text",
			SubType:     "number",
		}, {
			DisplayName: "Display Name Long Text Area",
			Name:        "realnametextarea",
			Type:        "textarea",
			Placeholder: "placeholder",
			Optional:    true,
			MinLength:   5,
			MaxLength:   100,
		}, {
			DisplayName: "User Selector",
			Name:        "someuserselector",
			Type:        "select",
			Placeholder: "Select a user...",
			HelpText:    "Choose a user from the list.",
			Optional:    true,
			MinLength:   5,
			MaxLength:   100,
			DataSource:  "users",
		}, {
			DisplayName: "Channel Selector",
			Name:        "somechannelselector",
			Type:        "select",
			Placeholder: "Select a channel...",
			HelpText:    "Choose a channel from the list.",
			Optional:    true,
			MinLength:   5,
			MaxLength:   100,
			DataSource:  "channels",
		}, {
			DisplayName: "Option Selector",
			Name:        "someoptionselector",
			Type:        "select",
			Placeholder: "Select an option...",
			HelpText:    "Choose an option from the list.",
			Options: []*model.PostActionOptions{{
				Text:  "Option1",
				Value: "opt1",
			}, {
				Text:  "Option2",
				Value: "opt2",
			}, {
				Text:  "Option3",
				Value: "opt3",
			}},
		}, {
			DisplayName: "Option Selector with default",
			Name:        "someoptionselector2",
			Type:        "select",
			Default:     "opt2",
			Placeholder: "Select an option...",
			HelpText:    "Choose an option from the list.",
			Options: []*model.PostActionOptions{{
				Text:  "Option1",
				Value: "opt1",
			}, {
				Text:  "Option2",
				Value: "opt2",
			}, {
				Text:  "Option3",
				Value: "opt3",
			}},
		}, {
			DisplayName: "Boolean Selector",
			Name:        "someboolean",
			Type:        "bool",
			Placeholder: "Agree to the terms of service",
			HelpText:    "You must agree to the terms of service to proceed.",
		}, {
			DisplayName: "Boolean Selector",
			Name:        "someboolean_optional",
			Type:        "bool",
			Placeholder: "Sign up for monthly emails?",
			HelpText:    "It's up to you if you want to get monthly emails.",
			Optional:    true,
		}, {
			DisplayName: "Boolean Selector (default true)",
			Name:        "someboolean_default_true",
			Type:        "bool",
			Placeholder: "Enable secure login",
			HelpText:    "You must enable secure login to proceed.",
			Default:     "true",
		}, {
			DisplayName: "Boolean Selector (default true)",
			Name:        "someboolean_default_true_optional",
			Type:        "bool",
			Placeholder: "Enable painfully secure login",
			HelpText:    "You may optionally enable painfully secure login.",
			Default:     "true",
			Optional:    true,
		}, {
			DisplayName: "Boolean Selector (default false)",
			Name:        "someboolean_default_false",
			Type:        "bool",
			Placeholder: "Agree to the annoying terms of service",
			HelpText:    "You must also agree to the annoying terms of service to proceed.",
			Default:     "false",
		}, {
			DisplayName: "Boolean Selector (default false)",
			Name:        "someboolean_default_false_optional",
			Type:        "bool",
			Placeholder: "Throw-away account",
			HelpText:    "A throw-away account will be deleted after 24 hours.",
			Default:     "false",
			Optional:    true,
		}, {
			DisplayName: "Radio Option Selector",
			Name:        "someradiooptionselector",
			Type:        "radio",
			HelpText:    "Choose an option from the list.",
			Options: []*model.PostActionOptions{{
				Text:  "Option1",
				Value: "opt1",
			}, {
				Text:  "Option2",
				Value: "opt2",
			}, {
				Text:  "Option3",
				Value: "opt3",
			}},
		}},
		SubmitLabel:    "Submit",
		NotifyOnCancel: true,
		State:          dialogStateSome,
	}
}

func getDialogWithoutElements(state string) model.Dialog {
	return model.Dialog{
		CallbackId:     "somecallbackid",
		Title:          "Sample Confirmation Dialog",
		IconURL:        "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		Elements:       nil,
		SubmitLabel:    "Confirm",
		NotifyOnCancel: true,
		State:          state,
	}
}

func getDialogWithIntroductionText(introductionText string) model.Dialog {
	dialog := getDialogWithSampleElements()
	dialog.IntroductionText = introductionText
	return dialog
}

func getDialogBasic() model.Dialog {
	return model.Dialog{
		CallbackId:     "basiccallbackid",
		Title:          "Simple Dialog Test",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Submit Test",
		NotifyOnCancel: true,
		State:          "somestate",
		Elements: []model.DialogElement{{
			DisplayName: "Optional Text Field",
			Name:        "optional_text",
			Type:        "text",
			Default:     "",
			Placeholder: "Enter some text (optional)...",
			HelpText:    "This field is optional for basic testing",
			Optional:    true,
			MinLength:   0,
			MaxLength:   100,
		}},
	}
}

func getDialogBoolean() model.Dialog {
	return model.Dialog{
		CallbackId:     "booleancallbackid",
		Title:          "Boolean Fields Dialog Test",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Submit Test",
		NotifyOnCancel: true,
		State:          "somestate",
		Elements: []model.DialogElement{{
			DisplayName: "Required Boolean",
			Name:        "required_boolean",
			Type:        "bool",
			Placeholder: "This field is required",
			HelpText:    "This boolean field is required and has no default (initially false)",
			Optional:    false,
		}, {
			DisplayName: "Optional Boolean",
			Name:        "optional_boolean",
			Type:        "bool",
			Placeholder: "This field is optional",
			HelpText:    "This boolean field is optional and has no default (initially false)",
			Optional:    true,
		}, {
			DisplayName: "Boolean Default True",
			Name:        "boolean_default_true",
			Type:        "bool",
			Placeholder: "This defaults to true",
			HelpText:    "This boolean field has default value true",
			Default:     "true",
			Optional:    false,
		}, {
			DisplayName: "Boolean Default False",
			Name:        "boolean_default_false",
			Type:        "bool",
			Placeholder: "This defaults to false",
			HelpText:    "This boolean field has default value false",
			Default:     "false",
			Optional:    false,
		}},
	}
}

func getDialogTextFields() model.Dialog {
	return model.Dialog{
		CallbackId:     "textfieldscallbackid",
		Title:          "Text Fields Dialog Test",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Submit Test",
		NotifyOnCancel: true,
		State:          "somestate",
		Elements: []model.DialogElement{{
			DisplayName: "Regular Text Field",
			Name:        "text_field",
			Type:        "text",
			Default:     "",
			Placeholder: "Enter some text...",
			HelpText:    "This is a regular text input",
			Optional:    true,
			MinLength:   0,
			MaxLength:   100,
		}, {
			DisplayName: "Required Text Field",
			Name:        "required_text",
			Type:        "text",
			Default:     "",
			Placeholder: "This field is required",
			HelpText:    "This field must be filled",
			Optional:    false,
			MinLength:   1,
			MaxLength:   50,
		}, {
			DisplayName: "Email Field",
			Name:        "email_field",
			Type:        "text",
			SubType:     "email",
			Default:     "",
			Placeholder: "user@example.com",
			HelpText:    "Enter a valid email address",
			Optional:    true,
			MinLength:   0,
			MaxLength:   100,
		}, {
			DisplayName: "Number Field",
			Name:        "number_field",
			Type:        "text",
			SubType:     "number",
			Default:     "",
			Placeholder: "123",
			HelpText:    "Enter a number",
			Optional:    true,
			MinLength:   0,
			MaxLength:   10,
		}, {
			DisplayName: "Password Field",
			Name:        "password_field",
			Type:        "text",
			SubType:     "password",
			Default:     "",
			Placeholder: "Enter password...",
			HelpText:    "Password field test",
			Optional:    true,
			MinLength:   0,
			MaxLength:   50,
		}, {
			DisplayName: "Text Area Field",
			Name:        "textarea_field",
			Type:        "text",
			SubType:     "textarea",
			Default:     "",
			Placeholder: "Enter multiline text...",
			HelpText:    "Text area for longer content",
			Optional:    true,
			MinLength:   0,
			MaxLength:   500,
		}},
	}
}

func getDialogSelectFields() model.Dialog {
	return model.Dialog{
		CallbackId:     "selectfieldscallbackid",
		Title:          "Select Fields Dialog Test",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Submit Test",
		NotifyOnCancel: true,
		State:          "somestate",
		Elements: []model.DialogElement{{
			DisplayName: "Radio Option Selector",
			Name:        "someradiooptions",
			Type:        "radio",
			HelpText:    "Choose your department",
			Optional:    false,
			Options: []*model.PostActionOptions{{
				Text:  "Engineering",
				Value: "engineering",
			}, {
				Text:  "Sales",
				Value: "sales",
			}},
		}, {
			DisplayName: "Option Selector",
			Name:        "someoptionselector",
			Type:        "select",
			Default:     "",
			Placeholder: "Select an option...",
			HelpText:    "",
			Optional:    false,
			Options: []*model.PostActionOptions{{
				Text:  "Option1",
				Value: "opt1",
			}, {
				Text:  "Option2",
				Value: "opt2",
			}, {
				Text:  "Option3",
				Value: "opt3",
			}},
		}, {
			DisplayName: "User Selector",
			Name:        "someuserselector",
			Type:        "select",
			Default:     "",
			Placeholder: "Select a user...",
			HelpText:    "",
			Optional:    false,
			DataSource:  "users",
		}, {
			DisplayName: "Channel Selector",
			Name:        "somechannelselector",
			Type:        "select",
			Default:     "",
			Placeholder: "Select a channel...",
			HelpText:    "Choose a channel from the list.",
			Optional:    true,
			DataSource:  "channels",
		}},
	}
}

func getDialogWithDynamicSelectElements() model.Dialog {
	return model.Dialog{
		CallbackId: "somecallbackid",
		Title:      "Dynamic Select Dialog Demo",
		IconURL:    "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		Elements: []model.DialogElement{{
			DisplayName:   "Dynamic Products",
			Name:          "dynamic_products",
			Type:          "select",
			Placeholder:   "Type to search products...",
			HelpText:      "Search for products dynamically from external API.",
			DataSource:    "dynamic",
			DataSourceURL: fmt.Sprintf("/plugins/%s/api/products", manifest.Id),
		}, {
			DisplayName:   "Dynamic Companies",
			Name:          "dynamic_companies",
			Type:          "select",
			Placeholder:   "Type to search companies...",
			HelpText:      "Search for companies dynamically based on your input.",
			DataSource:    "dynamic",
			DataSourceURL: fmt.Sprintf("/plugins/%s/api/companies", manifest.Id),
		}, {
			DisplayName:   "Dynamic Countries",
			Name:          "dynamic_countries",
			Type:          "select",
			Placeholder:   "Type to search countries...",
			HelpText:      "Search for countries dynamically with real-time filtering.",
			DataSource:    "dynamic",
			DataSourceURL: fmt.Sprintf("/plugins/%s/api/countries", manifest.Id),
			Optional:      true,
		}},
		SubmitLabel:    "Submit Dynamic Select",
		NotifyOnCancel: true,
		State:          dialogStateSome,
	}
}

// Helper function to convert interface{} to string safely
func interfaceToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
