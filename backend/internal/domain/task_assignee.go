package domain

import "time"

type TaskAssignee struct {
	TaskID int64
	UserID int64
	AssignedBy int64
	CreatedAt time.Time
	UpdatedAt time.Time
}


func NewTaskAssignee(clock Clock , taskID int64 , userID int64 , assignedBy int64) (*TaskAssignee , error) {
	now := clock.Now()
	assignee := &TaskAssignee{
		TaskID: taskID,
		UserID: userID,
		AssignedBy: assignedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return assignee , nil
}
