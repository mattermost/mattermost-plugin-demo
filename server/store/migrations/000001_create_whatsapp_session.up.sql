CREATE TABLE IF NOT EXISTS {{.prefix}}session (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    {{if .postgres}}created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),{{end}}
    {{if .mysql}}created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,{{end}}
    {{if .postgres}}closed_at TIMESTAMP WITH TIME ZONE{{end}}
    {{if .mysql}}closed_at TIMESTAMP NULL{{end}}
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};
