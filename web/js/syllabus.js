(function () {
  const loading = document.getElementById('loading');
  const content = document.getElementById('courseContent');
  const errorBox = document.getElementById('errorBox');

  function showError(msg) {
    loading.classList.add('d-none');
    errorBox.textContent = msg;
    errorBox.classList.remove('d-none');
  }

  function progressLabel(mod) {
    const parts = [];
    if (mod.lessonDone) parts.push('Lesson ✓');
    parts.push(`Exercises ${mod.doneExercises}/${mod.totalExercises}`);
    if (mod.quizDone) parts.push('Quiz ✓');
    return parts.join(' · ');
  }

  function renderModules(modules) {
    const list = document.getElementById('moduleList');
    list.innerHTML = '';

    modules.forEach((mod, index) => {
      const item = document.createElement('div');
      item.className = 'list-group-item list-group-item-action' + (mod.locked ? ' module-locked' : '');

      const row = document.createElement('div');
      row.className = 'd-flex w-100 justify-content-between align-items-start gap-3';

      const left = document.createElement('div');
      const title = document.createElement('h2');
      title.className = 'h5 mb-1';
      title.textContent = `${index + 1}. ${mod.title}`;

      const desc = document.createElement('p');
      desc.className = 'mb-1 text-muted small';
      desc.textContent = mod.description;

      const progress = document.createElement('span');
      progress.className = 'badge text-bg-secondary progress-pill';
      progress.textContent = progressLabel(mod);

      left.append(title, desc, progress);

      const right = document.createElement('div');
      if (mod.locked) {
        const badge = document.createElement('span');
        badge.className = 'badge text-bg-warning';
        badge.textContent = 'Locked';
        right.appendChild(badge);
      } else if (mod.completed) {
        const badge = document.createElement('span');
        badge.className = 'badge text-bg-success';
        badge.textContent = 'Complete';
        right.appendChild(badge);
      }

      row.append(left, right);
      item.appendChild(row);

      if (!mod.locked) {
        item.style.cursor = 'pointer';
        item.addEventListener('click', () => {
          window.location.href = `/lesson.html?module=${encodeURIComponent(mod.id)}`;
        });
      }

      list.appendChild(item);
    });
  }

  async function init() {
    try {
      const course = await window.LearnGoAPI.getCourse();
      document.getElementById('courseTitle').textContent = course.title;
      document.getElementById('courseDescription').textContent = course.description;
      renderModules(course.modules);
      loading.classList.add('d-none');
      content.classList.remove('d-none');
    } catch (err) {
      showError(err.message);
    }
  }

  init();
})();
