package goadmin

import (
	"context"
	"net/http"

	"github.com/zhenyangze/goadmin/form"
	"github.com/zhenyangze/goadmin/grid"
	"github.com/zhenyangze/goadmin/show"
	"github.com/zhenyangze/goadmin/tree"
)

// Resource describes one admin module with Dcat-like page builders.
type Resource struct {
	Name              string
	Path              string
	Title             string
	Description       string
	Icon              string
	Permission        string
	EmptyText         string
	CapabilityTags    []string
	VerificationSteps []string
	Repository        Repository
	BuildGrid         func(*grid.Builder)
	BuildForm         func(*form.Builder)
	BuildShow         func(*show.Builder)
	BuildTree         func(*tree.Builder)
}

// HasTree reports whether the resource exposes a tree page.
func (r Resource) HasTree() bool {
	return r.BuildTree != nil
}

// DashboardPage is a read-only shell page rendered with the dashboard template.
type DashboardPage struct {
	Path        string
	Title       string
	Description string
	Permission  string
	Build       func(context.Context, *http.Request, *Identity) (DashboardData, error)
}

// Route is a custom non-resource endpoint mounted under the admin prefix.
type Route struct {
	Path       string
	Methods    []string
	Permission string
	Handler    func(http.ResponseWriter, *http.Request, *Identity) error
}
