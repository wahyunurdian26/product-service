-- +goose Up
-- +goose StatementBegin
INSERT INTO products (id, name, price, type, created_at, updated_at) VALUES 
('d3b07384-d113-4f52-870b-8dc0436a5784', 'Bayam', 5000, 'Sayuran', NOW(), NOW()),
('f2a58b41-e97d-419a-9e1b-b461877d9c12', 'Daging Ayam 1kg', 35000, 'Protein', NOW(), NOW()),
('8451927e-cfc2-4a73-9a3c-b1d6e1fa5cd3', 'Apel Fuji', 20000, 'Buah', NOW(), NOW()),
('c987a0f3-b26a-4d76-928f-7c18a59b3420', 'Chitato', 10000, 'Snack', NOW(), NOW())
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE id IN (
	'd3b07384-d113-4f52-870b-8dc0436a5784',
	'f2a58b41-e97d-419a-9e1b-b461877d9c12',
	'8451927e-cfc2-4a73-9a3c-b1d6e1fa5cd3',
	'c987a0f3-b26a-4d76-928f-7c18a59b3420'
);
-- +goose StatementEnd
