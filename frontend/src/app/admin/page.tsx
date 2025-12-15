'use client';

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getAllUsers, updateUserRole, deleteUser, getAllTodosAdmin } from '@/lib/api';
import { User, Todo } from '@/types';
import { useAuthGuard } from '@/hooks/useAuth';
import styles from './admin.module.css';

export default function AdminPage() {
  const router = useRouter();
  const {
    isLoading: authLoading,
    isAuthenticated,
    logout: logoutAndRedirect,
  } = useAuthGuard({
    redirectTo: '/dashboard',
    requireAdmin: true,
  });
  const [users, setUsers] = useState<User[]>([]);
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'users' | 'todos'>('users');

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [usersData, todosData] = await Promise.all([
        getAllUsers(),
        getAllTodosAdmin(),
      ]);
      setUsers(usersData);
      setTodos(todosData);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      loadData();
    }
  }, [isAuthenticated, loadData]);

  const handleUpdateRole = async (userId: string, newRole: string) => {
    try {
      await updateUserRole(userId, newRole);
      setUsers(users.map(user =>
        user.id === userId ? { ...user, role: newRole } : user
      ));
      alert('ユーザーのロールを更新しました');
    } catch (error) {
      console.error('Failed to update role:', error);
      alert('ロールの更新に失敗しました');
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (!confirm('このユーザーを削除しますか？')) return;

    try {
      await deleteUser(userId);
      setUsers(users.filter(user => user.id !== userId));
      alert('ユーザーを削除しました');
    } catch (error: any) {
      console.error('Failed to delete user:', error);
      alert(error.message || 'ユーザーの削除に失敗しました');
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
        <h1>管理画面</h1>
        <div className={styles.actions}>
          <button
            onClick={() => router.push('/dashboard')}
            className={styles.backButton}
          >
            ダッシュボードに戻る
          </button>
          <button onClick={logoutAndRedirect} className={styles.logoutButton}>
            ログアウト
          </button>
        </div>
      </header>

      <main className={styles.main}>
        <div className={styles.tabs}>
          <button
            className={`${styles.tab} ${activeTab === 'users' ? styles.activeTab : ''}`}
            onClick={() => setActiveTab('users')}
          >
            ユーザー一覧 ({users.length})
          </button>
          <button
            className={`${styles.tab} ${activeTab === 'todos' ? styles.activeTab : ''}`}
            onClick={() => setActiveTab('todos')}
          >
            全TODO ({todos.length})
          </button>
        </div>

        {activeTab === 'users' && (
          <div className={styles.section}>
            <h2>ユーザー管理</h2>
            <div className={styles.table}>
              <table>
                <thead>
                  <tr>
                    <th>メールアドレス</th>
                    <th>ロール</th>
                    <th>登録日時</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map(user => (
                    <tr key={user.id}>
                      <td>{user.email}</td>
                      <td>
                        <select
                          value={user.role}
                          onChange={(e) => handleUpdateRole(user.id, e.target.value)}
                          className={styles.select}
                        >
                          <option value="user">user</option>
                          <option value="admin">admin</option>
                        </select>
                      </td>
                      <td>{new Date(user.created_at).toLocaleDateString('ja-JP')}</td>
                      <td>
                        <button
                          onClick={() => handleDeleteUser(user.id)}
                          className={styles.deleteButton}
                        >
                          削除
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {activeTab === 'todos' && (
          <div className={styles.section}>
            <h2>全TODO一覧</h2>
            <div className={styles.todoGrid}>
              {todos.map(todo => (
                <div key={todo.id} className={styles.todoCard}>
                  <h3>{todo.title}</h3>
                  {todo.description && <p>{todo.description}</p>}
                  <div className={styles.todoMeta}>
                    <span className={todo.completed ? styles.completed : styles.pending}>
                      {todo.completed ? '完了' : '未完了'}
                    </span>
                    <span className={styles.userId}>
                      ユーザーID: {todo.user_id.substring(0, 8)}...
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
