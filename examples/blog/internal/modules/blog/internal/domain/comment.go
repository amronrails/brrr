// Package domain holds the blog module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Comment is the comment entity.
type Comment struct {
	ID uuid.UUID
	Body string
	Approved bool
	PostID uuid.UUID
	AuthorID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateCommentInput carries the data needed to create a comment.
type CreateCommentInput struct {
	Body string
	Approved bool
	PostID uuid.UUID
	AuthorID uuid.UUID
}

// UpdateCommentInput carries the data needed to update a comment.
type UpdateCommentInput struct {
	Body string
	Approved bool
	PostID uuid.UUID
	AuthorID uuid.UUID
}

// ErrCommentNotFound is returned when a comment does not exist.
var ErrCommentNotFound = errors.New("comment not found")

// ErrCommentConflict is returned when a comment violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrCommentConflict = errors.New("comment already exists")
