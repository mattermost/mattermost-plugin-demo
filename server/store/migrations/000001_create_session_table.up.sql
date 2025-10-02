CREATE TABLE IF NOT EXISTS {{.prefix}}session (
    id VARCHAR(26) UNIQUE NOT NULL,
    user_id VARCHAR(26) NOT NULL,
    create_at BIGINT NOT NULL,
    PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
