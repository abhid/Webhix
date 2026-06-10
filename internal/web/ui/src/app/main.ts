import '@fontsource-variable/inter';
import '@fontsource/monaspace-neon/400.css';
import '@fontsource/monaspace-neon/500.css';
import './styles.css';

import {
  connectEvents,
  createEndpoint,
  fetchEndpoints,
  fetchRequests,
} from '../features/endpoint-session/api/endpoint-api';
import type { Endpoint } from '../features/endpoint-session/api/endpoint-api';
import { getElements } from './dom';
import { renderRequestList, refreshRelativeTimes } from '../widgets/request-list/request-list';
import { renderSelectedDetail, showPlaceholder } from '../widgets/request-detail/request-detail';
import { renderWebhookSettings } from '../widgets/webhook-settings/webhook-settings';
import {
  addRequests,
  createInitialState,
  resetForToken,
  selectRequest,
  selectedRequest,
  setActiveTab,
  setMethodFilter,
  setSearchQuery,
  uniqueMethods,
} from '../entities/request/model/request-state';
import type { RequestTab } from '../entities/request/model/types';
import { buildCurlCommand } from '../shared/lib/format';

const state = createInitialState();
const elements = getElements();
let eventSource: EventSource | null = null;
let toastTimer: ReturnType<typeof setTimeout> | undefined;
let currentToken: string | null = null;

init();

function init(): void {
  const params = new URLSearchParams(location.search);
  const token = params.get('token');
  if (token) {
    elements.tokenInput.value = token;
    activateToken(token);
  }

  elements.tokenInput.addEventListener('keydown', (event) => {
    if (event.key === 'Enter') loadToken();
  });
  elements.copyButton.addEventListener('click', copyURL);
  elements.curlButton.addEventListener('click', copyCurl);
  elements.replayButton.addEventListener('click', replayRequest);
  elements.newEndpointButton.addEventListener('click', showOverlay);
  elements.loadTokenButton.addEventListener('click', loadToken);
  elements.createEndpointButton.addEventListener('click', createNewEndpoint);
  elements.requestList.addEventListener('click', handleRequestClick);
  elements.searchInput.addEventListener('input', () => {
    setSearchQuery(state, elements.searchInput.value);
    renderRequestList(elements, state);
  });
  elements.methodFilterButton.addEventListener('click', () => {
    const methods = uniqueMethods(state);
    const current = state.methodFilter;
    const idx = current ? methods.indexOf(current) : -1;
    const next = methods[idx + 1] ?? null;
    setMethodFilter(state, next);
    elements.methodFilterButton.textContent = next ?? 'All Methods';
    renderRequestList(elements, state);
  });

  for (const button of elements.tabButtons) {
    button.addEventListener('click', () => {
      if (isRequestTab(button.dataset.tab)) switchTab(button.dataset.tab);
    });
  }

  for (const button of elements.sectionTabs) {
    button.addEventListener('click', () => switchSection(button.dataset.section));
  }

  setInterval(() => refreshRelativeTimes(elements.requestList), 15000);

  void loadEndpoints();
}

function switchSection(section: string | undefined): void {
  const isSettings = section === 'settings';
  for (const button of elements.sectionTabs) {
    button.classList.toggle('active', button.dataset.section === section);
  }
  elements.requestsSection.classList.toggle('hidden', isSettings);
  elements.settingsSection.classList.toggle('hidden', !isSettings);
  if (isSettings) renderWebhookSettings(elements.webhookSettings, currentToken);
}

function showOverlay(): void {
  elements.overlay.classList.remove('hidden');
  elements.mainArea.classList.add('hidden');
}

function loadToken(): void {
  const token = elements.tokenInput.value.trim();
  if (!token) return;
  activateToken(token);
}

function activateToken(token: string): void {
  currentToken = token;
  resetForToken(state, token);

  const url = new URL(location.href);
  url.searchParams.set('token', token);
  history.replaceState(null, '', url.toString());

  elements.pillURL.textContent = `${location.origin}/r/${token}`;
  document.getElementById('endpointTitle')?.replaceChildren(document.createTextNode(token));
  elements.pillArea.classList.remove('hidden');
  elements.overlay.classList.add('hidden');
  elements.mainArea.classList.remove('hidden');
  switchSection('requests');

  renderRequestList(elements, state);
  showPlaceholder(elements);
  void loadHistory(token);
  void loadEndpoints();
  connectSSE(token);
}

async function loadEndpoints(): Promise<void> {
  try {
    const eps = await fetchEndpoints();
    renderEndpointsList(eps);
  } catch {
    // silently fail
  }
}

function renderEndpointsList(endpoints: Endpoint[]): void {
  const list = document.getElementById('endpointsList');
  const countBadge = document.getElementById('endpointsPanelCount');
  if (!list) return;

  if (countBadge) countBadge.textContent = String(endpoints.length);

  elements.overviewEndpoints.textContent = String(endpoints.length);
  const totalRequests = endpoints.reduce((sum, ep) => sum + (ep.requestCount ?? 0), 0);
  elements.overviewRequests.textContent = totalRequests.toLocaleString();

  list.replaceChildren();
  for (const ep of endpoints) {
    const btn = document.createElement('button');
    btn.className = 'endpoint-card' + (ep.token === currentToken ? ' active' : '');
    btn.dataset.token = ep.token;

    const label = document.createElement('strong');
    label.textContent = ep.name || ep.token;

    const count = document.createElement('span');
    count.className = 'endpoint-count';
    const received = ep.requestCount ?? 0;
    count.textContent = `${received} ${received === 1 ? 'event' : 'events'}`;

    const path = document.createElement('small');
    path.textContent = `/r/${ep.token}`;

    btn.append(label, count, path);
    btn.addEventListener('click', () => activateToken(ep.token));
    list.appendChild(btn);
  }
}

async function createNewEndpoint(): Promise<void> {
  try {
    const token = await createEndpoint();
    elements.tokenInput.value = token;
    activateToken(token);
  } catch (error) {
    toast(`Error: ${error instanceof Error ? error.message : String(error)}`);
  }
}

async function loadHistory(token: string): Promise<void> {
  try {
    const requests = await fetchRequests(token);
    addRequests(state, requests);
    renderRequestList(elements, state);
  } catch {
    toast('Failed to load request history');
  }
}

function connectSSE(token: string): void {
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }

  elements.statusDot.classList.remove('connected');
  eventSource = connectEvents(token, {
    onOpen: () => elements.statusDot.classList.add('connected'),
    onRequest: (request) => {
      const [added] = addRequests(state, [request], { prepend: true });
      if (!added) return;
      renderRequestList(elements, state, { highlightRequestId: added.id, scrollTop: true });
      renderSelectedDetail(elements, state);
    },
    onError: () => elements.statusDot.classList.remove('connected'),
  });
}

function handleRequestClick(event: MouseEvent): void {
  if (!(event.target instanceof Element)) return;
  const item = event.target.closest<HTMLButtonElement>('.request-item');
  if (!item?.dataset.requestId) return;
  selectRequest(state, item.dataset.requestId);
  renderRequestList(elements, state);
  renderSelectedDetail(elements, state);
}

function switchTab(tab: RequestTab): void {
  setActiveTab(state, tab);
  renderSelectedDetail(elements, state);
}

function replayRequest(): void {
  const request = selectedRequest(state);
  if (!request) return;

  const headers: Record<string, string> = {};
  if (request.headers) {
    try {
      const parsed = JSON.parse(request.headers) as Record<string, string | string[]>;
      for (const [key, value] of Object.entries(parsed)) {
        headers[key] = Array.isArray(value) ? (value[0] ?? '') : value;
      }
    } catch {
      // ignore
    }
  }

  let body: BodyInit | undefined;
  if (request.body && typeof request.body === 'string') {
    body = Uint8Array.from(atob(request.body), (c) => c.charCodeAt(0));
  }

  elements.replayButton.disabled = true;
  fetch(`${location.origin}${request.path}`, { method: request.method, headers, body })
    .then(() => toast('Replayed!'))
    .catch(() => toast('Replay failed'))
    .finally(() => {
      elements.replayButton.disabled = false;
    });
}

function copyCurl(): void {
  const request = selectedRequest(state);
  if (!request) return;
  const curl = buildCurlCommand(request.method, request.path, request.headers, request.body);
  navigator.clipboard
    .writeText(curl)
    .then(() => toast('Copied as curl!'))
    .catch(() => toast('Copy failed'));
}

function copyURL(): void {
  const url = elements.pillURL.textContent || '';
  navigator.clipboard
    .writeText(url)
    .then(() => toast('URL copied!'))
    .catch(() => toast('Copy failed'));
}

function toast(message: string): void {
  elements.toast.textContent = message;
  elements.toast.classList.add('show');
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => elements.toast.classList.remove('show'), 2500);
}

function isRequestTab(value: string | undefined): value is RequestTab {
  return value === 'body' || value === 'headers' || value === 'details';
}
