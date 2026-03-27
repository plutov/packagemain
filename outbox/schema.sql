CREATE TABLE orders (
    id uuid PRIMARY KEY,
    product varchar(255) NOT NULL,
    quantity int NOT NULL
);

CREATE TABLE outbox (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    topic varchar(255) NOT NULL,
    message jsonb NOT NULL,
    state varchar(50) NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT now(),
    processed_at timestamptz
);
