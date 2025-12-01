import React from 'react';
import { PendingEventsList } from './PendingEventsList';

export const AdminApprovalList: React.FC = () => {
    return <PendingEventsList isAdmin={true} />;
};
