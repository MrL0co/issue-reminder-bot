// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package issue_reminder_bot

type Integration interface {
	IsIntegration()
}

type AssignIssue struct {
	IssueID string `json:"issueId"`
	UserID  string `json:"userId"`
}

type EditIssue struct {
	IssueID     string `json:"issueId"`
	Description string `json:"description"`
}

type Issue struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Assignee    *User  `json:"assignee"`
}

type MatterMostIntegration struct {
	Driver        string `json:"driver"`
	ServerAddress string `json:"serverAddress"`
	Name          string `json:"name"`
}

func (MatterMostIntegration) IsIntegration() {}

type NewIssue struct {
	Description string  `json:"description"`
	UserID      *string `json:"userId"`
}

type NewUser struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
}

type ServerConfig struct {
	Integrations []Integration `json:"integrations"`
}

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    *string  `json:"email"`
	Username *string  `json:"username"`
	Issues   []*Issue `json:"issues"`
}
