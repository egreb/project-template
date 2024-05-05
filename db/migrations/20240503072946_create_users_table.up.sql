CREATE TABLE IF NOT EXISTS users (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	username VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	salt bytea NOT NULL, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE(username)
);
