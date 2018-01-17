import mocks from './mocks';

const defaultServerUrl =
  process.env.REACT_APP_SERVER_URL || 'http://127.0.0.1:8080/api';

const apiUrl = url => `${defaultServerUrl}${url}`;

function checkStatus(resp) {
  if (resp.status < 200 || resp.status >= 300) {
    const error = new Error(resp.statusText);
    error.response = resp;
    throw error;
  }
  return resp;
}

function normalizeError(err) {
  if (typeof err === 'object') {
    // error from server
    if (err.title) {
      return err.title;
    }
    // javascript error
    if (err.message) {
      return err.message;
    }
    // weird object as error, shouldn't really happen
    return JSON.stringify(err);
  }
  if (typeof err === 'string') {
    return err;
  }
  return 'Internal error';
}

function normalizeErrors(err) {
  if (Array.isArray(err)) {
    return err.map(e => normalizeError(e));
  }
  return [normalizeError(err)];
}

function apiCall(url, options) {
  return fetch(apiUrl(url), options)
    .then(checkStatus)
    .then(resp => resp.json())
    .then(json => {
      if (json.errors) {
        throw json.errors;
      }
      return json.data;
    })
    .catch(err => Promise.reject(normalizeErrors(err)));
}

function signIn() {}

function login() {
  return Promise.resolve({
    userId: 1,
    username: 'Maxim',
    avatarUrl: 'https://avatars3.githubusercontent.com/u/406916?s=40&v=4',
  });
}

function getExperiment(experimentId) {
  return apiCall(`/experiments/${experimentId}`);
}

function getAssignments(experimentId) {
  return apiCall(`/experiments/${experimentId}/assignments`);
}

function getFilePair(experimentId, pairId) {
  return apiCall(`/experiments/${experimentId}/filePairs/${pairId}`);
}

// eslint-disable-next-line
let exportObject = {
  signIn,
  login,

  getExperiment,
  getAssignments,
  getFilePair,
};

if (process.env.NODE_ENV !== 'test') {
  exportObject = mocks;
}

export default exportObject;
