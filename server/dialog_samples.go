package main

import (
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
)

// Sample Dialog with Field Refresh Functionality
// This demonstrates how field changes can trigger form refreshes
func getDialogWithFieldRefresh(projectType string) model.Dialog {
	dialog := model.Dialog{
		CallbackId:     "field_refresh_demo",
		Title:          "Project Configuration",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Create Project",
		NotifyOnCancel: true,
		State:          "field_refresh_state",
		SourceURL:      fmt.Sprintf("/plugins/%s/dialog/field-refresh", manifest.Id), // NEW: Source URL for refreshing
		Elements: []model.DialogElement{{
			DisplayName: "Project Type",
			Name:        "project_type",
			Type:        "select",
			Placeholder: "Select a project type...",
			HelpText:    "Choose the type of project you want to create.",
			Refresh:     true, // NEW: This field will refresh the form when changed
			Options: []*model.PostActionOptions{{
				Text:  "Web Application",
				Value: "web",
			}, {
				Text:  "Mobile Application",
				Value: "mobile",
			}, {
				Text:  "Desktop Application",
				Value: "desktop",
			}, {
				Text:  "API Service",
				Value: "api",
			}},
			Default: projectType,
		}},
	}

	// Add project-type-specific fields based on current selection
	switch projectType {
	case "web":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Frontend Framework",
				Name:        "frontend_framework",
				Type:        "select",
				Placeholder: "Select a framework...",
				HelpText:    "Choose your preferred frontend framework.",
				Options: []*model.PostActionOptions{{
					Text:  "React",
					Value: "react",
				}, {
					Text:  "Vue.js",
					Value: "vue",
				}, {
					Text:  "Angular",
					Value: "angular",
				}},
			},
			model.DialogElement{
				DisplayName: "Enable PWA",
				Name:        "enable_pwa",
				Type:        "bool",
				Placeholder: "Enable Progressive Web App features",
				HelpText:    "Add service worker and offline capabilities.",
				Optional:    true,
			},
		)
	case "mobile":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Platform",
				Name:        "mobile_platform",
				Type:        "select",
				Placeholder: "Select platform...",
				HelpText:    "Choose the mobile platform to target.",
				Options: []*model.PostActionOptions{{
					Text:  "React Native",
					Value: "react-native",
				}, {
					Text:  "Flutter",
					Value: "flutter",
				}, {
					Text:  "Native iOS",
					Value: "ios",
				}, {
					Text:  "Native Android",
					Value: "android",
				}},
			},
			model.DialogElement{
				DisplayName: "Minimum OS Version",
				Name:        "min_os_version",
				Type:        "text",
				Placeholder: "e.g., iOS 14.0 or Android 10",
				HelpText:    "Specify the minimum supported OS version.",
			},
		)
	case "desktop":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Desktop Framework",
				Name:        "desktop_framework",
				Type:        "select",
				Placeholder: "Select a framework...",
				HelpText:    "Choose your desktop application framework.",
				Options: []*model.PostActionOptions{{
					Text:  "Electron",
					Value: "electron",
				}, {
					Text:  "Tauri",
					Value: "tauri",
				}, {
					Text:  "Qt",
					Value: "qt",
				}, {
					Text:  ".NET WPF",
					Value: "wpf",
				}},
			},
			model.DialogElement{
				DisplayName: "Auto-Updater",
				Name:        "auto_updater",
				Type:        "bool",
				Placeholder: "Enable automatic updates",
				HelpText:    "Allow the application to update itself automatically.",
				Default:     "true",
			},
		)
	case "api":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "API Type",
				Name:        "api_type",
				Type:        "radio",
				HelpText:    "Choose the type of API to create.",
				Options: []*model.PostActionOptions{{
					Text:  "REST API",
					Value: "rest",
				}, {
					Text:  "GraphQL API",
					Value: "graphql",
				}, {
					Text:  "gRPC Service",
					Value: "grpc",
				}},
			},
			model.DialogElement{
				DisplayName: "Database",
				Name:        "database",
				Type:        "select",
				Placeholder: "Select a database...",
				HelpText:    "Choose your preferred database.",
				Options: []*model.PostActionOptions{{
					Text:  "PostgreSQL",
					Value: "postgresql",
				}, {
					Text:  "MySQL",
					Value: "mysql",
				}, {
					Text:  "MongoDB",
					Value: "mongodb",
				}, {
					Text:  "Redis",
					Value: "redis",
				}},
			},
		)
	}

	// Add common fields for all project types
	if projectType != "" {
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Project Name",
				Name:        "project_name",
				Type:        "text",
				Placeholder: "Enter project name...",
				HelpText:    "Name of your new project.",
				MinLength:   3,
				MaxLength:   50,
			},
			model.DialogElement{
				DisplayName: "Description",
				Name:        "description",
				Type:        "textarea",
				Placeholder: "Describe your project...",
				HelpText:    "Brief description of what this project will do.",
				Optional:    true,
				MaxLength:   500,
			},
		)
	}

	return dialog
}

// Sample Dialog with Form Refresh (Multi-Step) Functionality
// This demonstrates how submit can return a new form for multi-step workflows
func getDialogStep1() model.Dialog {
	return model.Dialog{
		CallbackId:     "multistep_demo_step1",
		Title:          "User Registration - Step 1",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Next Step",
		NotifyOnCancel: true,
		State:          "step1",
		Elements: []model.DialogElement{{
			DisplayName: "User Type",
			Name:        "user_type",
			Type:        "radio",
			HelpText:    "What type of user account do you want to create?",
			Options: []*model.PostActionOptions{{
				Text:  "Individual Developer",
				Value: "individual",
			}, {
				Text:  "Team/Organization",
				Value: "organization",
			}, {
				Text:  "Student",
				Value: "student",
			}},
		}, {
			DisplayName: "Primary Use Case",
			Name:        "use_case",
			Type:        "select",
			Placeholder: "Select your primary use case...",
			HelpText:    "What will you primarily use this platform for?",
			Options: []*model.PostActionOptions{{
				Text:  "Software Development",
				Value: "development",
			}, {
				Text:  "Project Management",
				Value: "project_mgmt",
			}, {
				Text:  "Team Communication",
				Value: "communication",
			}, {
				Text:  "Learning/Education",
				Value: "education",
			}},
		}, {
			DisplayName: "First Name",
			Name:        "first_name",
			Type:        "text",
			Placeholder: "Enter your first name...",
			HelpText:    "Your first name.",
			MinLength:   2,
			MaxLength:   50,
		}, {
			DisplayName: "Last Name",
			Name:        "last_name",
			Type:        "text",
			Placeholder: "Enter your last name...",
			HelpText:    "Your last name.",
			MinLength:   2,
			MaxLength:   50,
		}},
	}
}

func getDialogStep2(userType, useCase string) model.Dialog {
	dialog := model.Dialog{
		CallbackId:     "multistep_demo_step2",
		Title:          "User Registration - Step 2",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		SubmitLabel:    "Complete Registration",
		NotifyOnCancel: true,
		State:          "step2",
		Elements:       []model.DialogElement{},
	}

	// Add fields based on user type from step 1
	switch userType {
	case "individual":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Years of Experience",
				Name:        "experience_years",
				Type:        "text",
				SubType:     "number",
				Placeholder: "e.g., 5",
				HelpText:    "How many years of professional experience do you have?",
				MinLength:   1,
				MaxLength:   2,
			},
			model.DialogElement{
				DisplayName: "Programming Languages",
				Name:        "languages",
				Type:        "textarea",
				Placeholder: "e.g., JavaScript, Python, Go...",
				HelpText:    "List the programming languages you're familiar with.",
				Optional:    true,
				MaxLength:   200,
			},
		)
	case "organization":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Company Name",
				Name:        "company_name",
				Type:        "text",
				Placeholder: "Enter company name...",
				HelpText:    "Name of your organization.",
				MinLength:   2,
				MaxLength:   100,
			},
			model.DialogElement{
				DisplayName: "Team Size",
				Name:        "team_size",
				Type:        "select",
				Placeholder: "Select team size...",
				HelpText:    "How many people are in your team?",
				Options: []*model.PostActionOptions{{
					Text:  "1-5 people",
					Value: "small",
				}, {
					Text:  "6-20 people",
					Value: "medium",
				}, {
					Text:  "21-100 people",
					Value: "large",
				}, {
					Text:  "100+ people",
					Value: "enterprise",
				}},
			},
			model.DialogElement{
				DisplayName: "Industry",
				Name:        "industry",
				Type:        "text",
				Placeholder: "e.g., Technology, Healthcare, Finance...",
				HelpText:    "What industry is your organization in?",
				Optional:    true,
				MaxLength:   50,
			},
		)
	case "student":
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Educational Institution",
				Name:        "school",
				Type:        "text",
				Placeholder: "Enter your school/university...",
				HelpText:    "Name of your educational institution.",
				MinLength:   2,
				MaxLength:   100,
			},
			model.DialogElement{
				DisplayName: "Study Level",
				Name:        "study_level",
				Type:        "select",
				Placeholder: "Select your level...",
				HelpText:    "What level are you studying at?",
				Options: []*model.PostActionOptions{{
					Text:  "High School",
					Value: "high_school",
				}, {
					Text:  "Undergraduate",
					Value: "undergraduate",
				}, {
					Text:  "Graduate",
					Value: "graduate",
				}, {
					Text:  "PhD",
					Value: "phd",
				}},
			},
			model.DialogElement{
				DisplayName: "Field of Study",
				Name:        "field_of_study",
				Type:        "text",
				Placeholder: "e.g., Computer Science, Engineering...",
				HelpText:    "What is your major/field of study?",
				Optional:    true,
				MaxLength:   50,
			},
		)
	}

	// Add fields based on use case
	if useCase == "development" {
		dialog.Elements = append(dialog.Elements,
			model.DialogElement{
				DisplayName: "Preferred Development Environment",
				Name:        "dev_environment",
				Type:        "select",
				Placeholder: "Select environment...",
				HelpText:    "What development environment do you prefer?",
				Options: []*model.PostActionOptions{{
					Text:  "VS Code",
					Value: "vscode",
				}, {
					Text:  "IntelliJ IDEA",
					Value: "intellij",
				}, {
					Text:  "Vim/Neovim",
					Value: "vim",
				}, {
					Text:  "Emacs",
					Value: "emacs",
				}, {
					Text:  "Other",
					Value: "other",
				}},
				Optional: true,
			},
		)
	}

	// Add common notification preferences
	dialog.Elements = append(dialog.Elements,
		model.DialogElement{
			DisplayName: "Email Notifications",
			Name:        "email_notifications",
			Type:        "bool",
			Placeholder: "Receive email notifications",
			HelpText:    "Get notified about important updates via email.",
			Default:     "true",
		},
		model.DialogElement{
			DisplayName: "Newsletter Subscription",
			Name:        "newsletter",
			Type:        "bool",
			Placeholder: "Subscribe to our newsletter",
			HelpText:    "Receive monthly updates about new features and tips.",
			Optional:    true,
		},
	)

	return dialog
}

// Generate a summary dialog showing collected information
func getDialogStep3Summary(formData map[string]interface{}) model.Dialog {
	// Create a summary of the collected information from all accumulated values
	summaryText := "## Registration Summary\n\n"

	// Add key information from accumulated form data
	if userType, ok := formData["user_type"].(string); ok {
		summaryText += fmt.Sprintf("**User Type:** %s\n", userType)
	}
	if useCase, ok := formData["use_case"].(string); ok {
		summaryText += fmt.Sprintf("**Use Case:** %s\n", useCase)
	}
	if firstName, ok := formData["first_name"].(string); ok {
		if lastName, ok := formData["last_name"].(string); ok {
			summaryText += fmt.Sprintf("**Name:** %s %s\n", firstName, lastName)
		}
	}

	// Add other accumulated information
	summaryText += "\n**Additional Information:**\n"
	for key, value := range formData {
		// Skip already displayed fields
		if key == "user_type" || key == "use_case" || key == "first_name" || key == "last_name" {
			continue
		}
		if str := interfaceToString(value); str != "" {
			summaryText += fmt.Sprintf("- **%s:** %s\n", key, str)
		}
	}

	return model.Dialog{
		CallbackId:       "multistep_demo_final",
		Title:            "Confirm Registration",
		IconURL:          "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		IntroductionText: summaryText,
		SubmitLabel:      "Confirm & Complete",
		NotifyOnCancel:   true,
		State:            "final",
		Elements: []model.DialogElement{{
			DisplayName: "Terms & Conditions",
			Name:        "accept_terms",
			Type:        "bool",
			Placeholder: "I accept the Terms & Conditions",
			HelpText:    "You must accept our terms to complete registration.",
		}, {
			DisplayName: "Privacy Policy",
			Name:        "accept_privacy",
			Type:        "bool",
			Placeholder: "I accept the Privacy Policy",
			HelpText:    "You must accept our privacy policy to complete registration.",
		}},
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
