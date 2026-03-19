package chat

import (
	"context"
	"fmt"
	"strings"

	"orbita-api/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type ConversationService struct{}

func (s *ConversationService) GetByUser(userID, companyID int) ([]Conversation, error) {
	db := db.DB()
	conversations := make([]Conversation, 0)
	stm := db.Builder.Select("*").From(ConversationsTable).
		Where(sq.Eq{"user_id": userID, "company_id": companyID}).
		OrderBy("updated_at DESC")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), db.C, &conversations, sql, args...)
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (s *ConversationService) GetByID(id int) (*Conversation, error) {
	db := db.DB()
	conv := &Conversation{}
	stm := db.Builder.Select("*").From(ConversationsTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(context.Background(), db.C, conv, sql, args...)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *ConversationService) Create(userID, companyID int, title string) (*Conversation, error) {
	db := db.DB()
	conv := &Conversation{}
	stm := db.Builder.Insert(ConversationsTable).
		Columns("user_id", "company_id", "title").
		Values(userID, companyID, title).
		Suffix("RETURNING id, user_id, company_id, title, created_at, updated_at")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(context.Background(), db.C, conv, sql, args...)
	return conv, err
}

func (s *ConversationService) UpdateTitle(id int, title string) error {
	db := db.DB()
	stm := db.Builder.Update(ConversationsTable).
		Set("title", title).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *ConversationService) Touch(id int) error {
	db := db.DB()
	stm := db.Builder.Update(ConversationsTable).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *ConversationService) Delete(id int) error {
	db := db.DB()
	stm := db.Builder.Delete(ConversationsTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

// MessageService

type MessageService struct{}

func (s *MessageService) GetByConversation(conversationID int) ([]Message, error) {
	db := db.DB()
	messages := make([]Message, 0)
	stm := db.Builder.Select("*").From(MessagesTable).
		Where(sq.Eq{"conversation_id": conversationID}).
		OrderBy("created_at ASC")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), db.C, &messages, sql, args...)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) GetLastN(conversationID, n int) ([]Message, error) {
	db := db.DB()
	messages := make([]Message, 0)
	subquery := db.Builder.Select("*").From(MessagesTable).
		Where(sq.Eq{"conversation_id": conversationID}).
		OrderBy("created_at DESC").
		Limit(uint64(n))
	sql, args, err := subquery.ToSql()
	if err != nil {
		return nil, err
	}
	wrappedSQL := "SELECT * FROM (" + sql + ") sub ORDER BY created_at ASC"
	err = pgxscan.Select(context.Background(), db.C, &messages, wrappedSQL, args...)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) Create(conversationID int, role, content string) (*Message, error) {
	db := db.DB()
	msg := &Message{}
	stm := db.Builder.Insert(MessagesTable).
		Columns("conversation_id", "role", "content").
		Values(conversationID, role, content).
		Suffix("RETURNING id, conversation_id, role, content, created_at")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(context.Background(), db.C, msg, sql, args...)
	return msg, err
}

// ChunkService

type ChunkService struct{}

func (s *ChunkService) DeleteByFile(companyID int, filePath string) error {
	db := db.DB()
	stm := db.Builder.Delete(DocumentChunksTable).
		Where(sq.Eq{"company_id": companyID, "file_path": filePath})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *ChunkService) Insert(chunk *DocumentChunk, embedding []float32) error {
	db := db.DB()
	sql := `INSERT INTO chat_document_chunks (company_id, file_path, file_name, chunk_index, content, embedding)
		VALUES ($1, $2, $3, $4, $5, $6::vector) RETURNING id`
	return db.C.QueryRow(context.Background(), sql,
		chunk.CompanyID, chunk.FilePath, chunk.FileName, chunk.ChunkIndex, chunk.Content, vectorToString(embedding),
	).Scan(&chunk.ID)
}

func vectorToString(v []float32) string {
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = fmt.Sprintf("%g", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func (s *ChunkService) SearchSimilar(companyID int, embedding []float32, limit int) ([]DocumentChunk, error) {
	db := db.DB()
	chunks := make([]DocumentChunk, 0)
	sql := `SELECT id, company_id, file_path, file_name, chunk_index, content, created_at
		FROM chat_document_chunks
		WHERE company_id = $1
		ORDER BY embedding <=> $2::vector
		LIMIT $3`
	err := pgxscan.Select(context.Background(), db.C, &chunks, sql, companyID, vectorToString(embedding), limit)
	if err != nil {
		return nil, err
	}
	return chunks, nil
}
