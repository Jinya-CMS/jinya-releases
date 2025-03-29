/**
 * Displays a alert modal dialog
 * @param title {string}
 * @param message {string}
 * @param closeLabel {string}
 * @param negative {boolean}
 * @return {Promise<boolean>}
 */
export default async function alert({ title = window.location.href, message, closeLabel, negative = false }) {
  return new Promise((resolve) => {
    const container = document.createElement('div');
    document.body.appendChild(container);
    const modalId = crypto.randomUUID();

    container.innerHTML = `
      <div class="cosmo-modal__container">
        <div class="cosmo-modal ${negative ? 'is--negative' : ''}">
          <h1 class="cosmo-modal__title">${title}</h1>
          <p class="cosmo-modal__content">${message}</p>
          <div class="cosmo-modal__button-bar">
            <button id="${modalId}CloseButton" class="cosmo-button">${closeLabel}</button>
          </div>
        </div>
      </div>`;

    document.body.appendChild(container);

    document.getElementById(`${modalId}CloseButton`).addEventListener('click', (e) => {
      e.preventDefault();
      container.remove();
      resolve();
    });
  });
}
