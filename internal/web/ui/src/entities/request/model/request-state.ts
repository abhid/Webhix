import type { RequestTab, WebhookRequest } from './types';

export interface AppState {
  token: string | null;
  requests: WebhookRequest[];
  selectedRequestId: string | null;
  activeTab: RequestTab;
  seenIds: Set<string>;
}

export function createInitialState(): AppState {
  return {
    token: null,
    requests: [],
    selectedRequestId: null,
    activeTab: 'headers',
    seenIds: new Set(),
  };
}

export function resetForToken(state: AppState, token: string): void {
  state.token = token;
  state.requests = [];
  state.selectedRequestId = null;
  state.seenIds = new Set();
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
  return state.requests.find((request) => request.id === state.selectedRequestId) || null;
}

export function setActiveTab(state: AppState, tab: RequestTab): void {
  state.activeTab = tab;
}
