CREATE TABLE outbox (
    id uuid PRIMARY KEY,
    topic varchar(255) NOT NULL,
    message jsonb NOT NULL,
    state varchar(50) NOT NULL DEFAULT 'pending', -- e.g., pending, processed
    created_at timestamptz NOT NULL DEFAULT now(),
    processed_at timestamptz
);

