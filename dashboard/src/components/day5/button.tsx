import { Button } from "@/components/ui/button";

interface Props {
  className?: string;
  children?: React.ReactElement | string;
  icon?: React.ReactElement;
}

export default function DayFiveButton({ className, children, icon = <TreesIcon /> }: Props) {
  return (
    <div className={`flex flex-wrap gap-4 p-4 bg-[#a0d2db] ${className}`}>
      <Button className="flex justify-center bg-[#f4a261] text-white px-6 py-2 rounded-lg shadow-md">
        {icon}
        {children}
      </Button>
    </div>
  );
}

DayFiveButton.defaultProps = {
  className: "",
};

function TreesIcon(props: React.SVGProps<SVGSVGElement>) {
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
      <path d="M10 10v.2A3 3 0 0 1 8.9 16v0H5v0h0a3 3 0 0 1-1-5.8V10a3 3 0 0 1 6 0Z" />
      <path d="M7 16v6" />
      <path d="M13 19v3" />
      <path d="M12 19h8.3a1 1 0 0 0 .7-1.7L18 14h.3a1 1 0 0 0 .7-1.7L16 9h.2a1 1 0 0 0 .8-1.7L13 3l-1.4 1.5" />
    </svg>
  );
}
