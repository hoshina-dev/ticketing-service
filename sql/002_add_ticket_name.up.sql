-- Optional human-friendly name; auto-generated when not provided on creation.
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS name TEXT;
