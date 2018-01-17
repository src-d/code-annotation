import smallDiff from './respMock';

function mock(res) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(res);
    }, 500);
  });
}

function signIn() {}

function login() {
  return Promise.resolve({
    userId: 1,
    username: 'Maxim',
    avatarUrl: 'https://avatars3.githubusercontent.com/u/406916?s=40&v=4',
  });
}

function getExperiment() {
  return mock({
    id: 1,
    name: 'test experiment',
  });
}

function getAssignments() {
  return mock([
    {
      id: 1,
      pairId: 1,
      answer: 'Yes',
    },
    {
      id: 2,
      pairId: 2,
      answer: 'No',
    },
    {
      id: 3,
      pairId: 3,
      answer: null,
    },
    {
      id: 4,
      pairId: 4,
      answer: 'Yes',
    },
    {
      id: 5,
      pairId: 5,
      answer: null,
    },
  ]);
}

function getFilePair(experimentId, id) {
  return mock({
    id,
    diff: smallDiff,
  });
}

export default {
  signIn,
  login,

  getExperiment,
  getAssignments,
  getFilePair,
};
