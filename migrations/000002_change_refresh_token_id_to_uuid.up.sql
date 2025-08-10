-- Add new UUID column
ALTER TABLE refresh_tokens ADD COLUMN uuid_id UUID DEFAULT gen_random_uuid();

-- Update all existing records to have UUID values
UPDATE refresh_tokens SET uuid_id = gen_random_uuid() WHERE uuid_id IS NULL;

-- Make uuid_id NOT NULL
ALTER TABLE refresh_tokens ALTER COLUMN uuid_id SET NOT NULL;

-- Drop the old primary key constraint
ALTER TABLE refresh_tokens DROP CONSTRAINT refresh_tokens_pkey;

-- Drop the old id column
ALTER TABLE refresh_tokens DROP COLUMN id;

-- Rename uuid_id to id
ALTER TABLE refresh_tokens RENAME COLUMN uuid_id TO id;

-- Add new primary key constraint
ALTER TABLE refresh_tokens ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);
