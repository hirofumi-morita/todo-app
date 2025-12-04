'use client';

import { useState } from 'react';
import { Todo } from '@/types';
import styles from './TodoItem.module.css';

interface TodoItemProps {
  todo: Todo;
  onUpdate: (id: string, updates: any) => void;
  onDelete: (id: string) => void;
}

export default function TodoItem({ todo, onUpdate, onDelete }: TodoItemProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState(todo.title);
  const [description, setDescription] = useState(todo.description || '');

  const handleSave = () => {
    onUpdate(todo.id, { title, description });
    setIsEditing(false);
  };

  const handleCancel = () => {
    setTitle(todo.title);
    setDescription(todo.description || '');
    setIsEditing(false);
  };

  const handleToggleComplete = () => {
    onUpdate(todo.id, { completed: !todo.completed });
  };

  if (isEditing) {
    return (
      <div className={styles.todoItem}>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className={styles.input}
        />
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className={styles.textarea}
          placeholder="説明（任意）"
        />
        <div className={styles.actions}>
          <button onClick={handleSave} className={styles.saveButton}>
            保存
          </button>
          <button onClick={handleCancel} className={styles.cancelButton}>
            キャンセル
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`${styles.todoItem} ${todo.completed ? styles.completed : ''}`}>
      <div className={styles.checkbox}>
        <input
          type="checkbox"
          checked={todo.completed}
          onChange={handleToggleComplete}
        />
      </div>
      <div className={styles.content}>
        <h3 className={styles.title}>{todo.title}</h3>
        {todo.description && <p className={styles.description}>{todo.description}</p>}
      </div>
      <div className={styles.actions}>
        <button onClick={() => setIsEditing(true)} className={styles.editButton}>
          編集
        </button>
        <button onClick={() => onDelete(todo.id)} className={styles.deleteButton}>
          削除
        </button>
      </div>
    </div>
  );
}
