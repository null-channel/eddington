import {
  FrontendApi,
  LoginFlow,
  RecoveryFlow,
  RegistrationFlow,
  SettingsFlow,
  VerificationFlow,
} from "@ory/client";
import { AxiosInstance } from "axios";

declare global {
  interface Window {
    $axios: AxiosInstance;
    $ory: FrontendApi;
  }
}
type Flow =
  | LoginFlow
  | RegistrationFlow
  | RecoveryFlow
  | SettingsFlow
  | VerificationFlow;
export { Flow };
