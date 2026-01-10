import { useEffect, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { usePostZenodoLink, useGetProfile } from "@api/moris";
import { useToast } from "@/hooks/use-toast";

const ZenodoCallbackRoute = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { mutateAsync: linkZenodo } = usePostZenodoLink();
  const { refetch: fetchProfile } = useGetProfile({
    query: {
      enabled: false,
    },
  });
  const processedRef = useRef(false);

  useEffect(() => {
    const code = searchParams.get("code");
    const error = searchParams.get("error");

    if (processedRef.current) return;
    processedRef.current = true;

    if (error) {
      toast({
        title: "Zenodo Connection Failed",
        description: "You denied access or an error occurred.",
        variant: "destructive",
      });
      navigate("/dashboard/profile", { replace: true });
      return;
    }

    if (!code) {
      toast({
        title: "Error",
        description: "No authorization code received.",
        variant: "destructive",
      });
      navigate("/dashboard/profile", { replace: true });
      return;
    }

    const link = async () => {
      try {
        await linkZenodo({ data: { code } });
        await fetchProfile();

        toast({
          title: "Success",
          description: "Your Zenodo account has been successfully linked.",
        });
      } catch (err: any) {
        toast({
          title: "Connection Failed",
          description:
            err.response?.data?.message || "Failed to link Zenodo account.",
          variant: "destructive",
        });
      } finally {
        navigate("/dashboard/profile", { replace: true });
      }
    };

    link();
  }, [searchParams, navigate, toast, linkZenodo, fetchProfile]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-bold">Connecting to Zenodo...</h2>
        <p className="text-muted-foreground">
          Please wait while we link your account.
        </p>
      </div>
    </div>
  );
};

export default ZenodoCallbackRoute;
