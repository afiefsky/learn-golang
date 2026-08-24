(function () {
  const params = new URLSearchParams(window.location.search);
  const exerciseId = params.get('id');
  const moduleId = params.get('module');
  const loading = document.getElementById('loading');
  const content = document.getElementById('exerciseContent');
  const errorBox = document.getElementById('errorBox');
  let starterCode = '';

  if (!exerciseId) {
    showError('Missing exercise id.');
    return;
  }

  function showError(msg) {
    loading.classList.add('d-none');
    errorBox.textContent = msg;
    errorBox.classList.remove('d-none');
  }

  function renderFeedback(result) {
    const box = document.getElementById('feedback');
    if (result.passed) {
      box.innerHTML = `<div class="alert alert-success">${escapeHtml(result.message)}</div>`;
      return;
    }

    let html = `<div class="alert alert-warning">${escapeHtml(result.message)}</div>`;
    if (result.errors && result.errors.length) {
      html += '<ul class="small text-danger">';
      result.errors.forEach((e) => { html += `<li>${escapeHtml(e)}</li>`; });
      html += '</ul>';
    }
    if (result.hint) {
      html += `<div class="alert alert-info"><strong>Hint:</strong> ${escapeHtml(result.hint)}</div>`;
    }
    box.innerHTML = html;
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  async function init() {
    try {
      const ex = await window.LearnGoAPI.getExercise(exerciseId);
      document.title = `${ex.title} — Learn Go`;
      document.getElementById('exerciseTitle').textContent = ex.title;
      document.getElementById('exerciseBreadcrumb').textContent = ex.title;
      document.getElementById('exerciseNarrative').textContent = ex.narrative;
      if (moduleId) {
        document.getElementById('moduleLink').href = `/lesson.html?module=${encodeURIComponent(moduleId)}`;
      }

      starterCode = ex.starterCode || '';
      window.learnGoEditorInit(document.getElementById('codeEditor'), starterCode);

      document.getElementById('resetBtn').addEventListener('click', () => {
        window.learnGoEditor.setValue(starterCode);
        document.getElementById('feedback').innerHTML = '';
      });

      document.getElementById('submitBtn').addEventListener('click', async () => {
        const btn = document.getElementById('submitBtn');
        btn.disabled = true;
        try {
          const code = window.learnGoEditor.getValue();
          const result = await window.LearnGoAPI.submitExercise(exerciseId, code);
          renderFeedback(result);
        } catch (err) {
          renderFeedback({ passed: false, message: err.message });
        } finally {
          btn.disabled = false;
        }
      });

      loading.classList.add('d-none');
      content.classList.remove('d-none');
    } catch (err) {
      showError(err.message);
    }
  }

  init();
})();
