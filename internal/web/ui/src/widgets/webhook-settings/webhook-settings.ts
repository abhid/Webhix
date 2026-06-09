import {
  fetchHookResponse,
  saveHookResponse,
} from '../../features/endpoint-session/api/endpoint-api';

export function renderWebhookSettings(container: HTMLElement, token: string | null): void {
  container.replaceChildren();
  if (!token) return;

  const intro = document.createElement('p');
  intro.className = 'settings-intro';
  intro.textContent = 'Configure the response Webhix returns to senders of this endpoint.';

  const form = document.createElement('div');
  form.className = 'settings-form';

  const statusLabel = document.createElement('label');
  statusLabel.textContent = 'Response Status Code';
  const statusInput = document.createElement('input');
  statusInput.type = 'number';
  statusInput.min = '100';
  statusInput.max = '599';
  statusInput.value = '200';
  statusInput.className = 'settings-input';

  const headersLabel = document.createElement('label');
  headersLabel.textContent = 'Response Headers (JSON)';
  const headersInput = document.createElement('textarea');
  headersInput.className = 'settings-textarea';
  headersInput.placeholder = '{"Content-Type": "application/json"}';
  headersInput.rows = 3;

  const bodyLabel = document.createElement('label');
  bodyLabel.textContent = 'Response Body';
  const bodyInput = document.createElement('textarea');
  bodyInput.className = 'settings-textarea';
  bodyInput.placeholder = '{"ok": true}';
  bodyInput.rows = 5;

  const saveBtn = document.createElement('button');
  saveBtn.className = 'settings-save-btn';
  saveBtn.textContent = 'Save';

  form.append(statusLabel, statusInput, headersLabel, headersInput, bodyLabel, bodyInput, saveBtn);
  container.append(intro, form);

  void fetchHookResponse(token).then((resp) => {
    statusInput.value = String(resp.statusCode || 200);
    headersInput.value = Object.keys(resp.headers || {}).length
      ? JSON.stringify(resp.headers, null, 2)
      : '';
    bodyInput.value = resp.body || '';
  });

  saveBtn.addEventListener('click', () => {
    let headers: Record<string, string> = {};
    try {
      if (headersInput.value.trim()) {
        headers = JSON.parse(headersInput.value) as Record<string, string>;
      }
    } catch {
      saveBtn.textContent = 'Invalid JSON in headers';
      setTimeout(() => (saveBtn.textContent = 'Save'), 2000);
      return;
    }

    saveBtn.disabled = true;
    void saveHookResponse(token, {
      statusCode: parseInt(statusInput.value, 10) || 200,
      headers,
      body: bodyInput.value,
    })
      .then(() => {
        saveBtn.textContent = 'Saved!';
        setTimeout(() => (saveBtn.textContent = 'Save'), 2000);
      })
      .catch(() => {
        saveBtn.textContent = 'Error';
        setTimeout(() => (saveBtn.textContent = 'Save'), 2000);
      })
      .finally(() => {
        saveBtn.disabled = false;
      });
  });
}
