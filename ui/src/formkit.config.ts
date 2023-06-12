import { generateClasses } from "@formkit/themes";
import { genesisIcons } from "@formkit/icons";
import genesis from "@formkit/themes/tailwindcss/genesis";

export default {
  icons: {
    ...genesisIcons,
  },
  config: {
    classes: generateClasses({
      ...genesis,
      text: {
        outer: "mb-5",
        label: "input-label",
        input: "input",
        help: "text-xs text-gray-500",
        messages: "list-none p-0 mt-1 mb-0",
        message: "text-red-500 mb-1 text-xs",
      },
      submit: {
        label: "Log-in",
        outer: "mx-auto",
        input:
          "w-full !px-16 items-center flex flex-col button-primary !bg-gray-900 !text-gray-200",
      },
      email: {
        outer: "mb-5",
        label: "input-label",
        input: "input",
        help: "text-xs text-gray-500",
        messages: "list-none p-0 mt-1 mb-0",
        message: "text-red-500 mb-1 text-xs",
      },
      select: {
        outer: "mb-5",
        label: "input-label",
        input:
          'w-full pl-3 pr-8 py-2 border-none text-base text-gray-700 placeholder-gray-400 formkit-multiple:p-0 data-[placeholder="true"]:text-gray-400 formkit-multiple:data-[placeholder="true"]:text-inherit',
        help: "text-xs text-gray-500",
        messages: "list-none p-0 mt-1 mb-0",
        message: "text-red-500 mb-1 text-xs",
        option: "formkit-multiple:p-3 formkit-multiple:text-sm",
        selectIcon:
          "flex p-[3px] shrink-0 w-5 -ml-[1.5em] h-full pointer-events-none",
        inner:
          "input flex border border-gray-400 relative max-w-md items-center rounded mb-1 ring-1 ring-gray-400 focus-within:ring-blue-500 focus-within:ring-2 [&>span:first-child]:focus-within:text-blue-500",
      },
    }),
  },
};
