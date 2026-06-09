import type { RequestTab, WebhookRequest } from './types';

export interface AppState {
  token: string | null;
  requests: WebhookRequest[];
  selectedRequestId: string | null;
  activeTab: RequestTab;
  seenIds: Set<number>;
  searchQuery: string;
  methodFilter: string | null;
}

export function createInitialState(): AppState {
  return {
    token: null,
    requests: [],
    selectedRequestId: null,
    activeTab: 'body',
    seenIds: new Set<number>(),
    searchQuery: '',
    methodFilter: null,
  };
}

export function resetForToken(state: AppState, token: string): void {
  state.token = token;
  state.requests = [];
  state.selectedRequestId = null;
  state.seenIds = new Set<number>();
  state.searchQuery = '';
  state.methodFilter = null;
}

export function addRequests(
  state: AppState,
  requests: WebhookRequest[],
  options: { prepend?: boolean } = {},
): WebhookRequest[] {
  const nextRequests: WebhookRequest[] = [];
  for (const request of requests) {
    if (!request.id || state.seenIds.has(request.id)) continue;
    state.seenIds.add(request.id);
    nextRequests.push(request);
  }

  if (options.prepend) {
    state.requests.unshift(...nextRequests);
  } else {
    state.requests.push(...nextRequests);
  }

  return nextRequests;
}

export function selectRequest(state: AppState, requestId: string): void {
  state.selectedRequestId = requestId;
}

export function selectedRequest(state: AppState): WebhookRequest | null {
  return state.requests.find((request) => String(request.id) === state.selectedRequestId) || null;
}

export function setActiveTab(state: AppState, tab: RequestTab): void {
  state.activeTab = tab;
}

export function setSearchQuery(state: AppState, query: string): void {
  state.searchQuery = query;
}

export function setMethodFilter(state: AppState, method: string | null): void {
  state.methodFilter = method;
}

export function filteredRequests(state: AppState): WebhookRequest[] {
  const query = state.searchQuery.trim().toLowerCase();
  return state.requests.filter((request) => {
    if (state.methodFilter && request.method !== state.methodFilter) return false;
    if (!query) return true;
    return (
      request.path.toLowerCase().includes(query) ||
      (request.headers ?? '').toLowerCase().includes(query) ||
      (typeof request.body === 'string' ? atob(request.body) : '').toLowerCase().includes(query)
    );
  });
}

export function uniqueMethods(state: AppState): string[] {
  return [...new Set(state.requests.map((r) => r.method))].sort();
}
