CREATE TABLE IF NOT EXISTS invitation_tbl (
    id TEXT PRIMARY KEY,

    invitor_id TEXT NOT NULL REFERENCES user_tbl(id) ON DELETE CASCADE,
    invitee_qq TEXT NOT NULL,

    invitation_code VARCHAR(64) NOT NULL,

    assign_translator BOOLEAN DEFAULT FALSE,
    assign_proofreader BOOLEAN DEFAULT FALSE,
    assign_typesetter BOOLEAN DEFAULT FALSE,
    assign_redrawer BOOLEAN DEFAULT FALSE,
    assign_reviewer BOOLEAN DEFAULT FALSE,
    assign_uploader BOOLEAN DEFAULT FALSE,

    pending BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_invitation_invitee_qq ON invitation_tbl (invitee_qq);

CREATE INDEX IF NOT EXISTS idx_invitation_updated_at_pending_true 
    ON invitation_tbl (updated_at) 
    WHERE pending = TRUE;