CREATE TABLE IF NOT EXISTS moderation_policies(id TEXT NOT NULL,version INT NOT NULL,name TEXT NOT NULL,status TEXT NOT NULL,updated_at TIMESTAMPTZ NOT NULL,PRIMARY KEY(id,version));
CREATE TABLE IF NOT EXISTS moderation_rules(id TEXT PRIMARY KEY,policy_id TEXT NOT NULL,policy_version INT NOT NULL,kind TEXT NOT NULL,pattern TEXT NOT NULL,risk TEXT NOT NULL,priority INT NOT NULL,decision TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS moderation_dictionaries(id TEXT PRIMARY KEY,name TEXT NOT NULL,version INT NOT NULL,terms JSONB NOT NULL);
CREATE TABLE IF NOT EXISTS moderation_samples(id TEXT PRIMARY KEY,content_hash TEXT NOT NULL,labels JSONB NOT NULL,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS moderation_decisions(sample_id TEXT PRIMARY KEY,policy_id TEXT NOT NULL,policy_version INT NOT NULL,decision TEXT NOT NULL,risk TEXT NOT NULL,hits JSONB NOT NULL,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS moderation_review_tasks(id TEXT PRIMARY KEY,sample_id TEXT NOT NULL,status TEXT NOT NULL,reviewer TEXT,lease_until TIMESTAMPTZ,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS moderation_appeals(id TEXT PRIMARY KEY,sample_id TEXT NOT NULL,status TEXT NOT NULL,reason TEXT,created_at TIMESTAMPTZ NOT NULL);
CREATE INDEX IF NOT EXISTS idx_review_claim ON moderation_review_tasks(status,lease_until,created_at);
