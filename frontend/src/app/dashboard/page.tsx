'use client';

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getTodos, createTodo, updateTodo, deleteTodo } from '@/lib/api';
import { Todo } from '@/types';
import { useAuthGuard } from '@/hooks/useAuth';
import TodoItem from '@/components/TodoItem';
import styles from './dashboard.module.css';

export default function DashboardPage() {
  const router = useRouter();
  const {
    user,
    isLoading: authLoading,
    isAuthenticated,
    isAdmin: hasAdminRole,
    logout: logoutAndRedirect,
  } = useAuthGuard({
    redirectTo: '/login',
  });
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [newTodoTitle, setNewTodoTitle] = useState('');
  const [newTodoDescription, setNewTodoDescription] = useState('');

  const loadTodos = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getTodos();
      setTodos(data);
    } catch (error) {
      console.error('Failed to load todos:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      loadTodos();
    }
  }, [isAuthenticated, loadTodos]);

  const handleCreateTodo = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTodoTitle.trim()) return;

    try {
      const newTodo = await createTodo(newTodoTitle, newTodoDescription || undefined);
      setTodos([newTodo, ...todos]);
      setNewTodoTitle('');
      setNewTodoDescription('');
    } catch (error) {
      console.error('Failed to create todo:', error);
      alert('TODOの作成に失敗しました');
    }
  };

  const handleUpdateTodo = async (id: string, updates: any) => {
    try {
      const updatedTodo = await updateTodo(id, updates);
      setTodos(todos.map(todo => todo.id === id ? updatedTodo : todo));
    } catch (error) {
      console.error('Failed to update todo:', error);
      alert('TODOの更新に失敗しました');
    }
  };

  const handleDeleteTodo = async (id: string) => {
    if (!confirm('このTODOを削除しますか？')) return;

    try {
      await deleteTodo(id);
      setTodos(todos.filter(todo => todo.id !== id));
    } catch (error) {
      console.error('Failed to delete todo:', error);
      alert('TODOの削除に失敗しました');
    }
  };

  if (authLoading || loading) {
    return (
      <div className={styles.container}>
        <p>読み込み中...</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>TODO管理アプリ</h1>
        <div className={styles.userInfo}>
          <span>{user?.email}</span>
          {hasAdminRole && (
            <button
              onClick={() => router.push('/admin')}
              className={styles.adminButton}
            >
              管理画面
            </button>
          )}
          <button onClick={logoutAndRedirect} className={styles.logoutButton}>
            ログアウト
          </button>
        </div>
      </header>

      <main className={styles.main}>
        <form onSubmit={handleCreateTodo} className={styles.createForm}>
          <input
            type="text"
            value={newTodoTitle}
            onChange={(e) => setNewTodoTitle(e.target.value)}
            placeholder="新しいTODOのタイトル"
            className={styles.input}
          />
          <textarea
            value={newTodoDescription}
            onChange={(e) => setNewTodoDescription(e.target.value)}
            placeholder="説明（任意）"
            className={styles.textarea}
          />
          <button type="submit" className={styles.createButton}>
            追加
          </button>
        </form>

        <div className={styles.todoList}>
          {todos.length === 0 ? (
            <p className={styles.emptyMessage}>TODOがありません</p>
          ) : (
            todos.map(todo => (
              <TodoItem
                key={todo.id}
                todo={todo}
                onUpdate={handleUpdateTodo}
                onDelete={handleDeleteTodo}
              />
            ))
          )}
        </div>
      </main>
    </div>
  );
}
