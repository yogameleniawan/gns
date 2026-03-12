-- Seed default roles
INSERT INTO roles (name, description, is_system) VALUES
('Super Admin', 'Full system access with all permissions', TRUE),
('Admin', 'Administrative access with most permissions', TRUE),
('User', 'Standard user access', TRUE)
ON CONFLICT (name) DO NOTHING;

-- Seed default permissions
INSERT INTO permissions (name, module, description) VALUES
-- Auth permissions
('auth.login', 'auth', 'Login to the system'),
('auth.register', 'auth', 'Register new account'),

-- User management permissions
('users.read', 'users', 'View users'),
('users.create', 'users', 'Create new users'),
('users.update', 'users', 'Update user information'),
('users.delete', 'users', 'Delete users'),
('users.manage_status', 'users', 'Activate/deactivate users'),

-- Role management permissions
('roles.read', 'roles', 'View roles'),
('roles.create', 'roles', 'Create new roles'),
('roles.update', 'roles', 'Update role information'),
('roles.delete', 'roles', 'Delete roles'),
('roles.assign', 'roles', 'Assign roles to users'),

-- Permission management permissions
('permissions.read', 'permissions', 'View permissions'),
('permissions.create', 'permissions', 'Create new permissions'),
('permissions.update', 'permissions', 'Update permission information'),
('permissions.delete', 'permissions', 'Delete permissions'),
('permissions.assign', 'permissions', 'Assign permissions to roles'),

-- Broker permissions
('brokers.read', 'brokers', 'View brokers'),
('brokers.create', 'brokers', 'Create new brokers'),
('brokers.update', 'brokers', 'Update broker information'),
('brokers.delete', 'brokers', 'Delete brokers'),

-- Broker trades permissions
('broker_trades.read', 'broker_trades', 'View broker trades'),
('broker_trades.scrape', 'broker_trades', 'Trigger scraping'),

-- Dashboard permissions
('dashboard.view', 'dashboard', 'View dashboard')
ON CONFLICT (name) DO NOTHING;

-- Assign all permissions to Super Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Super Admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign most permissions to Admin (except super admin specific ones)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
AND p.name NOT IN ('permissions.delete', 'roles.delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign basic permissions to User
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'User'
AND p.name IN (
    'auth.login',
    'dashboard.view',
    'brokers.read',
    'broker_trades.read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Set up module access for Super Admin
INSERT INTO module_access (role_id, module_name, can_view, can_create, can_edit, can_delete)
SELECT r.id, module_name, TRUE, TRUE, TRUE, TRUE
FROM roles r
CROSS JOIN (
    VALUES 
        ('users'),
        ('roles'),
        ('permissions'),
        ('brokers'),
        ('broker_trades'),
        ('dashboard')
) AS modules(module_name)
WHERE r.name = 'Super Admin'
ON CONFLICT (role_id, module_name) DO NOTHING;

-- Set up module access for Admin
INSERT INTO module_access (role_id, module_name, can_view, can_create, can_edit, can_delete)
SELECT r.id, module_name, can_view, can_create, can_edit, can_delete
FROM roles r
CROSS JOIN (
    VALUES 
        ('users', TRUE, TRUE, TRUE, FALSE),
        ('roles', TRUE, TRUE, TRUE, FALSE),
        ('permissions', TRUE, FALSE, FALSE, FALSE),
        ('brokers', TRUE, TRUE, TRUE, TRUE),
        ('broker_trades', TRUE, TRUE, TRUE, TRUE),
        ('dashboard', TRUE, FALSE, FALSE, FALSE)
) AS modules(module_name, can_view, can_create, can_edit, can_delete)
WHERE r.name = 'Admin'
ON CONFLICT (role_id, module_name) DO NOTHING;

-- Set up module access for User
INSERT INTO module_access (role_id, module_name, can_view, can_create, can_edit, can_delete)
SELECT r.id, module_name, can_view, can_create, can_edit, can_delete
FROM roles r
CROSS JOIN (
    VALUES 
        ('dashboard', TRUE, FALSE, FALSE, FALSE),
        ('brokers', TRUE, FALSE, FALSE, FALSE),
        ('broker_trades', TRUE, FALSE, FALSE, FALSE)
) AS modules(module_name, can_view, can_create, can_edit, can_delete)
WHERE r.name = 'User'
ON CONFLICT (role_id, module_name) DO NOTHING;
