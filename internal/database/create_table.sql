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
    last_interaction TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contact_phone VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' NOT NULL,
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
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
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    direction VARCHAR(10) NOT NULL,
    read BOOLEAN DEFAULT false NOT NULL,
    status VARCHAR(20) DEFAULT 'sent' NOT NULL,
    CONSTRAINT fk_messages_conversation
        FOREIGN KEY (conversation_id)
        REFERENCES conversations (id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages (sent_at);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_sent_at 
ON messages (conversation_id, sent_at);

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

CREATE TABLE message_status_history (
    id BIGSERIAL PRIMARY KEY,
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE conversation_tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE conversation_tag_rel (
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    tag_id INT REFERENCES conversation_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (conversation_id, tag_id)
);

CREATE TABLE internal_notes (
    id SERIAL PRIMARY KEY,
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);