INSERT INTO products (id, name, price, type, created_at, updated_at) VALUES 
('11111111-1111-1111-1111-111111111111', 'Bayam', 5000, 'Sayuran', NOW(), NOW()),
('22222222-2222-2222-2222-222222222222', 'Daging Ayam 1kg', 35000, 'Protein', NOW(), NOW()),
('33333333-3333-3333-3333-333333333333', 'Apel Fuji', 20000, 'Buah', NOW(), NOW()),
('44444444-4444-4444-4444-444444444444', 'Chitato', 10000, 'Snack', NOW(), NOW())
ON CONFLICT DO NOTHING;
