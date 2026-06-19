/**
 * True when a keyboard event target is a text-editing element (`<input>` or
 * `<textarea>`), so global key handlers can ignore keystrokes the user is
 * typing into a field. Shared by `useEscapeKey` and the findings
 * keyboard-navigation handler so the guard is defined once.
 */
export function isEditableTarget(target: EventTarget | null): boolean {
  const el = target as HTMLElement | null;
  return el?.tagName === "INPUT" || el?.tagName === "TEXTAREA";
}
