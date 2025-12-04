import { Todo } from '@/types';

const API_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8000';

function getHeaders() {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
}

export async function getTodos(): Promise<Todo[]> {
  const response = await fetch(`${API_URL}/api/todos`, {
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch todos');
  }

  return response.json();
}

export async function createTodo(title: string, description?: string): Promise<Todo> {
  const response = await fetch(`${API_URL}/api/todos`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ title, description }),
  });

  if (!response.ok) {
    throw new Error('Failed to create todo');
  }

  return response.json();
}

export async function updateTodo(
  id: string,
  updates: { title?: string; description?: string; completed?: boolean }
): Promise<Todo> {
  const response = await fetch(`${API_URL}/api/todos/${id}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(updates),
  });

  if (!response.ok) {
    throw new Error('Failed to update todo');
  }

  return response.json();
}

export async function deleteTodo(id: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/todos/${id}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to delete todo');
  }
}

// Admin API
export async function getAllUsers() {
  const response = await fetch(`${API_URL}/api/admin/users`, {
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch users');
  }

  return response.json();
}

export async function updateUserRole(userId: string, role: string) {
  const response = await fetch(`${API_URL}/api/admin/users/${userId}/role`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify({ role }),
  });

  if (!response.ok) {
    throw new Error('Failed to update user role');
  }

  return response.json();
}

export async function deleteUser(userId: string) {
  const response = await fetch(`${API_URL}/api/admin/users/${userId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to delete user');
  }
}

export async function getAllTodosAdmin() {
  const response = await fetch(`${API_URL}/api/admin/todos`, {
    headers: getHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch all todos');
  }

  return response.json();
}
