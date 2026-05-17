export function escapeHtml(value: unknown): string {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

export function formatRelativeTime(iso: string | undefined): string {
  if (!iso) return '';
  const diff = Date.now() - new Date(iso).getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 5) return 'just now';
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return `${Math.floor(hours / 24)}d ago`;
}

export function formatDate(iso: string | undefined): string {
  if (!iso) return '';
  return new Date(iso).toLocaleString();
}

export function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
}

export function methodClass(method: string | undefined): string {
  const classes: Record<string, string> = {
    GET: 'method-GET',
    POST: 'method-POST',
    PUT: 'method-PUT',
    DELETE: 'method-DELETE',
    PATCH: 'method-PATCH',
  };
  return classes[(method || '').toUpperCase()] || 'method-OTHER';
}

export function isBase64(value: unknown): value is string {
  if (!value || typeof value !== 'string' || value.length % 4 !== 0) return false;
  return /^[A-Za-z0-9+/]*={0,2}$/.test(value);
}

export function looksLikeJSON(value: string | undefined): boolean {
  if (!value) return false;
  const text = value.trim();
  return (
    (text.startsWith('{') && text.endsWith('}')) || (text.startsWith('[') && text.endsWith(']'))
  );
}

export function syntaxHighlightJSON(json: string): string {
  return json.replace(
    /("(\\u[a-fA-F0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+-]?\d+)?)/g,
    (match) => {
      let className = 'json-num';
      if (/^"/.test(match)) {
        className = /:$/.test(match) ? 'json-key' : 'json-str';
      } else if (/true|false/.test(match)) {
        className = 'json-bool';
      } else if (/null/.test(match)) {
        className = 'json-null';
      }
      return `<span class="${className}">${escapeHtml(match)}</span>`;
    },
  );
}
