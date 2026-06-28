export interface Order {
  id: string;
  status: string;
  total: string;
  placed_at: string;
  customer_id: string;
  created_at: string;
  updated_at: string;
}

export interface OrderInput {
  status: string;
  total: string;
  placed_at: string;
  customer_id: string;
}
