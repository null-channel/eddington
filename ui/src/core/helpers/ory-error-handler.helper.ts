import type { AxiosError } from "axios";
import type { Router } from "vue-router";

const oryErrorHandler = (router: Router) => {
  const refreshFlow = () => {
    router.push({
      // for our use case, removing the flow
      // parameter from the search query and
      // reloading the page are sufficient
      // to refresh the flow
      query: {},
    });
    window.location.reload();
  };

  return (error: AxiosError) => {
    const responseData = error.response?.data as any;
    switch ((error.response?.data as any).error.id) {
      case "session_already_available": // User is already signed in, let's redirect them home!
        router.push("/");
        return Promise.resolve();
      case "session_aal2_required": // 2FA is enabled and enforced, but user did not perform 2fa yet!
        router.replace("/login?aal2=true");
        return Promise.resolve();
      case "session_refresh_required": // We need to re-authenticate to perform this action
        console.warn("sdkError 403: Redirect browser to");
        window.location = responseData.redirect_browser_to;
        return Promise.resolve();
      case "browser_location_change_required": // Ory Kratos asked us to point the user to this URL.
        window.location.href = responseData.redirect_browser_to;
        return Promise.resolve();
      case "self_service_flow_expired": // The flow expired, let's request a new one.
        break;
      case "self_service_flow_return_to_forbidden": // the return is invalid, we need a new flow
        break;
      case "security_csrf_violation": // A CSRF violation occurred. Best to just refresh the flow!
        break;
      case "security_identity_mismatch": // The requested item was intended for someone else. Let's request a new flow...
        //	refreshFlow(); This causes an infinite loop on chrome because it uses the same flow code each time
        router.push("/");
        return;
    }

    switch (error.response?.status) {
      case 410: // The flow expired, let's request a new one.
        refreshFlow();
        return;
      case 404: // User might be replaying old links, so go back to the start.
        router.push("/");
        return;
    }
    // We are not able to handle the error? Return it.
    return Promise.reject(error);
  };
};
export default oryErrorHandler;
