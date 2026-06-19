import "@testing-library/jest-dom/vitest";

// Node >=22 ships an experimental global `localStorage` (the test runner
// launches workers with `--localstorage-file`, which activates it). That
// global shadows jsdom's window.localStorage and throws
// "storage.setItem is not a function" the moment zustand's persist
// middleware writes — crashing every test that imports the persisted graph
// store. Install a deterministic in-memory Storage before any store is
// created so the persisted store works under the runner regardless of the
// Node version.
function createMemoryStorage(): Storage {
  const data = new Map<string, string>();
  return {
    get length() {
      return data.size;
    },
    clear() {
      data.clear();
    },
    getItem(key: string) {
      return data.get(key) ?? null;
    },
    key(index: number) {
      return [...data.keys()][index] ?? null;
    },
    removeItem(key: string) {
      data.delete(key);
    },
    setItem(key: string, value: string) {
      data.set(key, String(value));
    },
  } as unknown as Storage;
}

const memoryStorage = createMemoryStorage();
try {
  Object.defineProperty(globalThis, "localStorage", {
    value: memoryStorage,
    configurable: true,
    writable: true,
  });
} catch {
  (globalThis as { localStorage: Storage }).localStorage = memoryStorage;
}
