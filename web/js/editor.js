(function () {
  function initEditor(textarea, starterCode) {
    const theme = document.documentElement.getAttribute('data-bs-theme') === 'dark' ? 'dracula' : 'eclipse';
    const editor = CodeMirror.fromTextArea(textarea, {
      mode: 'go',
      theme,
      lineNumbers: true,
      indentWithTabs: false,
      indentUnit: 4,
      tabSize: 4,
      lineWrapping: false,
    });
    editor.setValue(starterCode || '');
    window.learnGoEditor = editor;
    return editor;
  }

  window.learnGoEditorInit = initEditor;
})();
