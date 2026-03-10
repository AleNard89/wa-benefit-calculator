DROP INDEX IF EXISTS idx_chat_document_chunks_file;
DROP INDEX IF EXISTS idx_chat_document_chunks_company;
DROP TABLE IF EXISTS chat_document_chunks CASCADE;

DROP INDEX IF EXISTS idx_chat_messages_conversation;
DROP TABLE IF EXISTS chat_messages CASCADE;

DROP INDEX IF EXISTS idx_chat_conversations_user_company;
DROP TABLE IF EXISTS chat_conversations CASCADE;

DROP EXTENSION IF EXISTS vector;
