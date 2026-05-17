import test from 'node:test';
import assert from 'node:assert/strict';

import { addRequests, createInitialState, resetForToken, selectRequest } from './request-state';

test('addRequests prepends live requests and ignores duplicate ids', () => {
  const state = createInitialState();

  addRequests(state, [{ id: 1, method: 'POST', path: '/old', receivedAt: '2026-05-16T10:00:00Z' }]);
  addRequests(
    state,
    [{ id: 2, method: 'POST', path: '/new', receivedAt: '2026-05-16T11:00:00Z' }],
    { prepend: true },
  );
  addRequests(
    state,
    [{ id: 2, method: 'POST', path: '/new', receivedAt: '2026-05-16T11:00:00Z' }],
    { prepend: true },
  );

  assert.deepEqual(
    state.requests.map((request) => request.id),
    [2, 1],
  );
});

test('selectRequest stores request id instead of list index', () => {
  const state = createInitialState();
  addRequests(state, [
    { id: 1, method: 'POST', path: '/first', receivedAt: '2026-05-16T10:00:00Z' },
    { id: 2, method: 'POST', path: '/second', receivedAt: '2026-05-16T11:00:00Z' },
  ]);

  selectRequest(state, '2');

  assert.equal(state.selectedRequestId, '2');
});

test('resetForToken clears request-specific state and keeps the active tab', () => {
  const state = createInitialState();
  state.activeTab = 'body';
  addRequests(state, [{ id: 1, method: 'POST', path: '/old', receivedAt: '2026-05-16T10:00:00Z' }]);
  selectRequest(state, '1');

  resetForToken(state, 'token-123');

  assert.equal(state.token, 'token-123');
  assert.equal(state.activeTab, 'body');
  assert.equal(state.selectedRequestId, null);
  assert.deepEqual(state.requests, []);
});
