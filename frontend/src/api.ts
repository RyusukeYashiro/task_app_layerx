const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface User {
  id: number;
  email: string;
  name: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Assignee {
  userId: number;
  assignedBy: number;
  assignedAt: string;
}

export interface Task {
  id: number;
  ownerId: number;
  title: string;
  description: string;
  dueDate: string | null;
  status: 'TODO' | 'IN_PROGRESS' | 'DONE';
  priority: number;
  assignees: Assignee[];
  createdAt: string;
  updatedAt: string;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

export class ApiException extends Error {
  statusCode: number;
  error: ApiError;

  constructor(statusCode: number, error: ApiError) {
    super(error.message);
    this.name = 'ApiException';
    this.statusCode = statusCode;
    this.error = error;
  }
}

// エラーコードに応じた日本語メッセージ
export const ERROR_MESSAGES: Record<string, string> = {
  // 認証エラー
  UNAUTHORIZED: '認証が必要です。ログインしてください。',
  INVALID_TOKEN: 'トークンが無効または期限切れです。再度ログインしてください。',
  TOKEN_EXPIRED: 'セッションが期限切れです。再度ログインしてください。',
  
  // 認可エラー
  FORBIDDEN: 'この操作を実行する権限がありません。',
  
  // Not Found
  NOT_FOUND: '指定されたリソースが見つかりません。',
  
  // バリデーションエラー
  VALIDATION_ERROR: '入力内容に誤りがあります。',
  INVALID_REQUEST: 'リクエストの形式が正しくありません。',
  INVALID_DATE_FORMAT: '日付の形式が正しくありません（ISO8601形式で入力してください）。',
  
  // 重複エラー
  CONFLICT: 'すでに存在するデータです。',
  
  // その他
  INTERNAL_ERROR: 'サーバーエラーが発生しました。しばらくしてから再度お試しください。',
};

class ApiClient {
  private getHeaders(): HeadersInit {
    const token = localStorage.getItem('token');
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      try {
        const errorData: ApiError = await response.json();
        throw new ApiException(response.status, errorData);
      } catch (e) {
        if (e instanceof ApiException) throw e;
        // JSONパースに失敗した場合
        throw new ApiException(response.status, {
          code: 'UNKNOWN_ERROR',
          message: `HTTP Error ${response.status}`,
        });
      }
    }
    
    // 204 No Content の場合
    if (response.status === 204) {
      return undefined as T;
    }
    
    return response.json();
  }

  async signup(email: string, password: string, name: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password, name }),
    });
    return this.handleResponse<AuthResponse>(response);
  }

  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    return this.handleResponse<AuthResponse>(response);
  }

  async logout(): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      headers: this.getHeaders(),
    });
    return this.handleResponse<void>(response);
  }

  async getUsers(): Promise<User[]> {
    const response = await fetch(`${API_BASE_URL}/users`, {
      headers: this.getHeaders(),
    });
    return this.handleResponse<User[]>(response);
  }

  async getTasks(): Promise<Task[]> {
    const response = await fetch(`${API_BASE_URL}/tasks`, {
      headers: this.getHeaders(),
    });
    return this.handleResponse<Task[]>(response);
  }

  async createTask(title: string, description: string, priority: number, assigneeIDs?: number[]): Promise<Task> {
    const response = await fetch(`${API_BASE_URL}/tasks`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ title, description, priority, assigneeIDs }),
    });
    return this.handleResponse<Task>(response);
  }

  async updateTask(id: number, updates: Partial<Task>): Promise<Task> {
    const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
      method: 'PATCH',
      headers: this.getHeaders(),
      body: JSON.stringify(updates),
    });
    return this.handleResponse<Task>(response);
  }

  async deleteTask(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/tasks/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    });
    return this.handleResponse<void>(response);
  }
}

export const api = new ApiClient();
