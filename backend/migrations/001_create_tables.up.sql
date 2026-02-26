-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_role CHECK (role IN ('admin', 'operator', 'viewer'))
);

-- Create incidents table
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    CONSTRAINT check_severity CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT check_status CHECK (status IN ('open', 'assigned', 'in_progress', 'resolved', 'closed'))
);

-- Create indexes
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_assigned_to ON incidents(assigned_to);
CREATE INDEX idx_incidents_created_by ON incidents(created_by);
CREATE INDEX idx_incidents_created_at ON incidents(created_at DESC);
CREATE INDEX idx_users_email ON users(email);

-- Insert default admin user for development
INSERT INTO users (id, email, name, role, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin@ops-incident-hub.com',
    'Admin User',
    'admin',
    NOW()
) ON CONFLICT (id) DO NOTHING;
