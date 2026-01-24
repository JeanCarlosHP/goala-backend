CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    
    -- RevenueCat integration
    revenuecat_user_id TEXT NOT NULL,
    revenuecat_original_transaction_id TEXT,
    
    -- Subscription state
    is_active BOOLEAN NOT NULL DEFAULT false,
    plan TEXT NOT NULL DEFAULT 'free',
    is_trial BOOLEAN NOT NULL DEFAULT false,
    
    -- Billing cycle
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    last_event_id TEXT,
    last_event_type TEXT,
    last_event_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_revenuecat_user_id ON subscriptions(revenuecat_user_id);
CREATE INDEX idx_subscriptions_is_active ON subscriptions(is_active);
CREATE INDEX idx_subscriptions_current_period_end ON subscriptions(current_period_end);
