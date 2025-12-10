import React, {
  createContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { UserAccount } from "@api/model";

interface AuthContextType {
  user: UserAccount | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: UserAccount) => void;
  logout: () => void;
  updateUser: (user: UserAccount) => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

import { STORAGE_KEY_AUTH_TOKEN, STORAGE_KEY_AUTH_USER } from "@/lib/constants";

// ... imports remain the same

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<UserAccount | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize auth state from localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem(STORAGE_KEY_AUTH_TOKEN);
    const storedUser = localStorage.getItem(STORAGE_KEY_AUTH_USER);

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser) as UserAccount;
        setToken(storedToken);
        setUser(parsedUser);
      } catch (error) {
        console.error("Failed to parse stored user data:", error);
        localStorage.removeItem(STORAGE_KEY_AUTH_TOKEN);
        localStorage.removeItem(STORAGE_KEY_AUTH_USER);
      }
    }
    setIsLoading(false);
  }, []);

  const login = useCallback((newToken: string, newUser: UserAccount) => {
    localStorage.setItem(STORAGE_KEY_AUTH_TOKEN, newToken);
    localStorage.setItem(STORAGE_KEY_AUTH_USER, JSON.stringify(newUser));
    setToken(newToken);
    setUser(newUser);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(STORAGE_KEY_AUTH_TOKEN);
    localStorage.removeItem(STORAGE_KEY_AUTH_USER);
    setToken(null);
    setUser(null);
  }, []);

  const updateUser = useCallback((updatedUser: UserAccount) => {
    localStorage.setItem(STORAGE_KEY_AUTH_USER, JSON.stringify(updatedUser));
    setUser(updatedUser);
  }, []);

  const value: AuthContextType = {
    user,
    token,
    isAuthenticated: !!token && !!user,
    isLoading,
    login,
    logout,
    updateUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
