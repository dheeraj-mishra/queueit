<script>
  import { addTask } from "../lib/store.js";
  let title = "";
  let deadline = "";

  function handleSubmit() {
    if (!title || !deadline) return;
    addTask({ title, deadline, done: false });
    title = "";
    deadline = "";
  }
</script>

<div class="flex gap-2 items-center my-4">
  <input
    type="text"
    placeholder="Task title"
    bind:value={title}
    class="flex-1 px-3 py-2 border rounded-lg dark:bg-gray-800 dark:text-gray-200"
  />
  <input
    type="datetime-local"
    bind:value={deadline}
    class="px-3 py-2 border rounded-lg dark:bg-gray-800 dark:text-gray-200"
  />
  <button
    on:click={handleSubmit}
    class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
  >
    Add
  </button>
</div>

