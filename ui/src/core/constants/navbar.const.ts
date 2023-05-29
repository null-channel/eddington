import { MAIN_ROUTES, ROUTES } from ".";
import { HomeIcon, ArrowLeftOnRectangleIcon } from "@heroicons/vue/24/outline";

const NAVBAR = [
  { ...MAIN_ROUTES.HOME, icon: HomeIcon },
  { ...ROUTES.SIGN_IN, icon: ArrowLeftOnRectangleIcon },
];
export default NAVBAR;
