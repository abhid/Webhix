import './styles.css';

import {
  connectEvents,
  createEndpoint,
  fetchRequests,
} from '../features/endpoint-session/api/endpoint-api';
import { getElements } from './dom';
import { renderRequestList, refreshRelativeTimes } from '../widgets/request-list/request-list';
import { renderSelectedDetail, showPlaceholder } from '../widgets/request-detail/request-detail';
import {
  addRequests,
  createInitialState,
  resetForToken,
  selectRequest,
  setActiveTab,
} from '../entities/request/model/request-state';
import type { RequestTab } from '../entities/request/model/types';

const state = createInitialState();
const elements = getElements();
let eventSource: EventSource | null = null;
let toastTimer: ReturnType<typeof setTimeout> | undefined;

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
  elements.newEndpointButton.addEventListener('click', showOverlay);
  elements.loadTokenButton.addEventListener('click', loadToken);
  elements.createEndpointButton.addEventListener('click', createNewEndpoint);
  elements.requestList.addEventListener('click', handleRequestClick);

  for (const button of elements.tabButtons) {
    button.addEventListener('click', () => {
      if (isRequestTab(button.dataset.tab)) switchTab(button.dataset.tab);
    });
  }

  setInterval(() => refreshRelativeTimes(elements.requestList), 15000);
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
  resetForToken(state, token);

  const url = new URL(location.href);
  url.searchParams.set('token', token);
  history.replaceState(null, '', url.toString());

  elements.pillURL.textContent = `${location.origin}/r/${token}`;
  document.getElementById('endpointName')?.replaceChildren(document.createTextNode(token));
  document.getElementById('endpointTitle')?.replaceChildren(document.createTextNode(token));
  document.getElementById('endpointPath')?.replaceChildren(document.createTextNode(`/r/${token}`));
  elements.pillArea.classList.remove('hidden');
  elements.overlay.classList.add('hidden');
  elements.mainArea.classList.remove('hidden');

  renderRequestList(elements, state);
  showPlaceholder(elements);
  void loadHistory(token);
  connectSSE(token);
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
    addRequests(state, requests.reverse());
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
  return value === 'headers' || value === 'body' || value === 'query' || value === 'info';
}
