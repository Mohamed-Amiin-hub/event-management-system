package db

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/EVENT-MANAGEMENT-SYSTEM/pkg/config"
	_ "github.com/lib/pq"
)

// ConnectDB initializes and returns a PostgreSQL database connection.
func ConnectDB(cfg *config.DBConfig) (*sql.DB, error) {
	connstr := cfg.ConnectionString()
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Ping the database to ensure the connection is successful
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	log.Println("Successfully connected to the database")

	return db, nil
}

// CreateTables ensures that the required tables are created.
func CreateTables(db *sql.DB) error {
	userTable := `CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	first_name VARCHAR(255),
	last_name VARCHAR(255),
	is_active BOOLEAN DEFAULT FALSE,
	last_login TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
	);`

	//
	eventTable := `CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			location VARCHAR(255),
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP,
			capacity INTEGER NOT NULL CHECK (capacity >= 0),
			is_public BOOLEAN DEFAULT true,
			status TEXT NOT NULL,
			organizer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
			);`

	// Create tokens table
	tokenTable := `CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Corrected: No 'UPDATE CURRENT_TIMESTAMP' here
    deleted_at TIMESTAMP NULL
);
`

	// Create roles table
	roleTable := `CREATE TABLE IF NOT EXISTS roles (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL
	);`

	// Create permissions table
	permissionTable := `CREATE TABLE IF NOT EXISTS permissions (
	id  SERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL 
	);`

	// Create user_roles table
	createRoleTable := `CREATE TABLE IF NOT EXISTS user_roles (
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	role_id INT REFERENCES roles(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, role_id)
	);`

	// Create user_permissions table
	userPermissionTable := `CREATE TABLE IF NOT EXISTS user_permissions (
	user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
	PRIMARY KEY(user_id, permission_id)
	);`

	// Execute the table creation queries
	query := []string{userTable, tokenTable, roleTable, permissionTable, createRoleTable, userPermissionTable, eventTable}
	for _, query := range query {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
		log.Println("Successfully created all tables")
	}
	return nil
}
