import { routerForBrowser } from 'redux-little-router';
import flattenRoutes from 'redux-little-router/lib/util/flatten-routes';
import UrlPattern from 'url-pattern';

// list of routes in the application
// every route must have unique property `name` to be able to easily create urls
// if a route should be available for unauthorized user it must contain `public: true`
const routes = {
  '/': {
    name: 'index',
    public: true,
    '/exp/:experiment': {
      name: 'experiment',
      '/finish': {
        name: 'finish',
      },
      '/:question': {
        name: 'question',
      },
    },
    '/review': {
      name: 'review',
    },
    '/export': {
      name: 'export',
    },
  },
};

const flatted = flattenRoutes(routes);
const namedRoutes = Object.keys(flatted).reduce((acc, path) => {
  const r = flatted[path];
  acc[r.name] = path;
  return acc;
}, {});

const namedPatterns = Object.keys(namedRoutes).reduce((acc, name) => {
  acc[name] = new UrlPattern(namedRoutes[name]);
  return acc;
}, {});

function makeUrl(name, args) {
  const pattern = namedPatterns[name];
  if (!pattern) {
    throw Error(`route with name ${name} doesn't exist`);
  }
  return pattern.stringify(args);
}

const { reducer, middleware, enhancer } = routerForBrowser({ routes });

export { middleware, enhancer, namedRoutes, makeUrl };
export default reducer;
