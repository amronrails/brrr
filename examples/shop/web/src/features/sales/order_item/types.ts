export interface OrderItem {
  id: string;
  quantity: number;
  unit_price: string;
  order_id: string;
  product_id: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItemInput {
  quantity: number;
  unit_price: string;
  order_id: string;
  product_id: string;
}
