(function () {
  const KEY = 'learn-go-theme';

  function currentTheme() {
    return document.documentElement.getAttribute('data-bs-theme') || 'light';
  }

  function applyTheme(theme) {
    document.documentElement.setAttribute('data-bs-theme', theme);
    localStorage.setItem(KEY, theme);
    const btn = document.getElementById('themeToggle');
    if (btn) {
      btn.textContent = theme === 'dark' ? '☀️' : '🌙';
    }
    if (window.learnGoEditor) {
      window.learnGoEditor.setTheme(theme === 'dark' ? 'dracula' : 'eclipse');
    }
  }

  function toggleTheme() {
    applyTheme(currentTheme() === 'dark' ? 'light' : 'dark');
  }

  const saved = localStorage.getItem(KEY);
  applyTheme(saved === 'dark' ? 'dark' : 'light');

  const btn = document.getElementById('themeToggle');
  if (btn) {
    btn.addEventListener('click', toggleTheme);
  }

  window.learnGoTheme = { applyTheme, currentTheme };
})();
