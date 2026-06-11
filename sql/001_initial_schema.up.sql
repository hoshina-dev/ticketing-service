CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE ticket_status AS ENUM (
        'REQUESTED',
        'PENDING',
        'EXPERIMENTING',
        'FINALIZING',
        'CLOSED'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Client that requested the experiment
    user_id         UUID NOT NULL,
    organization_id UUID NOT NULL,

    status          ticket_status NOT NULL DEFAULT 'REQUESTED',
    closed_reason   TEXT,
    
    -- Add sample_type_id column? To check if the delivered sample matches the request.
    -- Say a client requested an analysis of a piece of Coal on the website, but instead sent a Banana.

    sample_received_at      TIMESTAMPTZ,
    experiment_started_at   TIMESTAMPTZ,
    results_submitted_at    TIMESTAMPTZ,
    closed_at               TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Join table for which experiment templates the client requested for this ticket.
CREATE TABLE IF NOT EXISTS ticket_experiment_templates (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Used for creating an entry in experiments table
    ticket_id               UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    experiment_template_id  UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index
-- One-to-one: a ticket may have at most one active experiment template (one experiment form per ticket).
CREATE UNIQUE INDEX IF NOT EXISTS idx_tet_one_per_ticket ON ticket_experiment_templates(ticket_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tickets_user ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_organization ON tickets(organization_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

CREATE INDEX IF NOT EXISTS idx_tet_ticket ON ticket_experiment_templates(ticket_id);
CREATE INDEX IF NOT EXISTS idx_tet_experiment_template ON ticket_experiment_templates(experiment_template_id);