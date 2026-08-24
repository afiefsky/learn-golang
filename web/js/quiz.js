(function () {
  const params = new URLSearchParams(window.location.search);
  const quizId = params.get('id');
  const moduleId = params.get('module');
  const loading = document.getElementById('loading');
  const content = document.getElementById('quizContent');
  const errorBox = document.getElementById('errorBox');

  if (!quizId) {
    showError('Missing quiz id.');
    return;
  }

  function showError(msg) {
    loading.classList.add('d-none');
    errorBox.textContent = msg;
    errorBox.classList.remove('d-none');
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  function buildForm(questions) {
    const form = document.getElementById('quizForm');
    form.innerHTML = '';

    questions.forEach((q, i) => {
      const block = document.createElement('fieldset');
      block.className = 'mb-4';
      block.innerHTML = `<legend class="fs-6 fw-semibold">${i + 1}. ${escapeHtml(q.question)}</legend>`;

      if (q.type === 'mcq') {
        q.options.forEach((opt, j) => {
          const id = `q${i}_o${j}`;
          block.innerHTML += `
            <div class="form-check">
              <input class="form-check-input" type="radio" name="q${i}" id="${id}" value="${j}" required>
              <label class="form-check-label" for="${id}">${escapeHtml(opt)}</label>
            </div>`;
        });
      } else if (q.type === 'fill') {
        block.innerHTML += `
          <input type="text" class="form-control" name="q${i}" required autocomplete="off"
                 placeholder="Your answer">`;
      }

      form.appendChild(block);
    });
  }

  function collectAnswers(questions) {
    const answers = [];
    questions.forEach((q, i) => {
      if (q.type === 'mcq') {
        const selected = form.querySelector(`input[name="q${i}"]:checked`);
        answers.push(selected ? parseInt(selected.value, 10) : '');
      } else {
        const input = form.querySelector(`input[name="q${i}"]`);
        answers.push(input ? input.value : '');
      }
    });
    return answers;
  }

  let form;
  let questions = [];

  function renderResult(result) {
    const box = document.getElementById('quizResult');
    const cls = result.passed ? 'success' : 'warning';
    let html = `<div class="alert alert-${cls}">
      Score: ${result.score}% (${result.review.filter(r => r.correct).length}/${result.total} correct)
      — need ${result.minScore}% to pass.
    </div>`;

    html += '<div class="list-group">';
    result.review.forEach((r, i) => {
      const badge = r.correct ? 'success' : 'danger';
      html += `<div class="list-group-item">
        <div class="d-flex justify-content-between">
          <span>${i + 1}. ${escapeHtml(r.question)}</span>
          <span class="badge text-bg-${badge}">${r.correct ? 'Correct' : 'Wrong'}</span>
        </div>`;
      if (!r.correct) {
        html += `<small class="text-muted">Expected: ${escapeHtml(r.expected)} · You: ${escapeHtml(r.given)}</small>`;
      }
      html += '</div>';
    });
    html += '</div>';

    if (!result.passed) {
      html += '<button type="button" class="btn btn-outline-primary mt-3" id="retryQuiz">Try again</button>';
    } else if (moduleId) {
      html += `<a href="/lesson.html?module=${encodeURIComponent(moduleId)}" class="btn btn-success mt-3">Back to module</a>`;
      html += ` <a href="/" class="btn btn-outline-secondary mt-3">Course home</a>`;
    }

    box.innerHTML = html;

    const retry = document.getElementById('retryQuiz');
    if (retry) {
      retry.addEventListener('click', () => {
        box.innerHTML = '';
        document.getElementById('submitQuiz').classList.remove('d-none');
      });
    }
  }

  async function init() {
    try {
      const quiz = await window.LearnGoAPI.getQuiz(quizId);
      questions = quiz.questions;
      document.getElementById('quizTitle').textContent = quiz.title;
      if (moduleId) {
        document.getElementById('moduleLink').href = `/lesson.html?module=${encodeURIComponent(moduleId)}`;
      }

      form = document.getElementById('quizForm');
      buildForm(questions);

      document.getElementById('submitQuiz').addEventListener('click', async () => {
        if (!form.checkValidity()) {
          form.reportValidity();
          return;
        }
        const btn = document.getElementById('submitQuiz');
        btn.disabled = true;
        try {
          const answers = collectAnswers(questions);
          const result = await window.LearnGoAPI.submitQuiz(quizId, answers);
          renderResult(result);
          btn.classList.add('d-none');
        } catch (err) {
          document.getElementById('quizResult').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
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
