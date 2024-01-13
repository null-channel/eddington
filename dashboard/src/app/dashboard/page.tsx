import React from 'react';
import SideMenu from '@/components/sidemenu';
import TopBar from '@/components/topbar';

const MainDashboard = () => {
    const badgeNumber = 12; // Replace with the actual badge number
    return (
        <div className="flex h-screen bg-gray-900">
            <div className="flex flex-col">
                <SideMenu badgeNumber={badgeNumber} />
                <TopBar />
                <div className="flex flex-col flex-grow">
                    <div className="h-full overflow-y-auto">
                        {/* Main content goes here */}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default MainDashboard;