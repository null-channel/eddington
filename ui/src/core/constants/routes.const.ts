const ROUTES = {
  MAIN: { path: '/main', name: 'dashbord_home' },
  NOT_FOUND: { path: '/:catchAll(.*)', name: 'NotFound' },
  HOME: { name: 'home', path: '/' },
  SIGN_IN: { name: 'sign-in', path: '/sign-in' },
  CONTACT_US: { name: 'contact-us', path: '/contactus' },
  RESET_PASSWORD: { name: 'reset-password', path: '/resetpassword' },
  UPDATE_PASSWORD: { name: 'update-password', path: '/updatepassword' },
};
export default ROUTES;
