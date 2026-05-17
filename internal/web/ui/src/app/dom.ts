function requireElement<T extends Element>(id: string, constructor: { new (): T }): T {
  const element = document.getElementById(id);
  if (!(element instanceof constructor)) {
    throw new Error(`Expected #${id} to be a ${constructor.name}`);
  }
  return element;
}

export function getElements() {
  return {
    overlay: requireElement('overlay', HTMLDivElement),
    mainArea: requireElement('mainArea', HTMLDivElement),
    tokenInput: requireElement('tokenInput', HTMLInputElement),
    pillArea: requireElement('pillArea', HTMLDivElement),
    pillURL: requireElement('pillURL', HTMLSpanElement),
    statusDot: requireElement('statusDot', HTMLElement),
    requestList: requireElement('requestList', HTMLDivElement),
    emptyState: requireElement('emptyState', HTMLDivElement),
    countBadge: requireElement('countBadge', HTMLElement),
    detailPlaceholder: requireElement('detailPlaceholder', HTMLDivElement),
    detailContent: requireElement('detailContent', HTMLDivElement),
    detailMethod: requireElement('dtMethod', HTMLSpanElement),
    detailPath: requireElement('dtPath', HTMLSpanElement),
    detailTimestamp: requireElement('dtTs', HTMLDivElement),
    tabContent: requireElement('tabContent', HTMLDivElement),
    toast: requireElement('toast', HTMLDivElement),
    copyButton: requireElement('copyButton', HTMLButtonElement),
    curlButton: requireElement('curlButton', HTMLButtonElement),
    replayButton: requireElement('replayButton', HTMLButtonElement),
    newEndpointButton: requireElement('newEndpointButton', HTMLButtonElement),
    loadTokenButton: requireElement('loadTokenButton', HTMLButtonElement),
    createEndpointButton: requireElement('createEndpointButton', HTMLButtonElement),
    tabButtons: Array.from(document.querySelectorAll<HTMLButtonElement>('.tab[data-tab]')),
  };
}

export type Elements = ReturnType<typeof getElements>;
