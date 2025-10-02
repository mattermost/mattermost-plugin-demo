CREATE TABLE IF NOT EXISTS {{.prefix}}channel (
    id VARCHAR(26) UNIQUE NOT NULL,
    channel_id VARCHAR(26) NOT NULL,
    PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
