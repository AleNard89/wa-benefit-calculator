-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Chat conversations
CREATE TABLE IF NOT EXISTS chat_conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chat_conversations_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES auth_users (id) ON DELETE CASCADE,
    CONSTRAINT chat_conversations_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE
);

CREATE INDEX idx_chat_conversations_user_company ON chat_conversations (user_id, company_id);

-- Chat messages
CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL,
    role VARCHAR(20) NOT NULL, -- 'user' or 'assistant'
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chat_messages_conversation_id_fkey FOREIGN KEY (conversation_id)
        REFERENCES chat_conversations (id) ON DELETE CASCADE
);

CREATE INDEX idx_chat_messages_conversation ON chat_messages (conversation_id);

-- Document chunks for RAG
CREATE TABLE IF NOT EXISTS chat_document_chunks (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chat_document_chunks_company_id_fkey FOREIGN KEY (company_id)
        REFERENCES orgs_companies (id) ON DELETE CASCADE
);

CREATE INDEX idx_chat_document_chunks_company ON chat_document_chunks (company_id);
CREATE INDEX idx_chat_document_chunks_file ON chat_document_chunks (company_id, file_path);
