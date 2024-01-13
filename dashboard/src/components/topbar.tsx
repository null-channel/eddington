/**
 * v0 by Vercel.
 * @see https://v0.dev/t/VfyohelFTCM
 */
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar"

export default function Component() {
  return (
    <div className="flex items-center justify-between bg-black text-white py-2 px-4">
      <div className="flex items-center">
        <IconCloud className="w-6 h-6" />
        <span className="ml-2">CloudApp</span>
      </div>
      <div className="flex items-center">
        <div className="mr-4">John Doe</div>
        <Avatar className="h-8 w-8">
          <AvatarImage alt="John Doe" src="/placeholder-avatar.jpg" />
          <AvatarFallback>JD</AvatarFallback>
        </Avatar>
      </div>
    </div>
  )
}

function IconCloud(props: any) {
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
      <path d="M17.5 19H9a7 7 0 1 1 6.71-9h1.79a4.5 4.5 0 1 1 0 9Z" />
    </svg>
  )
}

