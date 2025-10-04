<script>
  import { onMount } from "svelte";

  let dark = false;

  onMount(() => {
    dark = localStorage.getItem("theme") === "dark";
    applyTheme();
  });

  function toggleTheme() {
    dark = !dark;
    localStorage.setItem("theme", dark ? "dark" : "light");
    applyTheme();
  }

  function applyTheme() {
    if (dark) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }
</script>

<button
  on:click={toggleTheme}
  class="px-3 py-1 rounded-lg bg-gray-200 dark:bg-gray-700 dark:text-gray-200 text-sm"
>
  {dark ? "☾ Dark" : "☀ Light"}
</button>

