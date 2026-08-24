(function () {
  const params = new URLSearchParams(window.location.search);
  const moduleId = params.get('module');
  const loading = document.getElementById('loading');
  const content = document.getElementById('lessonContent');
  const errorBox = document.getElementById('errorBox');

  if (!moduleId) {
    showError('Missing module parameter.');
    return;
  }

  function showError(msg) {
    loading.classList.add('d-none');
    errorBox.textContent = msg;
    errorBox.classList.remove('d-none');
  }

  function renderProject(mod, progress) {
    if (!mod.project || !mod.project.items || mod.project.items.length === 0) {
      return true;
    }

    const section = document.getElementById('projectSection');
    section.classList.remove('d-none');
    document.getElementById('projectTitle').textContent = mod.project.title || 'Project';

    const list = document.getElementById('projectChecklist');
    list.innerHTML = '';
    let allDone = true;

    mod.project.items.forEach((item) => {
      const itemId = `${moduleId}:project:${item.id}`;
      const done = progress.items[itemId];
      if (!done) allDone = false;

      const row = document.createElement('label');
      row.className = 'list-group-item d-flex gap-2 align-items-start';
      row.innerHTML = `
        <input class="form-check-input mt-1" type="checkbox" ${done ? 'checked' : ''} data-item-id="${itemId}">
        <span>${escapeHtml(item.text)}</span>`;
      list.appendChild(row);
    });

    list.querySelectorAll('input[type=checkbox]').forEach((cb) => {
      cb.addEventListener('change', async (e) => {
        const itemId = e.target.getAttribute('data-item-id');
        if (e.target.checked) {
          await window.LearnGoAPI.markComplete(itemId);
        }
        const fresh = await window.LearnGoAPI.getProgress();
        renderProject(mod, fresh);
        updateQuizGate(mod, fresh, document.getElementById('lessonDoneBadge').classList.contains('d-none') === false);
      });
    });

    if (allDone) {
      document.getElementById('projectDoneBadge').classList.remove('d-none');
    } else {
      document.getElementById('projectDoneBadge').classList.add('d-none');
    }

    return allDone;
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  function updateQuizGate(mod, progress, lessonDone) {
    let allExercisesDone = mod.exercises.length === 0;
    mod.exercises.forEach((ex) => {
      if (!progress.items[`exercise:${ex.id}`]) allExercisesDone = false;
    });

    let projectDone = true;
    if (mod.project && mod.project.items && mod.project.items.length > 0) {
      mod.project.items.forEach((item) => {
        if (!progress.items[`${moduleId}:project:${item.id}`]) projectDone = false;
      });
    }

    const quizLink = document.getElementById('quizLink');
    quizLink.href = `/quiz.html?id=${encodeURIComponent(mod.quizId)}&module=${encodeURIComponent(moduleId)}`;
    if (allExercisesDone && lessonDone && projectDone) {
      quizLink.classList.remove('disabled');
    } else {
      quizLink.classList.add('disabled');
    }
  }

  async function init() {
    try {
      const [mod, progress] = await Promise.all([
        window.LearnGoAPI.getModule(moduleId),
        window.LearnGoAPI.getProgress(),
      ]);

      document.getElementById('breadcrumbModule').textContent = mod.title;
      document.getElementById('moduleTitle').textContent = mod.title;
      document.getElementById('lessonBody').innerHTML = marked.parse(mod.lesson || '');

      const lessonItemId = `${moduleId}:lesson`;
      const lessonDone = progress.items[lessonItemId];
      if (lessonDone) {
        document.getElementById('markLessonDone').classList.add('d-none');
        document.getElementById('lessonDoneBadge').classList.remove('d-none');
      }

      document.getElementById('markLessonDone').addEventListener('click', async () => {
        await window.LearnGoAPI.markComplete(lessonItemId);
        document.getElementById('markLessonDone').classList.add('d-none');
        document.getElementById('lessonDoneBadge').classList.remove('d-none');
        const fresh = await window.LearnGoAPI.getProgress();
        updateQuizGate(mod, fresh, true);
      });

      const exerciseList = document.getElementById('exerciseList');
      exerciseList.innerHTML = '';

      if (mod.exercises.length === 0) {
        exerciseList.innerHTML = '<p class="text-muted small">No exercises — complete the project checklist below.</p>';
      }

      mod.exercises.forEach((ex) => {
        const itemId = `exercise:${ex.id}`;
        const done = progress.items[itemId];

        const a = document.createElement('a');
        a.className = 'list-group-item list-group-item-action d-flex justify-content-between align-items-center';
        a.href = lessonDone ? `/exercise.html?id=${encodeURIComponent(ex.id)}&module=${encodeURIComponent(moduleId)}` : '#';
        if (!lessonDone) {
          a.classList.add('disabled');
          a.addEventListener('click', (e) => e.preventDefault());
        }
        a.innerHTML = `<span>${ex.order}. ${ex.title}</span>` +
          (done ? '<span class="badge text-bg-success">Done</span>' : '<span class="badge text-bg-light text-dark">Todo</span>');
        exerciseList.appendChild(a);
      });

      renderProject(mod, progress);
      updateQuizGate(mod, progress, lessonDone);

      loading.classList.add('d-none');
      content.classList.remove('d-none');
    } catch (err) {
      showError(err.message);
    }
  }

  init();
})();
