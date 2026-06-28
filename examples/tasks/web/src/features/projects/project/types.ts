export interface Project {
  id: string;
  name: string;
  key: string;
  description: string;
  archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProjectInput {
  name: string;
  key: string;
  description: string;
  archived: boolean;
}
