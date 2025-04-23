-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id bigserial,
    name varchar(255) not null,
    email varchar(255) not null,
    bio text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
                             primary key (id)
    );

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_updated_at BEFORE UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd