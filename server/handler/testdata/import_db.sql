CREATE TABLE files (
		blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT, uast_a BLOB,
		blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT, uast_b BLOB,
		score DOUBLE PRECISION);

INSERT INTO files values (
'3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e',
'github.com/bblfsh/dashboard.git',
'922e922e922e922e922e922e922e922e922e922e',
'project/src/a',
'Some text',
'uast1',

'3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e',
'github.com/bblfsh/dashboard.git',
'922e922e922e922e922e922e922e922e922e922e',
'other_project/src/b',
'Some text
',
'uast2',

0.9512810301340767
);

INSERT INTO files values (
'3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e',
'github.com/bblfsh/dashboard.git',
'922e922e922e922e922e922e922e922e922e922e',
'dashboard/src/services/api.js',
'import log from ''./log'';

const defaultServerUrl =
  process.env.REACT_APP_SERVER_URL || ''http://0.0.0.0:9999/api'';

const apiUrl = url => `${defaultServerUrl}${url}`;

const unexpectedErrorMsg =
  ''Unexpected error contacting babelfish server. Please, try again.'';

export function parse(language, filename, code, query, serverUrl) {
  return new Promise((resolve, reject) => {
    return fetch(apiUrl(''/parse''), {
      method: ''POST'',
      headers: {
        ''Content-Type'': ''application/json'',
      },
      body: JSON.stringify({
        server_url: serverUrl,
        language,
        filename,
        content: code,
        query,
      }),
    })
      .then(resp => resp.json())
      .then(({ status, errors, uast, language }) => {
        if (status === 0) {
          resolve({ uast, language });
        } else {
          reject(errors ? errors.map(normalizeError) : [''unexpected error'']);
        }
      })
      .catch(err => {
        log.error(err);
        reject([unexpectedErrorMsg]);
      });
  });
}

export function listDrivers() {
  return fetch(apiUrl(''/drivers''))
    .then(checkStatus)
    .then(resp => resp.json());
}

function checkStatus(resp) {
  if (resp.status < 200 || resp.status >= 300) {
    const error = new Error(resp.statusText);
    error.response = resp;
    throw error;
  }
  return resp;
}

function normalizeError(err) {
  if (typeof err === ''object'' && err.hasOwnProperty(''message'')) {
    return err.message;
  } else if (typeof err === ''string'') {
    return err;
  }

  return null;
}

export function getGist(gist) {
  return new Promise((resolve, reject) => {
    return fetch(apiUrl(''/gist?url='' + gist), {
      method: ''GET'',
      headers: {
        ''Content-Type'': ''application/json'',
      },
    })
      .then(checkStatus)
      .then(resp => resp.text())
      .then(code => resolve(code))
      .catch(err => {
        log.error(err);
        reject([err].map(normalizeError));
      });
  });
}

export function version(serverUrl) {
  return fetch(apiUrl(`/version`), {
    method: ''POST'',
    headers: {
      ''Content-Type'': ''application/json'',
    },
    body: JSON.stringify({
      server_url: serverUrl,
    }),
  })
    .then(checkStatus)
    .then(resp => resp.json())
    .catch(err => {
      log.error(err);
      return Promise.reject([err].map(normalizeError));
    });
}',
'',

'3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e3a6e',
'github.com/bblfsh/dashboard.git',
'922e922e922e922e922e922e922e922e922e922e',
'dashboard/src/services/api.js',
'const defaultServerUrl =
  process.env.REACT_APP_SERVER_URL || ''http://0.0.0.0:9999/api'';

const apiUrl = url => `${defaultServerUrl}${url}`;

const unexpectedErrorMsg =
  ''Unexpected error contacting babelfish server. Please, try again.'';

export function parse(language, code, serverUrl) {
  return new Promise((resolve, reject) => {
    return fetch(apiUrl(''/parse''), {
      method: ''POST'',
      headers: {
        ''Content-Type'': ''application/json''
      },
      body: JSON.stringify({
        server_url: serverUrl,
        language,
        content: code
      })
    })
      .then(resp => resp.json())
      .then(({ status, errors, uast }) => {
        if (status === ''ok'') {
          resolve(uast);
        } else {
          reject(errors.map(normalizeError));
        }
      })
      .catch(err => {
        console.error(err);
        reject([unexpectedErrorMsg]);
      });
  });
}

export function listDrivers() {
  return fetch(apiUrl(''/drivers'')).then(checkStatus).then(resp => resp.json());
}

function checkStatus(resp) {
  if (resp.status < 200 || resp.status >= 300) {
    const error = new Error(resp.statusText);
    error.response = resp;
    throw error;
  }
  return resp;
}

function normalizeError(err) {
  if (typeof err === ''object'' && err.hasOwnProperty(''message'')) {
    return err.message;
  } else if (typeof err === ''string'') {
    return err;
  }

  return null;
}',
'',

0.9512810301340767);
