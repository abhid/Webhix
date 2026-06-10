import type { Elements } from '../../app/dom';
import { selectedRequest, type AppState } from '../../entities/request/model/request-state';
import type { RequestTab, WebhookRequest } from '../../entities/request/model/types';
import {
  formatBytes,
  formatDate,
  isBase64,
  looksLikeJSON,
  methodClass,
  relativeWebhookPath,
  syntaxHighlightJSON,
} from '../../shared/lib/format';

export function renderSelectedDetail(elements: Elements, state: AppState): void {
  const request = selectedRequest(state);
  if (!request) {
    showPlaceholder(elements);
    return;
  }

  elements.detailPlaceholder.classList.add('hidden');
  elements.detailContent.classList.remove('hidden');

  elements.detailMethod.textContent = request.method;
  elements.detailMethod.className = `method-badge ${methodClass(request.method)}`;
  elements.detailPath.textContent = relativeWebhookPath(request.path, state.token);
  elements.detailPath.title = request.path;
  elements.detailTimestamp.textContent = formatDate(request.receivedAt);

  renderActiveTab(elements, request, state.activeTab);
  renderTabButtons(elements, state.activeTab);
}

export function showPlaceholder(elements: Elements): void {
  elements.detailPlaceholder.classList.remove('hidden');
  elements.detailContent.classList.add('hidden');
}

function renderActiveTab(elements: Elements, request: WebhookRequest, tab: RequestTab): void {
  elements.tabContent.replaceChildren();

  switch (tab) {
    case 'body':
      elements.tabContent.appendChild(createBody(request.body, request.contentType));
      break;
    case 'headers':
      elements.tabContent.appendChild(createKeyValueTable(parseHeaders(request.headers)));
      break;
    case 'details':
      elements.tabContent.appendChild(createDetails(request));
      break;
  }
}

function renderTabButtons(elements: Elements, activeTab: RequestTab): void {
  for (const button of elements.tabButtons) {
    button.classList.toggle('active', button.dataset.tab === activeTab);
  }
}

function createKeyValueTable(pairs: Array<[string, string]>): HTMLElement {
  if (pairs.length === 0) {
    const message = document.createElement('p');
    message.className = 'empty-table-message';
    message.textContent = 'No entries';
    return message;
  }

  const table = document.createElement('table');
  table.className = 'kv-table';

  const thead = document.createElement('thead');
  const headerRow = document.createElement('tr');
  const keyHeader = document.createElement('th');
  keyHeader.textContent = 'Key';
  const valueHeader = document.createElement('th');
  valueHeader.textContent = 'Value';
  headerRow.append(keyHeader, valueHeader);
  thead.appendChild(headerRow);

  const tbody = document.createElement('tbody');
  for (const [key, value] of pairs) {
    const row = document.createElement('tr');
    const keyCell = document.createElement('td');
    keyCell.textContent = key;
    const valueCell = document.createElement('td');
    valueCell.textContent = value;
    row.append(keyCell, valueCell);
    tbody.appendChild(row);
  }

  table.append(thead, tbody);
  return table;
}

function createBody(body: unknown, contentType: string | undefined): HTMLPreElement {
  const pre = document.createElement('pre');
  pre.className = 'body-pre';

  if (!body || (typeof body === 'string' && body.trim() === '') || isEmptyArray(body)) {
    pre.classList.add('body-empty');
    pre.textContent = '(empty body)';
    return pre;
  }

  let text: string;
  if (typeof body === 'string') {
    text = body;
  } else {
    text = String(body);
  }

  if (isBase64(text)) {
    try {
      text = atob(text);
    } catch {
      // Keep original body text if base64 decoding fails.
    }
  }

  const contentTypeValue = (contentType || '').toLowerCase();
  if (contentTypeValue.includes('json') || looksLikeJSON(text)) {
    try {
      const parsed = JSON.parse(text) as unknown;
      // syntaxHighlightJSON escapes every token before wrapping it in controlled span markup.
      pre.innerHTML = syntaxHighlightJSON(JSON.stringify(parsed, null, 2));
      return pre;
    } catch {
      // Fall through to plain text rendering.
    }
  }

  pre.textContent = text;
  return pre;
}

function createDetails(request: WebhookRequest): HTMLDivElement {
  const wrap = document.createElement('div');
  wrap.className = 'details-tab';

  const query = parseQuery(request.query);
  if (query.length > 0) {
    const heading = document.createElement('h4');
    heading.className = 'details-heading';
    heading.textContent = 'Query Parameters';
    wrap.append(heading, createKeyValueTable(query));
  }

  const metaHeading = document.createElement('h4');
  metaHeading.className = 'details-heading';
  metaHeading.textContent = 'Metadata';
  wrap.append(metaHeading, createInfo(request));

  return wrap;
}

function createInfo(request: WebhookRequest): HTMLDivElement {
  const size = request.bodySize != null ? formatBytes(request.bodySize) : '0 B';
  const grid = document.createElement('div');
  grid.className = 'info-grid';
  grid.append(
    createInfoCard('Remote IP', request.remoteAddr || '—'),
    createInfoCard('Content-Type', request.contentType || '—'),
    createInfoCard('Body Size', size),
    createInfoCard('Received At', formatDate(request.receivedAt)),
    createInfoCard('Hook ID', String(request.hookId || '—')),
    createInfoCard('Request ID', String(request.id || '—')),
  );
  return grid;
}

function createInfoCard(label: string, value: string): HTMLDivElement {
  const card = document.createElement('div');
  card.className = 'info-card';

  const cardLabel = document.createElement('div');
  cardLabel.className = 'card-label';
  cardLabel.textContent = label;

  const cardValue = document.createElement('div');
  cardValue.className = 'card-value';
  cardValue.textContent = value;

  card.append(cardLabel, cardValue);
  return card;
}

function parseHeaders(raw: string | undefined): Array<[string, string]> {
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw) as Record<string, unknown>;
    const pairs: Array<[string, string]> = [];
    for (const [key, value] of Object.entries(parsed)) {
      if (Array.isArray(value)) value.forEach((item) => pairs.push([key, String(item)]));
      else pairs.push([key, String(value)]);
    }
    return pairs;
  } catch {
    return [];
  }
}

function parseQuery(raw: string | undefined): Array<[string, string]> {
  if (!raw) return [];
  try {
    return Array.from(new URLSearchParams(raw).entries());
  } catch {
    return [];
  }
}

function isEmptyArray(value: unknown): boolean {
  return Array.isArray(value) && value.length === 0;
}
