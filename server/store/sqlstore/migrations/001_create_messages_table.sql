-- Migration: Create messages table
-- Description: Creates the messages table to store plugin messages

CREATE TABLE IF NOT EXISTS whatsapp_plugin_session (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS whatsapp_plugin_channel (
   id VARCHAR(26) PRIMARY KEY,
   channel_id VARCHAR(26) NOT NULL,
   phone_number VARCHAR(20) NOT NULL,
   phone_number_id VARCHAR(20) NOT NULL
)
