import { type Plugin, type InjectionKey, type App } from "vue";
import { Configuration, FrontendApi } from "@ory/client";
import type { Session } from "@ory/client";
import { env } from "../constants";
export const Ory = new FrontendApi(
  new Configuration({
    basePath: env.ORY_URL.toString(),
    baseOptions: {
      withCredentials: true,
    },
  })
);

export const $ory: InjectionKey<typeof Ory> = Symbol("$ory");
export const $session: InjectionKey<Session> = Symbol("$session");
export const $ory_urls: InjectionKey<{
  logoutUrl: string;
}> = Symbol("$ory_urls");

export const OryPlugin: Plugin = {
  install: async (app: App) => {
    app.provide($ory, Ory);
  },
};
