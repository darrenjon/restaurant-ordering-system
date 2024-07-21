package database

import (
	"fmt"
)

// RunMigrations executes all database migrations
func (m *Manager) RunMigrations() error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "Create users table",
			sql: `CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(50) UNIQUE NOT NULL,
				password VARCHAR(100) NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				role VARCHAR(20) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create categories table",
			sql: `CREATE TABLE IF NOT EXISTS categories (
				id SERIAL PRIMARY KEY,
				name VARCHAR(50) UNIQUE NOT NULL,
				display_order INT NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create menu_items table",
			sql: `CREATE TABLE IF NOT EXISTS menu_items (
				id SERIAL PRIMARY KEY,
				category_id INT REFERENCES categories(id),
				name VARCHAR(100) NOT NULL,
				description TEXT,
				price DECIMAL(10, 2) NOT NULL,
				image_url VARCHAR(255),
				is_available BOOLEAN NOT NULL DEFAULT TRUE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create add_ons table",
			sql: `CREATE TABLE IF NOT EXISTS add_ons (
				id SERIAL PRIMARY KEY,
				name VARCHAR(50) NOT NULL,
				price DECIMAL(10, 2) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create menu_item_add_ons table",
			sql: `CREATE TABLE IF NOT EXISTS menu_item_add_ons (
				id SERIAL PRIMARY KEY,
				menu_item_id INT REFERENCES menu_items(id),
				add_on_id INT REFERENCES add_ons(id),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(menu_item_id, add_on_id)
			)`,
		},
		{
			name: "Create orders table",
			sql: `CREATE TABLE IF NOT EXISTS orders (
				id SERIAL PRIMARY KEY,
				table_number VARCHAR(10) NOT NULL,
				status VARCHAR(20) NOT NULL,
				total_amount DECIMAL(10, 2) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create order_items table",
			sql: `CREATE TABLE IF NOT EXISTS order_items (
				id SERIAL PRIMARY KEY,
				order_id INT REFERENCES orders(id),
				menu_item_id INT REFERENCES menu_items(id),
				quantity INT NOT NULL,
				unit_price DECIMAL(10, 2) NOT NULL,
				subtotal DECIMAL(10, 2) NOT NULL,
				special_instructions TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create order_item_add_ons table",
			sql: `CREATE TABLE IF NOT EXISTS order_item_add_ons (
				id SERIAL PRIMARY KEY,
				order_item_id INT REFERENCES order_items(id),
				add_on_id INT REFERENCES add_ons(id),
				price DECIMAL(10, 2) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
		{
			name: "Create restaurant_info table",
			sql: `CREATE TABLE IF NOT EXISTS restaurant_info (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				description TEXT,
				address TEXT,
				phone VARCHAR(20),
				email VARCHAR(100),
				logo_url VARCHAR(255),
				banner_url VARCHAR(255),
				opening_hours JSONB,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)`,
		},
	}

	for _, migration := range migrations {
		if err := m.executeMigration(migration.name, migration.sql); err != nil {
			return fmt.Errorf("migration '%s' failed: %w", migration.name, err)
		}
	}

	return nil
}

func (m *Manager) executeMigration(name, sql string) error {
	fmt.Printf("Running migration: %s\n", name)

	_, err := m.db.Exec(sql)
	if err != nil {
		return err
	}

	fmt.Printf("Migration completed: %s\n", name)
	return nil
}
