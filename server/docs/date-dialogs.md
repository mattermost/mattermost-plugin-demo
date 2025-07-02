## Date and DateTime Interactive Dialogs

The `/dialog date` command demonstrates the new Date and DateTime field types in [Interactive Dialogs](https://docs.mattermost.com/developer/interactive-dialogs.html).

### Features Demonstrated

This dialog showcases:
- **Date Picker Field**: Select a date using a calendar-style date picker
- **DateTime Picker Field**: Select both date and time using combined date and time selectors
- **Integration**: How date/datetime fields work alongside traditional form elements

### Usage

#### `/dialog date`:
Opens an Interactive Dialog with both date and datetime fields for testing:

- **Meeting Date**: A date-only field using `type: "date"`
- **Meeting Date & Time**: A datetime field using `type: "datetime"` 
- Additional form fields (text, textarea, boolean) to show integration

#### `/dialog help`:
Shows help text and usage information for all dialog variants

### Implementation Details

The date fields are implemented as new dialog element types:

```go
{
    DisplayName: "Meeting Date",
    Name:        "somedate",
    Type:        "date",
    Placeholder: "Select a meeting date",
    HelpText:    "Choose a date for the meeting using the date picker.",
}
```

```go
{
    DisplayName: "Meeting Date & Time", 
    Name:        "somedatetime",
    Type:        "datetime",
    Placeholder: "Select meeting date and time",
    HelpText:    "Choose both date and time for the meeting start.",
}
```

### Value Format

- **Date fields** return ISO date strings: `"2025-01-15"`
- **DateTime fields** return ISO datetime strings in UTC: `"2025-01-15T14:30:00Z"`

### User Experience

- Date fields display a calendar popup for date selection
- DateTime fields show both date picker and time selector
- Fields integrate seamlessly with existing form validation and submission
- Localized date formatting based on user preferences
- Timezone-aware datetime handling (stored in UTC, displayed in user timezone)

This implementation provides a comprehensive test environment for the new date and datetime field types in Mattermost Interactive Dialogs.