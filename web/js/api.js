const API = '/api';

async function apiFetch(path, options = {}) {
  const res = await fetch(API + path, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || res.statusText || 'Request failed');
  }
  return data;
}

window.LearnGoAPI = {
  getCourse: () => apiFetch('/course'),
  getModule: (id) => apiFetch(`/modules/${id}`),
  getExercise: (id) => apiFetch(`/exercises/${id}`),
  submitExercise: (id, code) =>
    apiFetch(`/exercises/${id}/submit`, { method: 'POST', body: JSON.stringify({ code }) }),
  getQuiz: (id) => apiFetch(`/quizzes/${id}`),
  submitQuiz: (id, answers) =>
    apiFetch(`/quizzes/${id}/submit`, { method: 'POST', body: JSON.stringify({ answers }) }),
  getProgress: () => apiFetch('/progress'),
  markComplete: (itemId) =>
    apiFetch('/progress/complete', { method: 'POST', body: JSON.stringify({ itemId }) }),
};
