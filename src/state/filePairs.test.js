import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_FILE_PAIR,
  SET_FILE_PAIR,
  SET_CURRENT_PAIR,
  load,
  getCurrentFilePair,
} from './filePairs';

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

    store.dispatch(load(1)).then(() => {
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
