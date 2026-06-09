import test from 'node:test';
import assert from 'node:assert/strict';

import {
  escapeHtml,
  formatBytes,
  looksLikeJSON,
  relativeWebhookPath,
  syntaxHighlightJSON,
} from './format';

test('escapeHtml escapes HTML-sensitive characters', () => {
  assert.equal(
    escapeHtml('<img src=x onerror="alert(1)">&'),
    '&lt;img src=x onerror=&quot;alert(1)&quot;&gt;&amp;',
  );
});

test('formatBytes formats bytes using compact units', () => {
  assert.equal(formatBytes(512), '512 B');
  assert.equal(formatBytes(1536), '1.5 KB');
  assert.equal(formatBytes(1048576), '1.00 MB');
});

test('looksLikeJSON only accepts object and array shaped strings', () => {
  assert.equal(looksLikeJSON('{"ok":true}'), true);
  assert.equal(looksLikeJSON('[1,2,3]'), true);
  assert.equal(looksLikeJSON('plain text'), false);
});

test('relativeWebhookPath strips the base webhook prefix', () => {
  assert.equal(relativeWebhookPath('/r/abc123', 'abc123'), '/');
  assert.equal(relativeWebhookPath('/r/abc123/github/push', 'abc123'), '/github/push');
  assert.equal(relativeWebhookPath('/r/abc123?ref=main', 'abc123'), '?ref=main');
  assert.equal(relativeWebhookPath('/r/other/path', 'abc123'), '/r/other/path');
  assert.equal(relativeWebhookPath('/r/abc123', null), '/r/abc123');
});

test('syntaxHighlightJSON escapes JSON string tokens before adding markup', () => {
  const highlighted = syntaxHighlightJSON('{"payload":"<img src=x onerror=\\"alert(1)\\">"}');

  assert.equal(highlighted.includes('<img src=x'), false);
  assert.equal(highlighted.includes('&lt;img src=x'), true);
  assert.equal(highlighted.includes('<span class="json-str">'), true);
});
