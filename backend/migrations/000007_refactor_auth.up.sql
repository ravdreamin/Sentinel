-- Remove NOT NULL constraint from password_hash
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- table to track login methods
CREATE TABLE IF NOT EXISTS user_identities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'local', 'google'
    provider_id VARCHAR(255),       -- OAuth subject ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, provider)
);
