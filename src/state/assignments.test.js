import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import { push } from 'redux-little-router';
import reducer, {
  ANSWER_SIMILAR,
  ANSWER_MAYBE,
  ANSWER_DIFFERENT,
  ANSWER_SKIP,
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_ERROR,
  SET_ASSIGNMENTS,
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
} from './assignments';
import { LOAD_FILE_PAIR, SET_FILE_PAIR } from './filePairs';
import { makeUrl } from './routes';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('assignments/reducer', () => {
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
        {
          type: SET_FILE_PAIR,
          id: 1,
          diff: 'test diff',
          score: 0.92,
          leftBlobId: 1,
          rightBlobId: 2,
        }
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

  describe('MARK_ASSIGNMENT', () => {
    const currentAssigmentStartTime = new Date(1516188409733);

    it('existing assigment', () => {
      mockDate(1516188500000);
      expect(
        reducer(
          {
            ...initialState,
            list: [{ id: 1 }],
            currentAssigmentStartTime,
          },
          { type: MARK_ASSIGNMENT, id: 1, answer: 'yes' }
        )
      ).toMatchSnapshot();
    });

    it('unknown assigment', () => {
      mockDate(1516188500000);
      expect(
        reducer(
          {
            ...initialState,
            list: [{ id: 1 }],
            currentAssigmentStartTime,
          },
          { type: MARK_ASSIGNMENT, id: 2, answer: 'yes' }
        )
      ).toMatchSnapshot();
    });
  });
});

describe('assignments/actions', () => {
  const expId = 1;

  describe('selectAssigment', () => {
    const assigment = { id: 1, pairId: 1 };

    it('with file pair in cache', () => {
      const store = mockStore({
        assignments: {
          ...initialState,
          list: [assigment],
        },
        filePairs: {
          pairs: { 1: 'do not care about content here' },
        },
      });

      return store.dispatch(selectAssigment(expId, 1)).then(() => {
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
        assignments: {
          ...initialState,
          list: [{ id: 1, pairId: 1 }],
        },
        filePairs: {
          pairs: {},
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

      return store.dispatch(selectAssigment(expId, 1)).then(() => {
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
        assignments: {
          ...initialState,
          list: [{ id: 1, pairId: 1 }],
        },
        filePairs: {
          pairs: {},
        },
      });

      fetch.mockReject('some error');

      return store.dispatch(selectAssigment(expId, 1)).then(() => {
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
      const list = [
        { id: 1, answer: 'yes' },
        { id: 2, answer: null },
        { id: 3, answer: 'no' },
      ];

      const store = mockStore({
        assignments: {
          ...initialState,
          list,
        },
      });

      store.dispatch(nextAssigment(expId));
      expect(store.getActions()).toEqual([push('/exp/1/2')]);
    });

    it('no unanswered assignments', () => {
      const list = [{ id: 1, answer: 'yes' }];

      const store = mockStore({
        assignments: {
          ...initialState,
          list,
        },
      });

      store.dispatch(nextAssigment(expId));
      expect(store.getActions()).toEqual([push('/exp/1/finish')]);
    });
  });

  describe('load', () => {
    it('success', () => {
      const store = mockStore({
        assignments: initialState,
      });

      fetch.mockResponse(
        JSON.stringify({
          data: [
            {
              id: 1,
              answer: null,
            },
          ],
        })
      );
      return store.dispatch(load(expId)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
          },
          {
            type: SET_ASSIGNMENTS,
            assignments: [{ answer: null, id: 1 }],
          },
          {
            type: LOAD_SUCCESS,
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({
        assignments: initialState,
      });

      fetch.mockReject('some error');
      return store.dispatch(load(expId)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
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

  describe('mark', () => {
    it('success', () => {
      const store = mockStore({
        assignments: {
          ...initialState,
          list: [{ id: 1 }],
        },
      });

      fetch.mockResponse(JSON.stringify({ data: {} }));
      return store.dispatch(mark(expId, 1, 'yes')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: MARK_ASSIGNMENT,
            id: 1,
            answer: 'yes',
          },
          push(makeUrl('finish', { experiment: 1 })),
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({
        assignments: {
          ...initialState,
          list: [{ id: 1 }],
        },
      });

      fetch.mockReject('some error');
      return store.dispatch(mark(expId, 1, 'yes')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: MARK_ASSIGNMENT,
            id: 1,
            answer: 'yes',
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

  it('markCurrent', () => {
    const store = mockStore({
      experiment: {
        id: expId,
      },
      assignments: {
        ...initialState,
        list: [{ id: 1 }],
        currentAssigment: { id: 1 },
      },
    });

    fetch.mockResponse(JSON.stringify({ data: {} }));
    return store.dispatch(markCurrent('yes')).then(() => {
      expect(store.getActions()[0]).toEqual({
        type: MARK_ASSIGNMENT,
        id: 1,
        answer: 'yes',
      });
    });
  });
});

describe('assignments/selectors', () => {
  const state = {
    assignments: {
      list: [
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

  it('getProgressPercent empty', () => {
    expect(
      getProgressPercent({
        assignments: {
          list: [],
        },
      })
    ).toEqual(0);
  });

  it('getOverallTime', () => {
    expect(getOverallTime(state)).toEqual(9503);
  });

  it('getAverageTime', () => {
    expect(getAverageTime(state)).toEqual(1358);
  });

  it('getCurrentFilePair', () => {
    const pairState = {
      assignments: {
        currentAssigment: { id: 1, pairId: 1 },
      },
      filePairs: {
        pairs: { 1: 'something' },
      },
    };
    expect(getCurrentFilePair(pairState)).toEqual('something');
  });
});
