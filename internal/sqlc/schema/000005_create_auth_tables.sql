-- +goose Up
-- +goose StatementBegin
-- Refresh token denylist for JWT security
CREATE TABLE IF NOT EXISTS refresh_token_denylist (
    jti VARCHAR(255) PRIMARY KEY, -- JWT ID claim
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    denylisted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason VARCHAR(100) DEFAULT 'logout' -- logout, revoked, expired, etc.
);

-- User sessions table for session management
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Scheduled deletions table for GDPR compliance
CREATE TABLE IF NOT EXISTS scheduled_deletions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deletion_type VARCHAR(50) NOT NULL CHECK (deletion_type IN ('account', 'data', 'partial')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'cancelled')),
    completed_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_refresh_denylist_user ON refresh_token_denylist(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_denylist_expires ON refresh_token_denylist(expires_at);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_active ON user_sessions(user_id, is_active) WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_scheduled_deletions_user ON scheduled_deletions(user_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_deletions_scheduled ON scheduled_deletions(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_deletions_status ON scheduled_deletions(status);

-- Function to clean up expired tokens automatically
CREATE OR REPLACE FUNCTION cleanup_expired_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM refresh_token_denylist
    WHERE expires_at < NOW() - INTERVAL '7 days';

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    -- Also clean up expired sessions
    DELETE FROM user_sessions
    WHERE expires_at < NOW();

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS cleanup_expired_tokens();
DROP INDEX IF EXISTS idx_scheduled_deletions_status;
DROP INDEX IF EXISTS idx_scheduled_deletions_scheduled;
DROP INDEX IF EXISTS idx_scheduled_deletions_user;
DROP INDEX IF EXISTS idx_user_sessions_active;
DROP INDEX IF EXISTS idx_user_sessions_expires;
DROP INDEX IF EXISTS idx_user_sessions_user;
DROP INDEX IF EXISTS idx_refresh_denylist_expires;
DROP INDEX IF EXISTS idx_refresh_denylist_user;
DROP TABLE IF EXISTS scheduled_deletions CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS refresh_token_denylist CASCADE;
-- +goose StatementEnd
