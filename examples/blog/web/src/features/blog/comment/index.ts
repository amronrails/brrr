// Public API of the comment feature: pages for routing plus the query hooks,
// API client and types for reuse elsewhere in the app.
export { CommentListPage } from "./CommentListPage";
export { CommentFormPage } from "./CommentFormPage";
export * from "./hooks";
export * from "./api";
export type { Comment, CommentInput } from "./types";
