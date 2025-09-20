CREATE TABLE permissions (
    id UUID NOT NULL PRIMARY KEY,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    UNIQUE(resource, action)
);

INSERT INTO permissions(id, resource, action) VALUES 
('750e8400-e29b-41d4-a716-446655440001', 'profile', 'read'),
('750e8400-e29b-41d4-a716-446655440002', 'profile', 'write'),
('750e8400-e29b-41d4-a716-446655440003', 'task', 'create'),
('750e8400-e29b-41d4-a716-446655440004', 'task', 'read'),
('750e8400-e29b-41d4-a716-446655440005', 'task', 'write'),
('750e8400-e29b-41d4-a716-446655440006', 'task', 'delete'),
('750e8400-e29b-41d4-a716-446655440007', 'user', 'read'),
('750e8400-e29b-41d4-a716-446655440008', 'user', 'write'),
('750e8400-e29b-41d4-a716-446655440009', 'user', 'delete');
