package services

import (
	"context"
	"errors"
	"strings"

	"mo-da-backend/internal/database"
)

// CollaborationService backs per-detail comments and activity history.
type CollaborationService struct{}

func NewCollaborationService() *CollaborationService {
	return &CollaborationService{}
}

func normalizeEntityKey(value string) string {
	return strings.TrimSpace(value)
}

func (s *CollaborationService) ListComments(ctx context.Context, clientID, entityType, entityID string) ([]map[string]interface{}, error) {
	clientID = normalizeClientID(clientID)
	entityType = normalizeEntityKey(entityType)
	entityID = normalizeEntityKey(entityID)
	if entityType == "" || entityID == "" {
		return []map[string]interface{}{}, nil
	}

	rows, err := database.Pool.Query(ctx, `
		SELECT id, entity_type, entity_id, username, COALESCE(full_name, ''),
		       content, to_char(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI'),
		       COALESCE(parent_id, 0),
		       COALESCE(to_char(updated_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI'), '')
		FROM entity_comments
		WHERE client_id = $1 AND entity_type = $2 AND entity_id = $3
		ORDER BY created_at DESC, id DESC
	`, clientID, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var (
			id                                   int
			entityTypeVal, entityIDVal, username string
			fullName, content, createdAt         string
			parentID                             int
			updatedAt                            string
		)
		if err := rows.Scan(&id, &entityTypeVal, &entityIDVal, &username, &fullName, &content, &createdAt, &parentID, &updatedAt); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":         id,
			"entityType": entityTypeVal,
			"entityId":   entityIDVal,
			"username":   username,
			"fullName":   fullName,
			"content":    content,
			"parentId":   parentID,
			"updatedAt":  updatedAt,
			"createdAt":  createdAt,
		})
	}
	return results, rows.Err()
}

func (s *CollaborationService) CreateComment(ctx context.Context, clientID, entityType, entityID, username, fullName, content string, parentID int) (map[string]interface{}, error) {
	clientID = normalizeClientID(clientID)
	entityType = normalizeEntityKey(entityType)
	entityID = normalizeEntityKey(entityID)
	content = strings.TrimSpace(content)
	if entityType == "" || entityID == "" {
		return nil, errors.New("Thiếu thông tin đối tượng bình luận")
	}
	if content == "" {
		return nil, errors.New("Nội dung bình luận không được để trống")
	}
	if len([]rune(content)) > 4000 {
		return nil, errors.New("Nội dung bình luận quá dài")
	}

	action := "Đã bình luận"
	detail := content
	if parentID > 0 {
		var parentAuthor, parentName string
		var parentParentID int
		if err := database.Pool.QueryRow(ctx, `
			SELECT username, COALESCE(NULLIF(full_name, ''), username), COALESCE(parent_id, 0)
			FROM entity_comments
			WHERE client_id = $1 AND id = $2 AND entity_type = $3 AND entity_id = $4
		`, clientID, parentID, entityType, entityID).Scan(&parentAuthor, &parentName, &parentParentID); err != nil {
			return nil, errors.New("Bình luận được trả lời không tồn tại")
		}
		if parentParentID > 0 {
			parentID = parentParentID
		}
		action = "Đã trả lời bình luận"
		detail = "@" + parentName + ": " + content
	}

	var (
		id        int
		createdAt string
	)
	err := database.Pool.QueryRow(ctx, `
		INSERT INTO entity_comments (client_id, entity_type, entity_id, username, full_name, content, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, 0))
		RETURNING id, to_char(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI')
	`, clientID, entityType, entityID, username, fullName, content, parentID).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	// Comments are actions too: record them in the detail activity timeline.
	_, _ = s.CreateActivity(ctx, clientID, entityType, entityID, username, fullName, action, "comment", truncateRunes(detail, 180))

	return map[string]interface{}{
		"id":         id,
		"entityType": entityType,
		"entityId":   entityID,
		"username":   username,
		"fullName":   fullName,
		"content":    content,
		"parentId":   parentID,
		"updatedAt":  "",
		"createdAt":  createdAt,
	}, nil
}

func (s *CollaborationService) UpdateComment(ctx context.Context, clientID string, commentID int, username, content string) (map[string]interface{}, error) {
	clientID = normalizeClientID(clientID)
	if commentID <= 0 {
		return nil, errors.New("Bình luận không hợp lệ")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("Nội dung bình luận không được để trống")
	}
	if len([]rune(content)) > 4000 {
		return nil, errors.New("Nội dung bình luận quá dài")
	}

	var (
		owner, entityType, entityID, fullName, createdAt string
		oldContent                                       string
		parentID                                         int
	)
	if err := database.Pool.QueryRow(ctx, `
		SELECT username, entity_type, entity_id, COALESCE(full_name, ''), content,
		       to_char(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI'),
		       COALESCE(parent_id, 0)
		FROM entity_comments
		WHERE client_id = $1 AND id = $2
	`, clientID, commentID).Scan(&owner, &entityType, &entityID, &fullName, &oldContent, &createdAt, &parentID); err != nil {
		return nil, errors.New("Không tìm thấy bình luận")
	}
	if owner != username {
		return nil, errors.New("Bạn chỉ có thể chỉnh sửa bình luận của chính mình")
	}

	var updatedAt string
	if err := database.Pool.QueryRow(ctx, `
		UPDATE entity_comments SET content = $3, updated_at = NOW()
		WHERE client_id = $1 AND id = $2
		RETURNING to_char(updated_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI')
	`, clientID, commentID, content).Scan(&updatedAt); err != nil {
		return nil, err
	}

	_, _ = s.CreateActivity(ctx, clientID, entityType, entityID, username, fullName, "Đã chỉnh sửa bình luận", "update", truncateRunes(content, 180))

	return map[string]interface{}{
		"id":         commentID,
		"entityType": entityType,
		"entityId":   entityID,
		"username":   username,
		"fullName":   fullName,
		"content":    content,
		"parentId":   parentID,
		"updatedAt":  updatedAt,
		"createdAt":  createdAt,
	}, nil
}

func (s *CollaborationService) DeleteComment(ctx context.Context, clientID string, commentID int, username string) error {
	clientID = normalizeClientID(clientID)
	if commentID <= 0 {
		return errors.New("Bình luận không hợp lệ")
	}
	var owner string
	if err := database.Pool.QueryRow(ctx,
		`SELECT username FROM entity_comments WHERE client_id = $1 AND id = $2`,
		clientID, commentID,
	).Scan(&owner); err != nil {
		return errors.New("Không tìm thấy bình luận")
	}
	if owner != username {
		return errors.New("Bạn chỉ có thể xoá bình luận của chính mình")
	}
	_, err := database.Pool.Exec(ctx,
		`DELETE FROM entity_comments WHERE client_id = $1 AND (id = $2 OR parent_id = $2)`,
		clientID, commentID,
	)
	return err
}

func (s *CollaborationService) ListActivities(ctx context.Context, clientID, entityType, entityID string) ([]map[string]interface{}, error) {
	clientID = normalizeClientID(clientID)
	entityType = normalizeEntityKey(entityType)
	entityID = normalizeEntityKey(entityID)
	if entityType == "" || entityID == "" {
		return []map[string]interface{}{}, nil
	}

	rows, err := database.Pool.Query(ctx, `
		SELECT id, entity_type, entity_id, username, COALESCE(full_name, ''),
		       action, action_type, COALESCE(detail, ''),
		       to_char(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI')
		FROM entity_activities
		WHERE client_id = $1 AND entity_type = $2 AND entity_id = $3
		ORDER BY created_at DESC, id DESC
		LIMIT 200
	`, clientID, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var (
			id                                              int
			entityTypeVal, entityIDVal, username            string
			fullName, action, actionType, detail, createdAt string
		)
		if err := rows.Scan(&id, &entityTypeVal, &entityIDVal, &username, &fullName, &action, &actionType, &detail, &createdAt); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":         id,
			"entityType": entityTypeVal,
			"entityId":   entityIDVal,
			"username":   username,
			"fullName":   fullName,
			"action":     action,
			"actionType": actionType,
			"detail":     detail,
			"createdAt":  createdAt,
		})
	}
	return results, rows.Err()
}

func (s *CollaborationService) CreateActivity(ctx context.Context, clientID, entityType, entityID, username, fullName, action, actionType, detail string) (map[string]interface{}, error) {
	clientID = normalizeClientID(clientID)
	entityType = normalizeEntityKey(entityType)
	entityID = normalizeEntityKey(entityID)
	action = strings.TrimSpace(action)
	actionType = strings.TrimSpace(actionType)
	if entityType == "" || entityID == "" || action == "" {
		return nil, errors.New("Thiếu thông tin nhật ký")
	}
	if actionType == "" {
		actionType = "update"
	}
	detail = truncateRunes(strings.TrimSpace(detail), 500)

	// Avoid flooding the timeline with repeated "view" entries for the same user.
	if actionType == "view" {
		var recent int
		_ = database.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM entity_activities
			WHERE client_id = $1 AND entity_type = $2 AND entity_id = $3
			  AND username = $4 AND action_type = 'view'
			  AND created_at > NOW() - INTERVAL '10 minutes'
		`, clientID, entityType, entityID, username).Scan(&recent)
		if recent > 0 {
			return nil, nil
		}
	}

	var (
		id        int
		createdAt string
	)
	err := database.Pool.QueryRow(ctx, `
		INSERT INTO entity_activities (client_id, entity_type, entity_id, username, full_name, action, action_type, detail)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, to_char(created_at AT TIME ZONE 'Asia/Ho_Chi_Minh', 'DD/MM/YYYY HH24:MI')
	`, clientID, entityType, entityID, username, fullName, action, actionType, detail).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         id,
		"entityType": entityType,
		"entityId":   entityID,
		"username":   username,
		"fullName":   fullName,
		"action":     action,
		"actionType": actionType,
		"detail":     detail,
		"createdAt":  createdAt,
	}, nil
}

func truncateRunes(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max]) + "…"
}
