import axios from 'axios';
import { create } from 'zustand'

interface UserState {
  isLoggedIn: boolean;
  loading: boolean;
  token: string;
  username: string;
  email: string;
}

export const useUserStore = create<UserState>((set) => ({
  isLoggedIn: false,
  loading: false,
  token: "",
  username: "",
  email: "",

  login: async (email: string, password: string) => {
    set({loading: true});
    try {
      const response = await axios.post<string>(
        'http://localhost:8080/login',
        {
          email,
          password
        }
      );
      if (!response.data.startsWith("Bearer ")) {
        throw "No 'Bearer ' prefix"
      }
      set({ token: response.data.split(" ")[1], email: email, loading: false, isLoggedIn: true });
    } catch (error) {
      const errorMessage = 'Failed to login';
      if (axios.isAxiosError(error)) {
        console.error(error.message || errorMessage)
      }
      set({ token: "", username: "", email: "", loading: false, isLoggedIn: false });
    }
  },

  register: async (username: string, email: string, password: string) => {
    set({loading: true});
    try {
        await axios.post<string>(
        'http://localhost:8080/register',
        {
          username,
          email,
          password
        }
      );
      set({ loading: false, isLoggedIn: false });
    } catch (error) {
      const errorMessage = 'Failed to register';
      if (axios.isAxiosError(error)) {
        console.error(error.message || errorMessage)
      }
      set({ loading: false, isLoggedIn: false });
    }
  }
}))
