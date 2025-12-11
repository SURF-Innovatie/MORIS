import { useEffect, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useToast } from "@/hooks/use-toast";
import { postAuthSurfconextLogin } from "@api/moris";
import { SURFconextLoginRequest } from "@api/model";

export default function SURFconextCallbackRoute() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { login } = useAuth();
  const { toast } = useToast();
  const processedRef = useRef(false);

  useEffect(() => {
    const code = searchParams.get("code");
    const error = searchParams.get("error");

    if (processedRef.current) return;
    
    // If we have neither code nor error, do nothing (shouldn't happen on callback)
    if (!code && !error) return;
    
    processedRef.current = true;

    if (error) {
      toast({
        title: "Login Failed",
        description: "SURFconext login was rejected or failed.",
        variant: "destructive",
      });
      navigate("/", { replace: true });
      return;
    }

    if (!code) {
      toast({
        title: "Login Failed",
        description: "No authorization code received from SURFconext.",
        variant: "destructive",
      });
      navigate("/", { replace: true });
      return;
    }

    const handleLogin = async () => {
      try {
        const request: SURFconextLoginRequest = { code };
        const response = await postAuthSurfconextLogin(request);
        
        if (response.token && response.user) {
            login(response.token, response.user);
            toast({
                title: "Login Successful",
                description: `Welcome back, ${response.user.person?.name || 'User'}!`,
            });
            navigate("/dashboard", { replace: true });
        } else {
            throw new Error("Invalid response from server");
        }

      } catch (err: any) {
        console.error("SURFconext login error:", err);
        toast({
          title: "Login Failed",
          description: "Failed to exchange code for token.",
          variant: "destructive",
        });
        navigate("/", { replace: true });
      }
    };

    handleLogin();
  }, [searchParams, navigate, login, toast]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-bold">Logging in with SURFconext...</h2>
        <p className="text-muted-foreground">
          Please wait while we verify your credentials.
        </p>
      </div>
    </div>
  );
}
