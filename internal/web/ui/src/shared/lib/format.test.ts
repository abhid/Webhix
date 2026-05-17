import test from 'node:test';
import assert from 'node:assert/strict';

import { escapeHtml, formatBytes, looksLikeJSON, syntaxHighlightJSON } from './format';

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

test('syntaxHighlightJSON escapes JSON string tokens before adding markup', () => {
  const highlighted = syntaxHighlightJSON('{"payload":"<img src=x onerror=\\"alert(1)\\">"}');

  assert.equal(highlighted.includes('<img src=x'), false);
  assert.equal(highlighted.includes('&lt;img src=x'), true);
  assert.equal(highlighted.includes('<span class="json-str">'), true);
});
