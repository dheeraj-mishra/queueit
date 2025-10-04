import { writable } from 'svelte/store';

// Example structure of a task
// { id: 1, title: "Do homework", deadline: "2025-10-05T18:00", done: false }
export const tasks = writable([]);

export function addTask(task) {
  tasks.update(ts => [...ts, { id: Date.now(), ...task }]);
}

export function updateTask(id, updated) {
  tasks.update(ts => ts.map(t => t.id === id ? { ...t, ...updated } : t));
}

export function deleteTask(id) {
  tasks.update(ts => ts.filter(t => t.id !== id));
}

