export type RequestTab = 'body' | 'headers' | 'details';

export interface WebhookRequest {
  id: number;
  method: string;
  path: string;
  receivedAt: string;
  headers?: string;
  query?: string;
  body?: unknown;
  contentType?: string;
  bodySize?: number;
  remoteAddr?: string;
  hookId?: string | number;
}

export interface ApiResponse<T> {
  success: boolean;
  body?: T;
  error?: {
    message?: string;
  };
}
