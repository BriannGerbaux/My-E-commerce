import axios from 'axios';
import { create } from 'zustand'

interface Product {
    id: string;
    name: string;
    price_in_dollar: string;
    thumbnail_url: string;
    created_at: string;
    updated_at: string;
}

interface ProductState {
  loading: boolean;
  products: Product[];
}

export const useProductStore = create<ProductState>((set) => ({
    loading: false,
    products: [],

    getProducts: async () => {
        set({loading: true});
        try {
          const response = await axios.get<Product[]>(
            'http://localhost:8080/products',
          );
          set({ loading: false, products: response.data });
        } catch (error) {
          const errorMessage = 'Failed to login';
          if (axios.isAxiosError(error)) {
            console.error(error.message || errorMessage)
          }
          set({ loading: false });
        }
    }
}));