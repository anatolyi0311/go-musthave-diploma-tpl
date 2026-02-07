package migrations_accrual

const (
	createEnum = `
		CREATE TYPE discount_type AS ENUM ('PERCENT', 'FIXED');
	`

	createTableCatalog = `
		CREATE TABLE IF NOT EXISTS catalog
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			description VARCHAR(255) NOT NULL UNIQUE,
			reward INT NOT NULL,
			reward_type discount_type NOT NULL,
			CONSTRAINT reward_percent_check 
				CHECK (
					(reward_type = 'PERCENT' AND reward <= 100 AND reward > 0) OR
					(reward_type = 'FIXED' AND reward > 0)
				)
		);
	`
	createTableOrders = `
		CREATE TABLE IF NOT EXISTS orders
		(
			id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			description VARCHAR(255) NOT NULL UNIQUE,
			price NUMERIC(10, 2),
			order_number VARCHAR(255) NOT NULL UNIQUE
		);
	`
	createTableAccrual = `
		CREATE TABLE IF NOT EXISTS accrual_orders
		(
		id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		order_number VARCHAR(255) NOT NULL UNIQUE,
		status VARCHAR(255) NOT NULL,
		accrual NUMERIC(10, 2)
		);
	`
	dropTableCatalog = `
		DROP TABLE IF EXISTS catalog;
	`
	dropTableOrders = `
		DROP TABLE IF EXISTS orders;
	`
	dropTableAccrual = `
		DROP TABLE IF EXISTS accrual_orders;
	`

	dropEnum = `
		DROP TYPE IF EXISTS discount_type;
	`

)