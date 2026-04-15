Create EXTENSION if not exists "uuid-ossp";

create table if not exists users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email varchar(255) unique not null,
    password varchar(255) not null,
    created_at timestamp with time zone default now() not null
);

create index if not exists idx_users_email on users (email);

create table if not exists contacts (
    phone varchar(20) PRIMARY key,
    name varchar(255) not null,
    last_interaction timestamp with time zone
);

create table if not exists conversation (
    id uuid primary key Default uuid_generate_v4(),
    contact_phone varchar(20) not null,
    status varchar(20) default 'active' not null,
    created_at timestamp with time zone default now() not null,
    Constraint fk_conversations_contact
        foreign key (contact_phone)
        references contacts (phone)
        on delete cascade
);

create index if not exists idx_conversation_contact_phone on conversation (contact_phone);

create table if not exists messages (
    id varchar(255) primary key,
    conversation_id UUID not null,
    from_phone varchar(20) not null,
    type varchar(20) not null,
    body text,
    media_id varchar(255),
    media_url text,
    timestamp Timestamp with time zone not null,
    direction varchar(10) not null,
    read Boolean default false not null,
    status varchar(20) default 'sent' not null,
    Constraint fk_messages_conversation
        foreign key (conversation_id)
        references conversation (id)
        on delete cascade
);

create index if not exists idx_messages_conversation_id on messages (conversation_id);
create index if not exists idx_messages_timestamp on messages (timestamp);

create table if not Exists media_file (
    id UUID PRIMARY key default uuid_generate_v4(),
    whatsapp_media_id varchar(255) unique not null,
    file_path text not null,
    mime_type varchar(100) not null,
    file_size BIGINT,

    duration_seconds int,
    is_voice_note Boolean default false,

    created_at timestamp with time zone default now() not null
);

created index if not exists idx_media_files_whatsapp_media_id on media_files (whatsapp_media_id);
created index if not exists idx_media_files_mime_type on media_file (mime_type)

create table if not exists webhook_events(
    id BIGSERIAL primary key,
    payload JSONB not null
    processed Boolean default false not null,
    error_log text,
    created_at timestamp with time zone default now() not null
);

create index if not exists idx_webhook_events_processed on webhook_events (processed)
create index if not exists idx_webhook_events_created_at on webhook_events (created_at)