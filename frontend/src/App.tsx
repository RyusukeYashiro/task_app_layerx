import { useEffect, useState } from 'react';
import { api, ApiException, ERROR_MESSAGES, type Task, type User } from './api';
import './App.css';

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [showLogin, setShowLogin] = useState(true);

  // Auth form
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');

  // Task form
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState(1);
  const [selectedAssignees, setSelectedAssignees] = useState<number[]>([]);

  // エラーメッセージを取得するヘルパー関数
  const getErrorMessage = (error: unknown): string => {
    if (error instanceof ApiException) {
      const customMessage = ERROR_MESSAGES[error.error.code];
      if (customMessage) {
        // 詳細情報がある場合は追加
        if (error.error.details?.field) {
          return `${customMessage}\n（項目: ${error.error.details.field}）`;
        }
        return customMessage;
      }
      return error.error.message || 'エラーが発生しました。';
    }
    return 'エラーが発生しました。';
  };

  useEffect(() => {
    const token = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    if (token && savedUser) {
      setUser(JSON.parse(savedUser));
      loadTasks();
      loadUsers();
    }
  }, []);

  const loadTasks = async () => {
    try {
      const data = await api.getTasks();
      setTasks(data);
    } catch (error) {
      console.error('Failed to load tasks:', error);
    }
  };

  const loadUsers = async () => {
    try {
      const data = await api.getUsers();
      setUsers(data);
    } catch (error) {
      console.error('Failed to load users:', error);
    }
  };

  const handleAuth = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = showLogin
        ? await api.login(email, password)
        : await api.signup(email, password, name);

      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      setUser(response.user);
      loadTasks();
      loadUsers();
      setEmail('');
      setPassword('');
      setName('');
    } catch (error) {
      const message = getErrorMessage(error);
      alert(`${showLogin ? 'ログイン' : 'ユーザー登録'}に失敗しました。\n\n${message}`);
    }
  };

  const handleLogout = async () => {
    await api.logout();
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    setTasks([]);
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createTask(title, description, priority, selectedAssignees.length > 0 ? selectedAssignees : undefined);
      setTitle('');
      setDescription('');
      setPriority(1);
      setSelectedAssignees([]);
      loadTasks();
      alert('✓ タスクを作成しました。');
    } catch (error) {
      const message = getErrorMessage(error);
      alert(`タスクの作成に失敗しました。\n\n${message}`);
    }
  };

  const handleUpdateStatus = async (taskId: number, status: Task['status']) => {
    try {
      await api.updateTask(taskId, { status });
      loadTasks();
    } catch (error) {
      const message = getErrorMessage(error);
      alert(`タスクの更新に失敗しました。\n\n${message}`);
    }
  };

  const handleDeleteTask = async (taskId: number) => {
    if (!confirm('このタスクを削除してもよろしいですか？\n\n削除すると元に戻せません。')) return;
    try {
      await api.deleteTask(taskId);
      loadTasks();
      alert('✓ タスクを削除しました。');
    } catch (error) {
      const message = getErrorMessage(error);
      if (error instanceof ApiException && error.error.code === 'FORBIDDEN') {
        alert(`タスクの削除に失敗しました。\n\n${message}\n\nこのタスクを削除できるのは、タスクの作成者のみです。`);
      } else {
        alert(`タスクの削除に失敗しました。\n\n${message}`);
      }
    }
  };

  if (!user) {
    return (
      <div className="container">
        <h1>Task Manager</h1>
        <div className="auth-container">
          <div className="tabs">
            <button
              className={showLogin ? 'active' : ''}
              onClick={() => setShowLogin(true)}
            >
              Login
            </button>
            <button
              className={!showLogin ? 'active' : ''}
              onClick={() => setShowLogin(false)}
            >
              Signup
            </button>
          </div>
          <form onSubmit={handleAuth}>
            {!showLogin && (
              <input
                type="text"
                placeholder="Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            )}
            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
            />
            <button type="submit">{showLogin ? 'Login' : 'Signup'}</button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="container">
      <header>
        <h1>Task Manager</h1>
        <div className="user-info">
          <span>Hello, {user.name}!</span>
          <button onClick={handleLogout}>Logout</button>
        </div>
      </header>

      <div className="task-form">
        <h2>Create New Task</h2>
        <form onSubmit={handleCreateTask}>
          <input
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
          <textarea
            placeholder="Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <select value={priority} onChange={(e) => setPriority(Number(e.target.value))}>
            <option value={1}>Priority: Low</option>
            <option value={2}>Priority: Medium</option>
            <option value={3}>Priority: High</option>
          </select>

          <div className="assignee-select">
            <label>Assign to: ({users.length} users available)</label>
            {users.length === 0 ? (
              <p style={{ color: '#999', fontSize: '14px' }}>No users available</p>
            ) : (
              users.map((u) => (
                <label key={u.id} className="checkbox-label">
                  <input
                    type="checkbox"
                    checked={selectedAssignees.includes(u.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedAssignees([...selectedAssignees, u.id]);
                      } else {
                        setSelectedAssignees(selectedAssignees.filter((id) => id !== u.id));
                      }
                    }}
                  />
                  {u.name} ({u.email})
                </label>
              ))
            )}
          </div>

          <button type="submit">Create Task</button>
        </form>
      </div>

      <div className="tasks-container">
        <h2>My Tasks ({tasks.length})</h2>
        {tasks.length === 0 ? (
          <p className="empty">No tasks yet. Create one above!</p>
        ) : (
          <div className="tasks-list">
            {tasks.map((task) => (
              <div key={task.id} className={`task-card ${task.status.toLowerCase()}`}>
                <div className="task-header">
                  <h3>{task.title}</h3>
                  <span className={`priority priority-${task.priority}`}>
                    P{task.priority}
                  </span>
                </div>
                <p>{task.description}</p>

                {task.assignees && task.assignees.length > 0 && (
                  <div className="task-assignees">
                    <strong>Assigned to:</strong>
                    <div className="assignee-list">
                      {task.assignees.map((assignee) => {
                        const assignedUser = users.find((u) => u.id === assignee.userId);
                        return (
                          <span key={assignee.userId} className="assignee-badge">
                            {assignedUser?.name || `User ${assignee.userId}`}
                          </span>
                        );
                      })}
                    </div>
                  </div>
                )}

                <div className="task-footer">
                  <select
                    value={task.status}
                    onChange={(e) => handleUpdateStatus(task.id, e.target.value as Task['status'])}
                  >
                    <option value="TODO">TODO</option>
                    <option value="IN_PROGRESS">IN PROGRESS</option>
                    <option value="DONE">DONE</option>
                  </select>
                  <button className="delete-btn" onClick={() => handleDeleteTask(task.id)}>
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
