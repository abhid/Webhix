import type { ApiResponse, WebhookRequest } from '../../../entities/request/model/types';

export async function createEndpoint(): Promise<string> {
  const response = await fetch('/api/endpoints', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: '' }),
  });
  const json = (await response.json()) as ApiResponse<{ token: string }>;
  if (!json.success) throw new Error(json.error?.message || 'Failed');
  if (!json.body?.token) throw new Error('Endpoint response did not include a token');
  return json.body.token;
}

export async function fetchRequests(token: string): Promise<WebhookRequest[]> {
  const response = await fetch(`/api/endpoints/${token}/requests`);
  const json = (await response.json()) as ApiResponse<WebhookRequest[]>;
  if (!json.success) return [];
  return json.body || [];
}

export function connectEvents(
  token: string,
  handlers: {
    onOpen?: () => void;
    onRequest?: (request: WebhookRequest) => void;
    onError?: () => void;
  },
): EventSource {
  const source = new EventSource(`/api/endpoints/${token}/events`);

  source.onopen = () => handlers.onOpen?.();
  source.onmessage = (event) => {
    try {
      handlers.onRequest?.(JSON.parse(event.data) as WebhookRequest);
    } catch {
      // Ignore malformed SSE payloads; the stream can continue.
    }
  };
  source.onerror = () => handlers.onError?.();

  return source;
}
