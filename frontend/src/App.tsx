import { Route, Routes } from "react-router-dom";
import { Layout } from "./components/Layout";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { DashboardPage } from "./pages/DashboardPage";
import { DocumentsPage } from "./pages/DocumentsPage";
import { ExpensesPage } from "./pages/ExpensesPage";
import { FamilyPage } from "./pages/FamilyPage";
import { LoginPage } from "./pages/LoginPage";
import { MarketingPage } from "./pages/MarketingPage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { ParentsPage } from "./pages/ParentsPage";
import { SettingsPage } from "./pages/SettingsPage";
import { SignupPage } from "./pages/SignupPage";
import { TasksPage } from "./pages/TasksPage";

function App() {
  return (
    <Routes>
      <Route path="/" element={<MarketingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignupPage />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<Layout />}>
          <Route path="/app" element={<DashboardPage />} />
          <Route path="/family" element={<FamilyPage />} />
          <Route path="/parents" element={<ParentsPage />} />
          <Route path="/tasks" element={<TasksPage />} />
          <Route path="/expenses" element={<ExpensesPage />} />
          <Route path="/documents" element={<DocumentsPage />} />
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
