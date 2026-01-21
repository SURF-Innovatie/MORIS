import { useEffect, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useToast } from "@/hooks/use-toast";
import { useAuth } from "@/hooks/useAuth";
import { AXIOS_INSTANCE } from "@/api/custom-axios";

interface SurfconextLoginResponse {
  token: string;
  user: {
    id: string;
    person_id: string;
    email: string;
    name: string;
    givenName?: string;
    familyName?: string;
    orcid?: string;
    is_sys_admin: boolean;
  };
}

const SurfconextCallbackRoute = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { login } = useAuth();
  const processedRef = useRef(false);

  useEffect(() => {
    const code = searchParams.get("code");
    const error = searchParams.get("error");

    if (processedRef.current) return;
    processedRef.current = true;

    if (error) {
      toast({
        title: "SURFconext Login Failed",
        description: "You denied access or an error occurred.",
        variant: "destructive",
      });
      navigate("/", { replace: true });
      return;
    }

    if (!code) {
      toast({
        title: "Error",
        description: "No authorization code received.",
        variant: "destructive",
      });
      navigate("/", { replace: true });
      return;
    }

    const exchangeCode = async () => {
      try {
        const response = await AXIOS_INSTANCE.post<SurfconextLoginResponse>(
          "/auth/surfconext/login",
          { code },
        );
        const { token, user } = response.data;

        if (token && user) {
          login(token, user);
          toast({
            title: "Welcome",
            description: `Logged in as ${user.name || user.email}`,
          });
          const returnUrl = sessionStorage.getItem("auth_return_url");
          if (returnUrl) {
            sessionStorage.removeItem("auth_return_url");
            navigate(returnUrl, { replace: true });
          } else {
            navigate("/dashboard", { replace: true });
          }
        } else {
          throw new Error("Invalid response from server");
        }
      } catch (err: any) {
        toast({
          title: "Login Failed",
          description:
            err.message ||
            "Failed to complete SURFconext login. Your account may not exist in MORIS.",
          variant: "destructive",
        });
        navigate("/", { replace: true });
      }
    };

    exchangeCode();
  }, [searchParams, navigate, toast, login]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-bold">Logging in with SURFconext...</h2>
        <p className="text-muted-foreground">
          Please wait while we complete your login.
        </p>
      </div>
    </div>
  );
};

export default SurfconextCallbackRoute;
