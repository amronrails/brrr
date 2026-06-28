-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    slug text NOT NULL UNIQUE,
    excerpt text NOT NULL DEFAULT '',
    body text NOT NULL DEFAULT '',
    published boolean NOT NULL DEFAULT false,
    views integer NOT NULL DEFAULT 0,
    author_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS posts;
