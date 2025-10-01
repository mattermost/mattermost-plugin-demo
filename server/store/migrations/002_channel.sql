CREATE TABLE IF NOT EXISTS whatsapp_plugin_channel (
   id VARCHAR(26) PRIMARY KEY,
   channel_id VARCHAR(26) NOT NULL,
   phone_number VARCHAR(20) NOT NULL,
   phone_number_id VARCHAR(20) NOT NULL
);
