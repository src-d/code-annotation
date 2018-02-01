import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import { push, replace } from 'redux-little-router';
import reducer, {
  ANSWER_SIMILAR,
  ANSWER_MAYBE,
  ANSWER_DIFFERENT,
  ANSWER_SKIP,
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_ERROR,
  SET_EXPERIMENT,
  SET_ASSIGNMENTS,
  LOAD_FILE_PAIR,
  SET_FILE_PAIR,
  SET_CURRENT_ASSIGMENT,
  MARK_ASSIGNMENT,
  selectAssigment,
  nextAssigment,
  load,
  mark,
  markCurrent,
  getIdenticalCount,
  getSimilarCount,
  getDifferentCount,
  getSkipCount,
  getProgressPercent,
  getOverallTime,
  getAverageTime,
  getCurrentFilePair,
} from './experiment';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('experiment/reducer', () => {
  const RealDate = Date;

  function mockDate(d) {
    global.Date = class extends RealDate {
      constructor() {
        return new RealDate(d);
      }
    };
  }

  afterEach(() => {
    global.Date = RealDate;
  });

  it('LOAD', () => {
    expect(reducer(initialState, { type: LOAD })).toMatchSnapshot();
  });

  it('LOAD_SUCCESS', () => {
    expect(
      reducer(
        { ...initialState, loading: true },
        {
          type: LOAD_SUCCESS,
          index: 1,
        }
      )
    ).toMatchSnapshot();
  });

  it('LOAD_ERROR', () => {
    expect(
      reducer(
        { ...initialState, loading: true, fileLoading: true },
        { type: LOAD_ERROR, error: 'test error' }
      )
    ).toMatchSnapshot();
  });

  it('SET_EXPERIMENT', () => {
    expect(
      reducer(initialState, { type: SET_EXPERIMENT, id: 1, name: 'test name' })
    ).toMatchSnapshot();
  });

  it('SET_ASSIGNMENTS', () => {
    expect(
      reducer(initialState, { type: SET_ASSIGNMENTS, assignments: [1, 2] })
    ).toMatchSnapshot();
  });

  it('LOAD_FILE_PAIR', () => {
    expect(reducer(initialState, { type: LOAD_FILE_PAIR })).toMatchSnapshot();
  });

  it('SET_FILE_PAIR', () => {
    expect(
      reducer(
        { ...initialState, loading: true },
        { type: SET_FILE_PAIR, id: 1, diff: 'test diff' }
      )
    ).toMatchSnapshot();
  });

  it('SET_CURRENT_ASSIGMENT', () => {
    mockDate(1516188409733);

    expect(
      reducer(initialState, {
        type: SET_CURRENT_ASSIGMENT,
        assigment: { id: 1, name: 'test' },
      })
    ).toMatchSnapshot();
  });

  it('MARK_ASSIGNMENT', () => {
    const currentAssigmentStartTime = new Date(1516188409733);
    mockDate(1516188500000);

    expect(
      reducer(
        {
          ...initialState,
          assignments: [{ id: 1 }],
          currentAssigmentStartTime,
        },
        { type: MARK_ASSIGNMENT, id: 1, answer: 'yes' }
      )
    ).toMatchSnapshot();
  });
});

describe('experiment/actions', () => {
  describe('selectAssigment', () => {
    const assigment = { id: 1, pairId: 1 };

    it('with file pair in cache', () => {
      const store = mockStore({
        experiment: {
          ...initialState,
          filePairs: { 1: 'do not care about content here' },
          assignments: [assigment],
        },
      });

      return store.dispatch(selectAssigment(1)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: SET_CURRENT_ASSIGMENT,
            assigment,
          },
        ]);
      });
    });

    it('success', () => {
      const store = mockStore({
        experiment: {
          ...initialState,
          assignments: [{ id: 1, pairId: 1 }],
        },
      });

      fetch.mockResponse(
        JSON.stringify({
          data: {
            id: 1,
            diff: 'some diff',
          },
        })
      );

      return store.dispatch(selectAssigment(1)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: SET_CURRENT_ASSIGMENT,
            assigment,
          },
          {
            type: LOAD_FILE_PAIR,
          },
          {
            type: SET_FILE_PAIR,
            id: 1,
            diff: 'some diff',
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({
        experiment: {
          ...initialState,
          assignments: [{ id: 1, pairId: 1 }],
        },
      });

      fetch.mockReject('some error');

      return store.dispatch(selectAssigment(1)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: SET_CURRENT_ASSIGMENT,
            assigment,
          },
          {
            type: LOAD_FILE_PAIR,
          },
          {
            type: ERROR_ADD,
            error: 'some error',
          },
          {
            type: LOAD_ERROR,
            error: ['some error'],
          },
        ]);
      });
    });
  });

  describe('nextAssigment', () => {
    it('has unanswered assignments', () => {
      const assignments = [
        { id: 1, answer: 'yes' },
        { id: 2, answer: null },
        { id: 3, answer: 'no' },
      ];

      const store = mockStore({
        experiment: {
          ...initialState,
          assignments,
        },
      });

      store.dispatch(nextAssigment());
      expect(store.getActions()).toEqual([push('/exp/1/2')]);
    });

    it('no unanswered assignments', () => {
      const assignments = [{ id: 1, answer: 'yes' }];

      const store = mockStore({
        experiment: {
          ...initialState,
          assignments,
        },
      });

      store.dispatch(nextAssigment());
      expect(store.getActions()).toEqual([push('/exp/1/finish')]);
    });
  });

  it('load', () => {
    const store = mockStore({
      experiment: initialState,
    });

    fetch.mockResponses(
      // experiment response
      [
        JSON.stringify({
          data: {
            id: 1,
            name: 'experiment name',
          },
        }),
      ],
      // assigments response
      [
        JSON.stringify({
          data: [
            {
              id: 1,
              answer: null,
            },
          ],
        }),
      ],
      // file pair response
      [
        JSON.stringify({
          data: {
            id: 1,
            diff: 'files diff',
          },
        }),
      ]
    );

    store.dispatch(load(1)).then(() => {
      expect(store.getActions()).toEqual([
        {
          type: LOAD,
        },
        {
          type: 'ca/experiment/SET_EXPERIMENT',
          id: 1,
          name: 'experiment name',
        },
        {
          type: 'ca/experiment/SET_ASSIGNMENTS',
          assignments: [{ answer: null, id: 1 }],
        },
        {
          type: 'ca/experiment/LOAD_SUCCESS',
        },
        replace('/exp/1/finish'),
      ]);
    });
  });

  // TODO: Test duration extension
  it('mark', () => {
    const store = mockStore({
      experiment: {
        ...initialState,
        assignments: [{ id: 1 }],
      },
    });

    store.dispatch(mark(1, 'yes'));
    expect(store.getActions()[0]).toEqual({
      type: MARK_ASSIGNMENT,
      id: 1,
      answer: 'yes',
    });
  });

  it('markCurrent', () => {
    const store = mockStore({
      experiment: {
        ...initialState,
        assignments: [{ id: 1 }],
        currentAssigment: { id: 1 },
      },
    });

    store.dispatch(markCurrent('yes'));
    expect(store.getActions()[0]).toEqual({
      type: MARK_ASSIGNMENT,
      id: 1,
      answer: 'yes',
    });
  });
});

describe('experiment/selectors', () => {
  const state = {
    experiment: {
      assignments: [
        { id: 1, answer: ANSWER_SIMILAR, duration: 2000 },
        { id: 2, answer: ANSWER_SIMILAR, duration: 1000 },
        { id: 3, answer: null },
        { id: 4, answer: ANSWER_DIFFERENT, duration: 5000 },
        { id: 5, answer: ANSWER_SKIP, duration: 1 },
        { id: 6, answer: ANSWER_SKIP, duration: 1 },
        { id: 7, answer: ANSWER_SKIP, duration: 1 },
        { id: 8, answer: ANSWER_MAYBE, duration: 1500 },
      ],
    },
  };

  it('getIdenticalCount', () => {
    expect(getIdenticalCount(state)).toEqual(2);
  });

  it('getSimilarCount', () => {
    expect(getSimilarCount(state)).toEqual(1);
  });

  it('getDifferentCount', () => {
    expect(getDifferentCount(state)).toEqual(1);
  });

  it('getSkipCount', () => {
    expect(getSkipCount(state)).toEqual(3);
  });

  it('getProgressPercent', () => {
    expect(getProgressPercent(state)).toEqual(50);
  });

  it('getOverallTime', () => {
    expect(getOverallTime(state)).toEqual(9503);
  });

  it('getAverageTime', () => {
    expect(getAverageTime(state)).toEqual(1358);
  });

  it('getCurrentFilePair', () => {
    const pairState = {
      experiment: {
        currentAssigment: { id: 1, pairId: 1 },
        filePairs: { 1: 'something' },
      },
    };
    expect(getCurrentFilePair(pairState)).toEqual('something');
  });
});
