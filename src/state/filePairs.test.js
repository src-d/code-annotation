import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_FILE_PAIR,
  SET_FILE_PAIR,
  SET_CURRENT_PAIR,
  LOAD_ANNOTATIONS,
  SET_ANNOTATIONS,
  load,
  getCurrentFilePair,
  loadFilePair,
  selectPair,
  loadAnnotations,
} from './filePairs';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('filePairs/reducer', () => {
  it('LOAD', () => {
    expect(reducer(initialState, { type: LOAD })).toMatchSnapshot();
  });

  it('LOAD_SUCCESS', () => {
    expect(
      reducer(
        { ...initialState, loading: true },
        {
          type: LOAD_SUCCESS,
          list: [1, 2, 3],
        }
      )
    ).toMatchSnapshot();
  });

  it('LOAD_FILE_PAIR', () => {
    expect(
      reducer(initialState, {
        type: LOAD_FILE_PAIR,
      })
    ).toMatchSnapshot();
  });

  it('SET_FILE_PAIR', () => {
    expect(
      reducer(initialState, {
        type: SET_FILE_PAIR,
        id: 1,
        diff: 'content',
        score: 0.987,
        leftBlobId: 'id1',
        rightBlobId: 'id2',
        leftLoc: 1,
        rightLoc: 2,
      })
    ).toMatchSnapshot();
  });

  it('SET_CURRENT_PAIR', () => {
    expect(
      reducer(initialState, {
        type: SET_CURRENT_PAIR,
        id: 1,
      })
    ).toMatchSnapshot();
  });

  it('LOAD_ANNOTATIONS', () => {
    expect(
      reducer(initialState, {
        type: LOAD_ANNOTATIONS,
      })
    ).toMatchSnapshot();
  });

  it('SET_ANNOTATIONS', () => {
    expect(
      reducer(initialState, {
        type: SET_ANNOTATIONS,
        id: 1,
      })
    ).toMatchSnapshot();
  });
});

describe('filePairs/actions', () => {
  it('load', () => {
    const store = mockStore({
      filePairs: initialState,
    });
    const pairs = [{ id: 1 }, { id: 2 }];

    fetch.mockResponse(
      JSON.stringify({
        data: pairs,
      })
    );

    return store.dispatch(load(1)).then(() => {
      expect(store.getActions()).toEqual([
        {
          type: LOAD,
        },
        {
          type: LOAD_SUCCESS,
          pairs,
        },
      ]);
    });
  });

  describe('loadFilePair', () => {
    it('success', () => {
      const store = mockStore({
        user: {
          showInvisible: false,
        },
      });
      const res = { id: 1 };

      const mockedFn = fetch.mockResponse(
        JSON.stringify({
          data: res,
        })
      );

      return store.dispatch(loadFilePair(1, 1)).then(() => {
        expect(store.getActions()).toEqual([
          { type: LOAD_FILE_PAIR },
          { type: SET_FILE_PAIR, ...res },
        ]);
        expect(mockedFn).toHaveBeenCalledWith(
          'http://127.0.0.1:8080/api/experiments/1/file-pairs/1',
          { credentials: 'include', headers: { Authorization: 'Bearer null' } }
        );
      });
    });

    it('success with spaces', () => {
      const store = mockStore({
        user: {
          showInvisible: true,
        },
      });
      const res = { id: 1 };

      const mockedFn = fetch.mockResponse(
        JSON.stringify({
          data: res,
        })
      );

      return store.dispatch(loadFilePair(1, 1)).then(() => {
        expect(store.getActions()).toEqual([
          { type: LOAD_FILE_PAIR },
          { type: SET_FILE_PAIR, ...res },
        ]);
        expect(mockedFn).toHaveBeenCalledWith(
          'http://127.0.0.1:8080/api/experiments/1/file-pairs/1?showInvisible=1',
          { credentials: 'include', headers: { Authorization: 'Bearer null' } }
        );
      });
    });

    it('error', () => {
      const store = mockStore({
        user: {
          showInvisible: false,
        },
      });
      const err = 'some error';

      fetch.mockReject(err);

      let failed = false;
      return store
        .dispatch(loadFilePair(1, 1))
        .catch(() => {
          failed = true;
        })
        .then(() => {
          expect(failed).toBeTruthy();
          expect(store.getActions()).toEqual([{ type: LOAD_FILE_PAIR }]);
        });
    });
  });

  describe('selectPair', () => {
    it('success', () => {
      const store = mockStore({
        filePairs: {
          ...initialState,
          pairs: { 1: 'pair' },
        },
        user: {
          showInvisible: false,
        },
      });
      return store.dispatch(selectPair(1, 1)).then(() => {
        expect(store.getActions()).toEqual([{ id: 1, type: SET_CURRENT_PAIR }]);
      });
    });

    it('error', () => {
      const store = mockStore({
        filePairs: initialState,
        user: {
          showInvisible: false,
        },
      });

      fetch.mockReject('some error');
      return store.dispatch(selectPair(1, 1)).then(() => {
        expect(store.getActions()).toEqual([
          { id: 1, type: SET_CURRENT_PAIR },
          { type: LOAD_FILE_PAIR },
          { error: 'some error', type: ERROR_ADD },
        ]);
      });
    });
  });

  describe('loadAnnotations', () => {
    it('success', () => {
      const store = mockStore({
        filePairs: {
          ...initialState,
          pairs: { 1: 'pair' },
        },
      });

      const annotations = { yes: 1, no: 2 };
      fetch.mockResponse(
        JSON.stringify({
          data: annotations,
        })
      );
      return store.dispatch(loadAnnotations(1, 1)).then(() => {
        expect(store.getActions()).toEqual([
          { type: LOAD_ANNOTATIONS },
          { type: SET_ANNOTATIONS, data: annotations },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({
        filePairs: {
          ...initialState,
          pairs: { 1: 'pair' },
        },
      });

      const errText = 'some error';
      fetch.mockReject(errText);

      return store.dispatch(loadAnnotations(1, 1)).then(() => {
        expect(store.getActions()).toEqual([
          { type: LOAD_ANNOTATIONS },
          {
            type: ERROR_ADD,
            error: errText,
          },
        ]);
      });
    });
  });
});

describe('filePairs/selector', () => {
  it('getCurrentFilePair', () => {
    const state = {
      filePairs: {
        currentPairId: 1,
        pairs: { 1: 'pair' },
      },
    };
    expect(getCurrentFilePair(state)).toEqual('pair');
  });
});
