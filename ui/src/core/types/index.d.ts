import { FrontendApi } from "@ory/client";

declare global {
  interface Window {
    $axios: any;
    $ory: FrontendApi;
  }
}
export {};
