package chat

import "time"

const (
	ConversationsTable = "chat_conversations"
	MessagesTable      = "chat_messages"
	DocumentChunksTable = "chat_document_chunks"
)

type Conversation struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"userId"`
	CompanyID int       `db:"company_id" json:"companyId"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type Message struct {
	ID             int       `db:"id" json:"id"`
	ConversationID int       `db:"conversation_id" json:"conversationId"`
	Role           string    `db:"role" json:"role"` // "user" or "assistant"
	Content        string    `db:"content" json:"content"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

type DocumentChunk struct {
	ID         int       `db:"id" json:"id"`
	CompanyID  int       `db:"company_id" json:"companyId"`
	FilePath   string    `db:"file_path" json:"filePath"`
	FileName   string    `db:"file_name" json:"fileName"`
	ChunkIndex int       `db:"chunk_index" json:"chunkIndex"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}
