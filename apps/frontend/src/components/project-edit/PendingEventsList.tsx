import React from 'react';
import { useParams } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import {
    useGetProjectsIdPendingEvents,
    usePostEventsIdApprove,
    usePostEventsIdReject,
    getGetProjectsIdPendingEventsQueryKey
} from '../../api/generated-orval/moris';

interface PendingEventsListProps {
    isAdmin?: boolean;
}

export const PendingEventsList: React.FC<PendingEventsListProps> = ({ isAdmin = false }) => {
    const { id } = useParams<{ id: string }>();
    const projectId = id || '';

    const { data, isLoading, error } = useGetProjectsIdPendingEvents(projectId, {
        query: {
            enabled: !!projectId
        }
    });

    const queryClient = useQueryClient();

    const { mutate: approve } = usePostEventsIdApprove();
    const { mutate: reject } = usePostEventsIdReject();

    const handleSuccess = () => {
        queryClient.invalidateQueries({
            queryKey: getGetProjectsIdPendingEventsQueryKey(projectId)
        });
    };

    if (isLoading) return <div className="p-4 text-center text-gray-500">Loading pending events...</div>;
    if (error) return <div className="p-4 text-center text-red-500">Error loading events</div>;

    const events = data?.events || [];

    if (events.length === 0) {
        return <div className="p-4 text-center text-gray-500 italic">No pending events</div>;
    }

    return (
        <div className="space-y-4">
            <h3 className="text-lg font-semibold text-gray-800">Pending Approvals</h3>
            <div className="grid gap-4">
                {events.map(event => (
                    <div key={event.id} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                        <div>
                            <div className="flex items-center gap-2">
                                <span className="font-medium text-gray-900">{event.type}</span>
                                <span className={`text-xs px-2 py-0.5 rounded-full ${event.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                        event.status === 'approved' ? 'bg-green-100 text-green-800' :
                                            'bg-red-100 text-red-800'
                                    }`}>
                                    {event.status}
                                </span>
                            </div>
                            <p className="text-sm text-gray-600 mt-1">{event.details}</p>
                            <p className="text-xs text-gray-400 mt-1">
                                {event.at ? new Date(event.at).toLocaleString() : 'Unknown date'}
                            </p>
                        </div>

                        {isAdmin && event.status === 'pending' && (
                            <div className="flex gap-2 shrink-0">
                                <button
                                    onClick={() => approve({ id: event.id! }, { onSuccess: handleSuccess })}
                                    className="px-3 py-1.5 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-1"
                                    title="Approve"
                                >
                                    Approve
                                </button>
                                <button
                                    onClick={() => reject({ id: event.id! }, { onSuccess: handleSuccess })}
                                    className="px-3 py-1.5 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-1"
                                    title="Reject"
                                >
                                    Reject
                                </button>
                            </div>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
};
