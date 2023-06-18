import { type Plugin, type InjectionKey, type App, inject } from "vue";
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
    // can now be used with inject($ory)
    app.provide($ory, Ory);
    // injection is headache Don't do it 
    // try {
    //   // can now be used with inject($session)
    //   const { data: session } = await Ory.toSession();
    //   app.provide($session, session);
    // } catch (error) {
    //   console.log("[Ory] User has no session.");
      
    // }
    // try{
    //   const [{data: {logout_url}}]= await Promise.all([Ory.createBrowserLogoutFlow(undefined,{
    //     params: {
    //       return_url: "/",
    //     },
    //   })])
    //   app.provide($ory_urls, {
    //     logoutUrl: logout_url,
    //   });
    
    // }catch(error){
    //   console.log("[Ory] Failed to retrieve logout URL.");
    // }
  },
};
