/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  'extends': [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript',
  ],
  parserOptions: {
    ecmaVersion: 'latest',
  },
  rules: {
    'semi': ['error'],
    'quotes': ['error', 'single'],
    'comma-dangle': [1, 'always-multiline'],
    'object-curly-spacing': [1, 'always'],
    'vue/no-mutating-props': 'off',
  },
};
