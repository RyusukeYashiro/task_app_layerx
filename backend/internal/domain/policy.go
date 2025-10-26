package domain

// ユーザーがタスクを閲覧できるかチェックする
func CanViewTask(task *Task, assignees []*TaskAssignee, userID int64) bool {
	// オーナーなら閲覧可能
	if task.OwnerID == userID {
		return true
	}

	// アサインされているかチェック
	for _, assignee := range assignees {
		if assignee.UserID == userID {
			return true
		}
	}
	// オーナーでもアサイン先でもない
	return false
}

// ユーザーがタスクを編集できるかチェックする
func CanEditTask(task *Task, userID int64) bool {
	return task.OwnerID == userID
}

// ユーザーがタスクを削除できるかチェックする
func CanDeleteTask(task *Task, userID int64) bool {
	return task.OwnerID == userID
}

// ユーザーがアサインを管理できるかチェックする
func CanManageAssignees(task *Task, userID int64) bool {
	return task.OwnerID == userID
}
