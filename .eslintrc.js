module.exports = {
  extends: [
    'airbnb-base',
    'plugin:react/recommended',
    'prettier',
    'prettier/react',
  ],
  env: {
    browser: true,
    es6: true,
    node: true,
    'jest/globals': true,
  },
  plugins: ['import', 'react', 'jest', 'prettier'],
  rules: {
    'prettier/prettier': ['error', { singleQuote: true, trailingComma: 'es5' }],
    'import/no-extraneous-dependencies': 0,
    'import/no-unresolved': 0,
    'import/extensions': 0,
    'func-names': 0,
    'no-underscore-dangle': 0, // because of _super
    'no-plusplus': ['error', { allowForLoopAfterthoughts: true }],
    'react/prop-types': 0,
    'class-methods-use-this': 0, // strange rule. It doesn't allow to create method render() without this
    'no-console': ['error', { allow: ['info', 'warn', 'error'] }], // later we will might need wrapper
  },
};
