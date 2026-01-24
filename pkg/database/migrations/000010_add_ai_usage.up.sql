CREATE TABLE IF NOT EXISTS ai_usage (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Feature tracking
    feature TEXT NOT NULL,
    
    -- Usage counters
    usage_count INTEGER NOT NULL DEFAULT 0,
    quota INTEGER NOT NULL,
    
    -- Reset tracking
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, feature, period_start)
);

CREATE INDEX idx_ai_usage_user_id ON ai_usage(user_id);
CREATE INDEX idx_ai_usage_feature ON ai_usage(feature);
CREATE INDEX idx_ai_usage_period_end ON ai_usage(period_end);
CREATE INDEX idx_ai_usage_user_feature ON ai_usage(user_id, feature);
