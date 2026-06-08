INSERT INTO users (email, password_hash, full_name, role)
VALUES (
    'admin@system.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: admin123456
    'System Admin',
    'admin'
)
ON CONFLICT (email) DO NOTHING;