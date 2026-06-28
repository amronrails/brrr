// Public API of the auth feature. Import auth from "@/features/auth"; the
// feature's internals (api, types) stay private to the folder.
export { AuthProvider, useAuth } from "./AuthContext";
export { ProtectedRoute } from "./ProtectedRoute";
export { LoginPage } from "./LoginPage";
export { RegisterPage } from "./RegisterPage";
