import { useEffect, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { usePostAuthOrcidLink, useGetProfile } from "@api/moris";
import { toast } from "sonner";
import { useAuth } from "@/hooks/useAuth";

const OrcidCallbackRoute = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { updateUser } = useAuth();
  const { mutateAsync: linkORCID } = usePostAuthOrcidLink();
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
      toast.error("ORCID Connection Failed", {
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
        await linkORCID({ data: { code } });
        try {
          const { data: user } = await fetchProfile();
          if (user) {
            updateUser({
              id: user.id,
              person_id: user.person_id,
              email: user.email,
              name: user.name,
              givenName: user.givenName,
              familyName: user.familyName,
              orcid: user.orcid,
              is_sys_admin: user.is_sys_admin,
            });
          }
        } catch (error) {
          console.error("Failed to update user context:", error);
        }

        toast.success("Success", {
          description: "Your ORCID iD has been successfully linked.",
        });
      } catch (err: any) {
        toast.error("Connection Failed", {
          description:
            err.response?.data?.message || "Failed to link ORCID iD.",
        });
      } finally {
        navigate("/dashboard/profile", { replace: true });
      }
    };

    link();
  }, [searchParams, navigate, toast, linkORCID, fetchProfile, updateUser]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-bold">Connecting to ORCID...</h2>
        <p className="text-muted-foreground">
          Please wait while we link your account.
        </p>
      </div>
    </div>
  );
};

export default OrcidCallbackRoute;
