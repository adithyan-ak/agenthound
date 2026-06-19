// shared/model barrel — cross-cutting client state stores consumed by both the
// app shell and multiple features (so they cannot live in any one feature).
export { useMarksStore } from "./marks";
export { useUIStore } from "./ui-store";
