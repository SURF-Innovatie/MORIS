import { useEffect, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { usePostZenodoLink, useGetProfile } from "@api/moris";
import { toast } from "sonner";

const ZenodoCallbackRoute = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
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
      toast.error("Zenodo Connection Failed", {
        description: "You denied access or an error occurred.",
      });
      navigate("/dashboard/profile", { replace: true });
      return;
    }

    if (!code) {
      toast.error("Error", {
        description: "No authorization code received.",
      });
      navigate("/dashboard/profile", { replace: true });
      return;
    }

    const link = async () => {
      try {
        await linkZenodo({ data: { code } });
        await fetchProfile();

        toast.success("Success", {
          description: "Your Zenodo account has been successfully linked.",
        });
      } catch (err: any) {
        toast.error("Connection Failed", {
          description:
            err.response?.data?.message || "Failed to link Zenodo account.",
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
