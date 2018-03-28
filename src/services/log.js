const noop = () => undefined;

export default {
  // eslint-disable-next-line no-console
  error: process.env.NODE_ENV !== 'test' ? console.error.bind(console) : noop,
};
