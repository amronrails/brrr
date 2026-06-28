export interface Product {
  id: string;
  name: string;
  sku: string;
  price: string;
  stock: number;
  active: boolean;
  metadata?: unknown;
  category_id: string;
  created_at: string;
  updated_at: string;
}

export interface ProductInput {
  name: string;
  sku: string;
  price: string;
  stock: number;
  active: boolean;
  metadata?: unknown;
  category_id: string;
}
