package services

import (
	"errors"

	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/repositories"
	appwebsocket "trithong.com/task-golang/internal/websocket"
)

type CommentService interface {
	Create(comment entities.Comment) (entities.Comment, error)
	GetByTaskID(taskID int, userID int) ([]entities.Comment, error) // ← thêm userID
}

type commentService struct {
	commentRepo repositories.CommentRepository
	hub         *appwebsocket.Hub
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	hub *appwebsocket.Hub,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		hub:         hub,
	}
}

func (s *commentService) Create(comment entities.Comment) (entities.Comment, error) {

	belongs, err := s.commentRepo.CheckTaskBelongsToUser(comment.TaskID, comment.UserID)
	if err != nil {
		return entities.Comment{}, err
	}
	if !belongs {
		return entities.Comment{}, errors.New("task not found or unauthorized")
	}

	newComment, err := s.commentRepo.Create(comment)
	if err != nil {
		return entities.Comment{}, err
	}

	// Chỉ gửi cho đúng user tạo comment thôi
	s.hub.SendToUser(
		comment.UserID,
		"task.comment.created",
		newComment,
		"Comment created successfully",
	)

	return newComment, nil
}

func (s *commentService) GetByTaskID(
	taskID int,
	userID int, // ← thêm tham số này
) ([]entities.Comment, error) {

	// Kiểm tra task có thuộc project của user không
	belongs, err := s.commentRepo.CheckTaskBelongsToUser(taskID, userID)

	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, errors.New("task not found or unauthorized")
	}

	return s.commentRepo.GetByTaskID(taskID)
}
