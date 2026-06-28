-- +goose Up
CREATE TABLE tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    status text NOT NULL,
    priority integer NOT NULL DEFAULT 0,
    due_date date NOT NULL DEFAULT CURRENT_DATE,
    done boolean NOT NULL DEFAULT false,
    project_id uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    assignee_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS tasks;
