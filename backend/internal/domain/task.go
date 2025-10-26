package domain

import (
	"strings"
	"time"
)

type TaskStatus string

const (
	TaskStatusTODO TaskStatus = "TODO"
	TaskStatusIN_PROGRESS TaskStatus = "IN_PROGRESS"
	TaskStatusDONE TaskStatus = "DONE"
)

type Task struct {
	ID int64
	OwnerID int64
	Title string
	Description *string
	DueDate *time.Time
	Status TaskStatus
	Priority int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewTask(clock Clock , ownerID int64 , title string) (*Task , error) {
	now := clock.Now()
	task := &Task{
		OwnerID: ownerID,
		Title: strings.TrimSpace(title),
		Status: TaskStatusTODO,
		Priority: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := task.ValidateTitle(); err != nil {
		return nil , err
	}
	if err := task.ValidatePriority(); err != nil {
		return nil , err
	}
	return task , nil
}

func (t *Task) ValidateTitle() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrTitleRequired
	}
	if len(t.Title) > 255 {
		return ErrTitleTooLong
	}
	return nil
}

func (t *Task) ValidatePriority() error {
	if t.Priority < 0 || t.Priority > 5 {
		return ErrInvalidPriority
	}
	return nil
}

func (t *Task) ValidateStatusTransaction(nextStatus TaskStatus) error {
	if !isValidStatus(nextStatus) {
		return ErrInvalidStatusTransition
	}

	switch t.Status {
	case TaskStatusTODO:
		if nextStatus == TaskStatusIN_PROGRESS || nextStatus == TaskStatusDONE {
			return nil
		}
	case TaskStatusIN_PROGRESS:
		if nextStatus == TaskStatusDONE || nextStatus == TaskStatusTODO {
			return nil
		}
	case TaskStatusDONE:
		if nextStatus == TaskStatusTODO {
			return nil
		}
	}
	return ErrInvalidStatusTransition
}

func isValidStatus(status TaskStatus) bool {
	return status == TaskStatusTODO || status == TaskStatusIN_PROGRESS || status == TaskStatusDONE
}

// touch 更新日時を更新
func (t *Task) touch(clock Clock) {
	t.UpdatedAt = clock.Now()
}

// タイトルの更新
func (t *Task) UpdateTitle (clock Clock , title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrTitleRequired
	}
	if len(title) > 255 {
		return ErrTitleTooLong
	}
	t.Title = title
	t.touch(clock)
	return nil
}


// 説明の更新
func (t *Task) UpdateDescription (clock Clock , description *string) {
	if description == nil {
		t.Description = nil
		t.touch(clock)
		return
	}
	trimmed := strings.TrimSpace(*description)
	if trimmed == "" {
		t.Description = nil
	} else {
		t.Description = &trimmed
	}
	t.touch(clock)
}

// 期日の更新
func (t *Task) UpdateDueDate (clock Clock , dueDate *time.Time) {
	t.DueDate = dueDate
	t.touch(clock)
}

// 優先度の更新
func (t *Task) UpdatePriority (clock Clock , priority int) error {
	t.Priority = priority
	if err := t.ValidatePriority(); err != nil {
		return err
	}
	t.touch(clock)
	return nil
}

func (t *Task) SoftDelete (clock Clock) {
	if t.DeletedAt != nil {
		return
	}
	now := clock.Now()
	t.DeletedAt = &now
	t.touch(clock)
}

func (t *Task) IsOwner(userID int64) bool {
	return t.OwnerID == userID
}

