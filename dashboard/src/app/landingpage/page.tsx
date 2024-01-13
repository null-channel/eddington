
/**
 * v0 by Vercel.
 * @see https://v0.dev/t/7BDnnsqAzFA
 */
import { Button } from "@/components/ui/button"
import DayFiveButton  from "@/components/day5/button"
import { Input } from "@/components/ui/input"
import Link from "next/link"

export default function Component() {
  return (
    <div className="min-h-screen bg-[#0D1117] text-white">
      <nav className="flex justify-between items-center p-4">
        <div>
          <img
            alt="Logo"
            className="h-10"
            height="40"
            src="/placeholder.svg"
            style={{
              aspectRatio: "40/40",
              objectFit: "cover",
            }}
            width="40"
          />
        </div>
        <div className="space-x-4">
          <Button
            className="rounded-full bg-transparent border border-white text-white hover:bg-gray-500"
            variant="outline"
          >
            Login
          </Button>
          <Button className="rounded-full bg-white text-black hover:bg-gray-500" variant="default">
            Sign Up
          </Button>
        </div>
        <div>
          <DayFiveButton className="rounded-full bg-[#1DB954] text-white hover:bg-gray-500">
            Get Started
          </DayFiveButton>
        </div>
      </nav>
      <header className="flex flex-col items-center justify-center h-[80vh] space-y-6">
        <h1 className="text-4xl font-bold">NULL CLOUD</h1>
        <p className="text-center max-w-md">
          Are you a budding developer, ready to soar into the world of cloud computing but don't know where to start?
          Look no further than Null Cloud, your trusty sidekick in the realm of deployment and cloud management.
        </p>
        <form className="flex space-x-3">
          <Input className="rounded-full px-4 py-2 bg-white text-black" placeholder="your@email.com" />
          <Button className="rounded-full bg-[#1DB954] text-white hover:bg-gray-500" variant="default">
            Subscribe
          </Button>
        </form>
      </header>
      <footer className="bg-[#161B22] p-4">
        <div className="flex justify-center space-x-6">
          <Link className="block" href="#">
            <IconGithub className="h-6 w-6 text-white" />
          </Link>
          <Link className="block" href="#">
            <IconDiscord className="h-6 w-6 text-white" />
          </Link>
          <Link className="block" href="#">
            <IconYoutube className="h-6 w-6 text-white" />
          </Link>
          <Link className="block" href="#">
            <IconBook className="h-6 w-6 text-white" />
          </Link>
          <Link className="block" href="#">
            <IconHeadset className="h-6 w-6 text-white" />
          </Link>
        </div>
      </footer>
    </div>
  )
}

function IconBook(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20" />
    </svg>
  )
}


function IconDiscord(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <circle cx="12" cy="12" r="10" />
      <circle cx="12" cy="12" r="2" />
    </svg>
  )
}


function IconGithub(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
      <path d="M9 18c-4.51 2-5-2-7-2" />
    </svg>
  )
}


function IconHeadset(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M3 14h3a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-7a9 9 0 0 1 18 0v7a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3" />
    </svg>
  )
}


function IconYoutube(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M2.5 17a24.12 24.12 0 0 1 0-10 2 2 0 0 1 1.4-1.4 49.56 49.56 0 0 1 16.2 0A2 2 0 0 1 21.5 7a24.12 24.12 0 0 1 0 10 2 2 0 0 1-1.4 1.4 49.55 49.55 0 0 1-16.2 0A2 2 0 0 1 2.5 17" />
      <path d="m10 15 5-3-5-3z" />
    </svg>
  )
}
