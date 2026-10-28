import { useCallback, useState } from "react";

interface ConfirmOptions {
  title: string;
  message: string;
  confirmLabel?: string;
  destructive?: boolean;
  onConfirm: () => void;
}

// Pairs with <ConfirmDialog>. Call confirm({...}) from an onClick instead
// of acting immediately — nothing destructive (delete a task, remove a
// member, ...) should fire straight off an icon click with no "are you
// sure", which the app did everywhere before this pass.
export function useConfirm() {
  const [state, setState] = useState<ConfirmOptions | null>(null);

  const confirm = useCallback((opts: ConfirmOptions) => setState(opts), []);
  const close = useCallback(() => setState(null), []);

  function handleConfirm() {
    state?.onConfirm();
    close();
  }

  return { state, confirm, close, handleConfirm };
}
