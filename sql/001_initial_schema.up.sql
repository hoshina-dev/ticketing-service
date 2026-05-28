CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE IF NOT EXISTS ticket_status AS ENUM (
    'requested',
    'pending',
    'experimenting',
    'finalizing'
    -- Questions for phases not in Nathan's seq diagram:
        -- Do we have a phase after `finalizing`. It sounds like there should be a completed or something
        -- Will we have cancelled?
);

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Client that requested the experiment
    user_id         UUID NOT NULL,
    organization_id UUID NOT NULL,

    status ticket_status NOT NULL DEFAULT 'request',

    -- Maybe we should include delivery method information? 
    -- As from what I remember they could either have the sample shipped or deliver it in person
    
    -- And maybe sample_type_id ? For like check if what deliver is correct.
    -- Say on website client requested an analysis on a piece of Coal. But sent a Banana

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
    -- Maybe also timestamps for phase changes? like sample_received_at, completed_at, ...
);

-- Join table for which experiment templates the client requested for this ticket.
CREATE TABLE IF NOT EXISTS ticket_experiment_templates (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id               UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    experiment_template_id  UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ

    UNIQUE (ticket_id, experiment_template_id) -- I forgot what you did for this Nathan, help
);

-- Index
CREATE INDEX IF NOT EXISTS idx_tickets_user ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_organization ON tickets(organization_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

CREATE INDEX IF NOT EXISTS idx_tet_ticket ON ticket_experiment_templates(ticket_id);
CREATE INDEX IF NOT EXISTS idx_tet_experiment_template ON ticket_experiment_templates(experiment_template_id);