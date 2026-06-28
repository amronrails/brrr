export interface Label {
  id: string;
  name: string;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface LabelInput {
  name: string;
  color: string;
}
