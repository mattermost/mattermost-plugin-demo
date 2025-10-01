CREATE TABLE IF NOT EXISTS {{.prefix}}channel (
   id VARCHAR(26) PRIMARY KEY,
   channel_id VARCHAR(26) NOT NULL,
   phone_number VARCHAR(20) NOT NULL,
   phone_number_id VARCHAR(20) NOT NULL
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
