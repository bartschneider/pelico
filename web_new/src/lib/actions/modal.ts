import { Modal } from 'bootstrap';

export function bootstrapModal(node: HTMLElement, { show }: { show: boolean }) {
  const modal = new Modal(node);

  if (show) {
    modal.show();
  }

  node.addEventListener('hidden.bs.modal', () => {
    node.dispatchEvent(new CustomEvent('close'));
  });

  return {
    update({ show }: { show: boolean }) {
      if (show) {
        modal.show();
      } else {
        modal.hide();
      }
    },
    destroy() {
      modal.dispose();
    }
  };
}
