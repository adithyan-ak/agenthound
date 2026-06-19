import { useEffect, useRef } from "react";

export interface UseEscapeKeyOptions {
  /** When false the listener is detached entirely. Defaults to true. */
  enabled?: boolean;
  /**
   * When true (default) Escape is ignored while focus is inside an
   * `<input>` or `<textarea>` so typing Esc in a field doesn't trigger
   * the handler.
   */
  ignoreInputs?: boolean;
}

/**
 * Subscribe to the Escape key once. Consolidates the hand-rolled
 * `window.addEventListener("keydown", ...)` Escape handlers that were
 * duplicated across the explorer surfaces. The latest `handler` is always
 * invoked via a ref, so passing an inline callback does not re-subscribe.
 */
export function useEscapeKey(
  handler: () => void,
  options: UseEscapeKeyOptions = {},
): void {
  const { enabled = true, ignoreInputs = true } = options;
  const handlerRef = useRef(handler);
  handlerRef.current = handler;

  useEffect(() => {
    if (!enabled) return;
    function onKey(e: KeyboardEvent) {
      if (e.key !== "Escape") return;
      if (ignoreInputs) {
        const target = e.target as HTMLElement | null;
        if (target?.tagName === "INPUT" || target?.tagName === "TEXTAREA") {
          return;
        }
      }
      handlerRef.current();
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [enabled, ignoreInputs]);
}
