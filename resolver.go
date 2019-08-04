//go:generate go run github.com/99designs/gqlgen

package issue_reminder_bot

import (
	"context"
	"errors"
	"strconv"
)

type Resolver struct {
	issueCount int
	issues     []*Issue
	users      map[string]*User
	config     ServerConfig
}

func NewResolver() *Resolver {
	resolver := new(Resolver)
	resolver.users = make(map[string]*User)
	return resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func filterIssue(a []*Issue, ID string) []*Issue {
	n := 0
	for _, x := range a {
		if x.ID != ID {
			a[n] = x
			n++
		}
	}
	return a[:n]
}

func findIssue(a []*Issue, ID string) *Issue {
	for _, x := range a {
		if x.ID == ID {
			return x
		}
	}

	return nil
}

func (issue *Issue) assign(user *User) (*Issue, error) {
	if issue.Assignee != nil {
		issue.unassign()
	}

	issue.Assignee = user

	user.Issues = append(user.Issues, issue)

	return issue, nil
}

func (issue *Issue) unassign() (*Issue, error) {
	if issue.Assignee != nil {
		issue.Assignee.Issues = filterIssue(issue.Assignee.Issues, issue.ID)

		issue.Assignee = nil
	}

	return issue, nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateIssues(ctx context.Context, input NewIssue) (*Issue, error) {
	r.issueCount++

	issue := &Issue{
		ID:          strconv.Itoa(r.issueCount),
		Description: input.Description,
	}

	if input.UserID != nil {
		if user, ok := r.users[*input.UserID]; ok {
			issue.assign(user)
		}
	}

	r.issues = append(r.issues, issue)
	return issue, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*User, error) {
	user := &User{
		ID:       *input.Email,
		Email:    input.Email,
		Name:     input.Name,
		Username: input.Username,
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *mutationResolver) AssignIssue(ctx context.Context, input AssignIssue) (*Issue, error) {
	user, ok := r.users[input.UserID]

	if !ok {
		return nil, errors.New("user not found")
	}

	var issue *Issue

	for _, searchIssue := range r.issues {
		if searchIssue.ID == input.IssueID {
			issue = searchIssue
		}
	}

	if issue == nil {
		return nil, errors.New("issue not found")
	}

	issue.assign(user)

	return issue, nil
}

func (r *mutationResolver) UnassignIssue(ctx context.Context, input string) (*Issue, error) {
	issue := findIssue(r.issues, input)

	if issue == nil {
		return nil, errors.New("issue not found")
	}

	issue.unassign()

	return issue, nil
}

func (r *mutationResolver) UpdateIssue(ctx context.Context, input *EditIssue) (*Issue, error) {
	issue := findIssue(r.issues, input.IssueID)

	if issue == nil {
		return nil, errors.New("issue not found")
	}

	issue.Description = input.Description

	return issue, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Issues(ctx context.Context) ([]*Issue, error) {
	return r.issues, nil
}
func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
	v := make([]*User, 0, len(r.users))
	for _, value := range r.users {
		v = append(v, value)
	}

	return v, nil
}
func (r *queryResolver) Config(ctx context.Context) (*ServerConfig, error) {
	return &r.config, nil
}
