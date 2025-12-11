import { useParams } from "react-router-dom";
import { useGetUsersId } from "@api/moris";
import { ProfileInfo } from "@/components/profile/ProfileInfo";

const AdminUserEditRoute = () => {
    const { id } = useParams<{ id: string }>();
    const { data: user, isLoading, refetch } = useGetUsersId(id!, { query: { enabled: !!id } });

    if (isLoading) {
        return <div>Loading user...</div>;
    }

    if (!user) {
        return <div>User not found</div>;
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center gap-4">
                <h1 className="text-2xl font-bold tracking-tight">Edit User: {user.name}</h1>
            </div>

            <div className="max-w-2xl">
                <ProfileInfo user={user} refetchProfile={refetch} />
            </div>
        </div>
    );
};

export default AdminUserEditRoute;
