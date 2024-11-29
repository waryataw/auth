-- +goose Up
-- +goose StatementBegin
CREATE TABLE accessible_roles
(
    id   SERIAL PRIMARY KEY,
    path VARCHAR  NOT NULL,
    role SMALLINT NOT NULL,
    CONSTRAINT unique_path_role UNIQUE (path, role)
);

INSERT INTO accessible_roles (path, role)
VALUES ('/chat_server_v1.ChatServerService/CreateChat', 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accessible_roles;
-- +goose StatementEnd
