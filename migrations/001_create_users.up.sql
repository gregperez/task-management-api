CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('Administrador', 'Ejecutor', 'Auditor')),
    is_temporary_password BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Crear usuario administrador por defecto
-- Contraseña: Admin123!
INSERT INTO users (id, username, email, password, role, is_temporary_password, created_at, updated_at)
VALUES (
    'admin-001',
    'admin',
    'admin@taskmanager.com',
    '$2a$10$xL3JsIdavhFIMVrft7NmmOLCr6zF75Yx7j.M0huvtM1H5pkdB2goK',
    'Administrador',
    FALSE,
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;