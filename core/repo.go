package core

import "context"

//go:generate mockgen -package mock -destination ../mock/repo_mock.go . RepoService,RepoStore

// Repo defined a repository structure
type Repo struct {
	ID        uint
	URL       string
	ReportID  string
	NameSpace string
	Name      string
	SCM       SCMProvider
}

// RepoService provides operations with SCM
type RepoService interface {
	NewReportID(repo *Repo) string
	// List repositories from SCM context
	List(ctx context.Context, scm SCMProvider, user *User) ([]*Repo, error)
}

// RepoStore repository in storage
type RepoStore interface {
	Create(repo *Repo) error
	Update(repo *Repo) error
	Find(repo *Repo) (*Repo, error)
	Finds(urls ...string) ([]*Repo, error)
}