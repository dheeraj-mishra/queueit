<script>
  import { updateTask, deleteTask } from "../lib/store.js";
  export let task;

  let editing = false;
  let newTitle = task.title;
  let newDeadline = task.deadline;

  function toggleDone() {
    updateTask(task.id, { done: !task.done });
  }

  function handleEdit() {
    if (editing) {
      updateTask(task.id, { title: newTitle, deadline: newDeadline });
    }
    editing = !editing;
  }

  function handleDelete() {
    deleteTask(task.id);
  }

  function timeLeft() {
    const now = new Date();
    const end = new Date(task.deadline);
    const diff = end - now;
    if (diff <= 0) return "Expired";
    const hours = Math.floor(diff / 1000 / 60 / 60);
    const mins = Math.floor((diff / 1000 / 60) % 60);
    return `${hours}h ${mins}m left`;
  }
</script>

<div
  class="p-4 rounded-xl shadow-md border dark:border-gray-700 dark:bg-gray-800 bg-white flex flex-col gap-2"
>
  {#if editing}
    <input
      type="text"
      bind:value={newTitle}
      class="px-2 py-1 border rounded-lg dark:bg-gray-700 dark:text-gray-200"
    />
    <input
      type="datetime-local"
      bind:value={newDeadline}
      class="px-2 py-1 border rounded-lg dark:bg-gray-700 dark:text-gray-200"
    />
  {:else}
    <h3 class="font-semibold text-lg {task.done ? 'line-through text-gray-400' : ''}">
      {task.title}
    </h3>
    <p class="text-sm text-gray-500 dark:text-gray-400">{timeLeft()}</p>
  {/if}

  <div class="flex gap-2 mt-2">
    <button
      on:click={toggleDone}
      class="px-3 py-1 text-sm rounded-lg bg-green-500 text-white"
    >
      {task.done ? "Undo" : "Done"}
    </button>
    <button
      on:click={handleEdit}
      class="px-3 py-1 text-sm rounded-lg bg-yellow-500 text-white"
    >
      {editing ? "Save" : "Edit"}
    </button>
    <button
      on:click={handleDelete}
      class="px-3 py-1 text-sm rounded-lg bg-red-500 text-white"
    >
      Delete
    </button>
  </div>
</div>

