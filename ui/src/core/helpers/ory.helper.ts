import type { Plugin, InjectionKey } from "vue";
import { Configuration, FrontendApi } from "@ory/client";
import type { Session } from "@ory/client";
import { env } from "../constants";
console.log(import.meta.env);
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
  install(app) {
    // can now be used with inject($ory)
    app.provide($ory, Ory);

    // can now be used with inject($session)
    Ory.toSession()
      .then(({ data }: any) => {
        app.provide($session, data);
      })
      .catch(() => {
        console.log("[Ory] User has no session.");
      });

    Promise.all([
      // get the logout url
      Ory.createBrowserLogoutFlow(undefined, {
        params: {
          return_url: "/",
        },
      }).catch(() => ({
        data: {
          logout_url: "",
        },
      })),
    ]).then(([{ data: logoutData }]) => {
      app.provide($ory_urls, {
        logoutUrl: logoutData.logout_url,
      });
    });
  },
};
