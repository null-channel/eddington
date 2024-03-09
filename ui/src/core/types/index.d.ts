import { AxiosInstance } from "axios";

declare global {
  interface Window {
    $axios: AxiosInstance;
  }
}
type Flow =
  | LoginFlow
  | RegistrationFlow
  | RecoveryFlow
  | SettingsFlow
  | VerificationFlow;
export { Flow };
