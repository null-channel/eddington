import React from 'react';
interface StatusCardProps {
    serviceName: string;
    bandwidth: number;
    cpuUsage: number;
    memoryUsage: number;
    status: 'Operational' | 'Non-Operational'; // Assuming status can only be 'Operational' or 'Non-Operational'
}

const StatusCard: React.FC<StatusCardProps> = ({ serviceName, bandwidth, cpuUsage, memoryUsage, status }) => {
    return (
      <div className="flex justify-between items-center gap-4 p-4 bg-zinc-950 rounded-lg shadow-sm w-full md:w-1/2">
        <div className="flex items-center gap-4">
          <IconCheckcircle className="h-6 w-6 text-emerald-400" />
          <div>
            <h2 className="text-lg font-semibold text-white">{serviceName}</h2>
            <p className="text-sm text-zinc-400">{status}</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-sm text-zinc-400">Bandwidth: {bandwidth}</p>
          <p className="text-sm text-zinc-400">CPU Usage: {cpuUsage}</p>
          <p className="text-sm text-zinc-400">Memory Usage: {memoryUsage}MB</p>
        </div>
      </div>
    )
  }
  
  function IconCheckcircle(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg
        {...props}
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
        <polyline points="22 4 12 14.01 9 11.01" />
      </svg>
    )
  }
  