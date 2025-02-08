-- Table for storing users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,  -- Unique identifier/name for the user
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing accounts
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    fk_users INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key to users table
    name VARCHAR(255) NOT NULL,  -- Unique Name of the account
    xpub TEXT NOT NULL,
    acc_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fk_users, name)  -- Ensuring a user can’t create multiple accounts with the same name
);

-- Table for storing addresses associated with accounts
CREATE TABLE account_addresses (
    id SERIAL PRIMARY KEY,
    fk_accounts INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,  -- Foreign key to accounts table
    address VARCHAR(255) NOT NULL,  -- Bitcoin address
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fk_accounts, address)  -- Ensuring unique addresses per account
);

-- Triggers to update 'updated_at' fields automatically
CREATE OR REPLACE FUNCTION update_timestamp() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_account_timestamp
    BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_account_address_timestamp
    BEFORE UPDATE ON account_addresses
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
