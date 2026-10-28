// Domain types mirroring the backend models (internal/models).
// Keep these in sync with the Go structs as fields are added.

export interface Family {
  id: string;
  name: string;
  invite_code: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface FamilyInvite {
  id: string;
  family_id: string;
  email: string;
  invited_by: string;
  status: "pending" | "accepted" | "revoked";
  created_at: string;
  accepted_at?: string;
}

export interface Member {
  id: string;
  family_id: string;
  user_id: string;
  email: string;
  full_name: string;
  role: "owner" | "member" | "viewer";
  joined_at: string;
}

export interface Parent {
  id: string;
  family_id: string;
  name: string;
  relation: string;
  doctor_name: string;
  doctor_phone: string;
  allergies: string;
  medications: string;
  insurance_provider: string;
  insurance_member_id: string;
  emergency_contact_name: string;
  emergency_contact_phone: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ParentKeyInfo {
  doctor_name: string;
  doctor_phone: string;
  allergies: string;
  medications: string;
  insurance_provider: string;
  insurance_member_id: string;
  emergency_contact_name: string;
  emergency_contact_phone: string;
}

export type TaskStatus = "todo" | "in_progress" | "done";
export type RecurrenceRule = "none" | "daily" | "weekly" | "monthly";

export interface Task {
  id: string;
  family_id: string;
  title: string;
  description: string;
  status: TaskStatus;
  assignee_id?: string;
  parent_id?: string;
  created_by: string;
  due_date?: string;
  recurrence_rule: RecurrenceRule;
  created_at: string;
  updated_at: string;
}

export interface Expense {
  id: string;
  family_id: string;
  paid_by: string;
  description: string;
  amount_cents: number;
  currency: string;
  category: string;
  parent_id?: string;
  document_id?: string;
  incurred_at: string;
  created_at: string;
  updated_at: string;
}

export interface ExpenseSplit {
  id: string;
  expense_id: string;
  member_id: string;
  share_cents: number;
  settled: boolean;
}

export interface Balance {
  member_id: string;
  net_cents: number;
}

export interface Document {
  id: string;
  family_id: string;
  uploaded_by: string;
  file_name: string;
  content_type: string;
  size_bytes: number;
  parent_id?: string;
  created_at: string;
}

export type ActivityVerb =
  | "task.created"
  | "task.completed"
  | "expense.created"
  | "document.uploaded"
  | "parent.created"
  | "member.joined";

export interface ActivityEntry {
  id: string;
  family_id: string;
  actor_id?: string;
  actor_name: string;
  verb: ActivityVerb;
  subject_type: string;
  subject_id?: string;
  subject_label: string;
  created_at: string;
}

export interface Contribution {
  member_id: string;
  member_name: string;
  tasks_completed: number;
  expenses_paid_count: number;
  expenses_paid_cents: number;
}

export interface NotificationPreferences {
  user_id: string;
  task_assigned_email: boolean;
  digest_cadence: string;
  updated_at: string;
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: AuthUser;
}
