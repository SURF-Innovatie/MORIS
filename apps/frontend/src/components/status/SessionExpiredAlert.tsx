import { useEffect, useState } from 'react';
import { LogOut, AlertTriangle } from 'lucide-react';
import { Button } from '../ui/button';
import { useAuth } from '../../contexts/AuthContext';

export const SessionExpiredAlert = () => {
    const [isOpen, setIsOpen] = useState(false);
    const { logout } = useAuth();

    useEffect(() => {
        const handleSessionExpired = () => {
            setIsOpen(true);
        };

        window.addEventListener('auth:session-expired', handleSessionExpired);
        return () => window.removeEventListener('auth:session-expired', handleSessionExpired);
    }, []);

    const handleLoginRedirect = () => {
        setIsOpen(false);
        logout();
        // The logout function in AuthContext clears state and storage.
        // Assuming the router or AuthContext handles redirection to login upon logout/null user.
    };

    if (!isOpen) {
        return null;
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
            <div className="mx-4 w-full max-w-md rounded-2xl border border-yellow-500/20 bg-card p-6 shadow-2xl">
                <div className="flex items-start gap-4">
                    <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-yellow-500/10">
                        <AlertTriangle className="h-6 w-6 text-yellow-500" />
                    </div>
                    <div className="flex-1">
                        <h2 className="text-lg font-semibold text-foreground">Session Expired</h2>
                        <p className="mt-2 text-sm text-muted-foreground">
                            Your session has expired. Please log in again to continue using the application.
                        </p>
                        <Button
                            onClick={handleLoginRedirect}
                            className="mt-4 w-full"
                            variant="default"
                        >
                            <LogOut className="mr-2 h-4 w-4" />
                            Log In
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
};
