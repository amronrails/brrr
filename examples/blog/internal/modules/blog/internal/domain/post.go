// Package domain holds the blog module's entities and value objects.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Post is the post entity.
type Post struct {
	ID uuid.UUID
	Title string
	Slug string
	Excerpt string
	Body string
	Published bool
	Views int32
	AuthorID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreatePostInput carries the data needed to create a post.
type CreatePostInput struct {
	Title string
	Slug string
	Excerpt string
	Body string
	Published bool
	Views int32
	AuthorID uuid.UUID
}

// UpdatePostInput carries the data needed to update a post.
type UpdatePostInput struct {
	Title string
	Slug string
	Excerpt string
	Body string
	Published bool
	Views int32
	AuthorID uuid.UUID
}

// ErrPostNotFound is returned when a post does not exist.
var ErrPostNotFound = errors.New("post not found")

// ErrPostConflict is returned when a post violates a unique
// constraint (e.g. a duplicate value in a unique column).
var ErrPostConflict = errors.New("post already exists")
