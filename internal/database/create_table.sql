CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS contacts (
    phone VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    last_interaction TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contact_phone VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' NOT NULL,
    last_message_at timestamp with time zone;
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    CONSTRAINT fk_conversations_contact
        FOREIGN KEY (contact_phone)
        REFERENCES contacts (phone)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_conversations_contact_phone ON conversations (contact_phone);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL,
    from_phone VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL,
    body TEXT,
    media_id VARCHAR(255),
    media_url TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    direction VARCHAR(10) NOT NULL,
    read BOOLEAN DEFAULT false NOT NULL,
    status VARCHAR(20) DEFAULT 'sent' NOT NULL,
    CONSTRAINT fk_messages_conversation
        FOREIGN KEY (conversation_id)
        REFERENCES conversations (id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages (timestamp);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_timestamp 
ON messages (conversation_id, timestamp);

CREATE TABLE IF NOT EXISTS media_file (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    whatsapp_media_id VARCHAR(255) UNIQUE NOT NULL,
    file_path TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size BIGINT,
    duration_seconds INT,
    is_voice_note BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_media_file_whatsapp_media_id ON media_file (whatsapp_media_id);
CREATE INDEX IF NOT EXISTS idx_media_file_mime_type ON media_file (mime_type);

CREATE TABLE IF NOT EXISTS webhook_events (
    id BIGSERIAL PRIMARY KEY,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT false NOT NULL,
    error_log TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_processed ON webhook_events (processed);
CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at ON webhook_events (created_at);

create table message_status_history (
    id BIGSERIAL PRIMARY key,
    message_id UUID,
    status VARCHAR(20),
    created_at TIMESTAMP with time zone DEFAULT now()
)

create table sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id BIGINT,
    token Text,
    expires_at timestamp
)

create table conversation_tags (
    id SERIAL PRIMARY key,
    name VARCHAR(50)
);

create table conversation_tag_rel (
    conversation_id UUID,
    tag_id int
)

create table internal_notes (
    id SERIAL PRIMARY key,
    conversation_id UUID,
    content TEXT,
    created_at timestamp
)