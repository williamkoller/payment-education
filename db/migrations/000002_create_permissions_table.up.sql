CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    modules TEXT[] NOT NULL,
    actions TEXT[] NOT NULL,
    level TEXT NOT NULL CHECK (level IN ('allowed', 'restricted', 'denied')),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);