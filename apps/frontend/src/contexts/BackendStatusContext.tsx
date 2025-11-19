import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { AXIOS_INSTANCE } from '../api/custom-axios';

interface BackendStatusContextType {
  isOnline: boolean;
  isChecking: boolean;
  lastError: string | null;
  checkHealth: () => Promise<void>;
}

const BackendStatusContext = createContext<BackendStatusContextType | undefined>(undefined);

export const useBackendStatus = () => {
  const context = useContext(BackendStatusContext);
  if (!context) {
    throw new Error('useBackendStatus must be used within a BackendStatusProvider');
  }
  return context;
};

interface BackendStatusProviderProps {
  children: ReactNode;
  checkInterval?: number;
}

export const BackendStatusProvider: React.FC<BackendStatusProviderProps> = ({
  children,
  checkInterval = 30000,
}) => {
  const [isOnline, setIsOnline] = useState(true);
  const [isChecking, setIsChecking] = useState(false);
  const [lastError, setLastError] = useState<string | null>(null);

  const checkHealth = async () => {
    setIsChecking(true);
    try {
      const response = await AXIOS_INSTANCE.get('/health', {
        timeout: 5000,
      });
      
      if (response.status === 200) {
        setIsOnline(true);
        setLastError(null);
      } else {
        setIsOnline(false);
        setLastError('Backend returned unexpected status');
      }
    } catch (error: any) {
      setIsOnline(false);
      if (error.code === 'ECONNABORTED') {
        setLastError('Connection timeout - backend might be slow or down');
      } else if (error.code === 'ERR_NETWORK') {
        setLastError('Network error - cannot reach backend');
      } else {
        setLastError(error.message || 'Unable to connect to backend');
      }
    } finally {
      setIsChecking(false);
    }
  };

  useEffect(() => {
    checkHealth();
    
    const interval = setInterval(() => {
      checkHealth();
    }, checkInterval);

    return () => clearInterval(interval);
  }, [checkInterval]);

  return (
    <BackendStatusContext.Provider value={{ isOnline, isChecking, lastError, checkHealth }}>
      {children}
    </BackendStatusContext.Provider>
  );
};
