export interface Task {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: number;
  due_date: string;
  done: boolean;
  project_id: string;
  assignee_id: string;
  created_at: string;
  updated_at: string;
}

export interface TaskInput {
  title: string;
  description: string;
  status: string;
  priority: number;
  due_date: string;
  done: boolean;
  project_id: string;
  assignee_id: string;
}
