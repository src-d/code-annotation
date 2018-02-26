import { LOCATION_CHANGED } from 'redux-little-router';

export const load = trackingId => {
  /* eslint-disable */
  // prettier-ignore
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
    (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
    m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
    })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', trackingId, 'auto');
};

export const middleware = () => next => action => {
  if (action.type === LOCATION_CHANGED) {
    const { payload } = action;
    const path = payload.pathname + payload.search;
    ga('set', 'page', path);
    ga('send', 'pageview');
  }

  return next(action);
};
