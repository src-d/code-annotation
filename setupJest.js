const { LocalStorage } = require('node-localstorage');

global.fetch = require('jest-fetch-mock');

global.localStorage = new LocalStorage('./localStorageTemp');

global.window = document.defaultView;
global.window.localStorage = global.localStorage;
