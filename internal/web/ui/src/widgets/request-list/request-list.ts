import type { Elements } from '../../app/dom';
import { filteredRequests, type AppState } from '../../entities/request/model/request-state';
import { formatRelativeTime, methodClass, relativeWebhookPath } from '../../shared/lib/format';

export function renderRequestList(
  elements: Elements,
  state: AppState,
  options: { highlightRequestId?: number; scrollTop?: boolean } = {},
): void {
  const visible = filteredRequests(state);

  elements.countBadge.textContent = String(state.requests.length);
  updateEndpointStats(elements, state);

  if (visible.length === 0) {
    elements.requestList.replaceChildren();
    elements.requestList.appendChild(elements.emptyState);
    return;
  }

  elements.emptyState.remove();
  elements.requestList.replaceChildren();

  for (const request of visible) {
    const item = document.createElement('button');
    item.type = 'button';
    item.className = `request-item${String(request.id) === state.selectedRequestId ? ' active' : ''}`;
    item.dataset.requestId = String(request.id);

    const method = document.createElement('span');
    method.className = `method-badge ${methodClass(request.method)}`;
    method.textContent = request.method;

    const meta = document.createElement('span');
    meta.className = 'request-meta';

    const path = document.createElement('span');
    path.className = 'request-path';
    path.title = request.path;
    path.textContent = relativeWebhookPath(request.path, state.token);

    const time = document.createElement('span');
    time.className = 'request-time';
    time.dataset.receivedAt = request.receivedAt || '';
    time.textContent = formatRelativeTime(request.receivedAt);

    meta.append(path, time);
    item.append(method, meta);

    if (options.highlightRequestId === request.id) {
      const dot = document.createElement('span');
      dot.className = 'new-dot';
      item.appendChild(dot);
    }

    elements.requestList.appendChild(item);
  }

  if (options.scrollTop) elements.requestList.scrollTop = 0;
}

export function refreshRelativeTimes(root: ParentNode = document): void {
  root.querySelectorAll<HTMLElement>('.request-time').forEach((element) => {
    element.textContent = formatRelativeTime(element.dataset.receivedAt);
  });
}

function updateEndpointStats(elements: Elements, state: AppState): void {
  const startOfToday = new Date();
  startOfToday.setHours(0, 0, 0, 0);

  let today = 0;
  let latest: string | undefined;
  for (const request of state.requests) {
    if (!request.receivedAt) continue;
    const received = new Date(request.receivedAt);
    if (received.getTime() >= startOfToday.getTime()) today += 1;
    if (!latest || received.getTime() > new Date(latest).getTime()) latest = request.receivedAt;
  }

  elements.statToday.textContent = String(today);
  elements.statLast.textContent = latest ? formatRelativeTime(latest) : 'Waiting';
}
