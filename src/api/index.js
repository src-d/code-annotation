import TokenService from '../services/token';

export const serverUrl =
  process.env.REACT_APP_SERVER_URL || 'http://127.0.0.1:8080';

const apiUrl = url => `${serverUrl}${url}`;

function checkStatus(resp) {
  // when server return Unauthorized we need to remove token
  if (resp.status === 401) {
    TokenService.remove();
  }
  if (resp.status < 200 || resp.status >= 300) {
    return resp
      .json()
      .catch(() => {
        throw new Error(resp.statusText);
      })
      .then(json => {
        if (json.errors) {
          throw json.errors;
        }
        throw new Error(resp.statusText);
      });
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

function apiCall(url, options = {}) {
  const token = TokenService.get();
  const fetchOptions = {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    },
  };
  if (options.body) {
    fetchOptions.body = JSON.stringify(options.body);
    fetchOptions.headers['Content-Type'] = 'application/json';
  }
  return fetch(apiUrl(url), fetchOptions)
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

function me() {
  return apiCall(`/api/me`);
}

function getExperiment(experimentId) {
  return apiCall(`/api/experiments/${experimentId}`);
}

function getAssignments(experimentId) {
  return apiCall(`/api/experiments/${experimentId}/assignments`);
}

function getFilePair(experimentId, pairId) {
  return apiCall(`/api/experiments/${experimentId}/file-pairs/${pairId}`);
}

function putAnswer(experimentId, assignmentId, answer) {
  return apiCall(
    `/api/experiments/${experimentId}/assignments/${assignmentId}`,
    {
      method: 'PUT',
      body: answer,
    }
  );
}

function exportList() {
  return apiCall('/api/exports');
}

function exportCreate() {
  return apiCall(`/api/exports`, {
    method: 'POST',
  });
}

function exportDownload(filename) {
  const token = TokenService.get();
  const url = apiUrl(`/api/exports/${filename}/download?jwt_token=${token}`);
  window.open(url, '_blank');
}

export default {
  me,
  getExperiment,
  getAssignments,
  getFilePair,
  putAnswer,
  exportList,
  exportCreate,
  exportDownload,
};
