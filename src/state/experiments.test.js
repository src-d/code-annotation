import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_ERROR,
  SET,
  load,
} from './experiments';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('experiments/reducer', () => {
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
        { ...initialState, loading: true },
        { type: LOAD_ERROR, error: 'test error' }
      )
    ).toMatchSnapshot();
  });

  it('SET', () => {
    expect(
      reducer(initialState, { type: SET, data: [{ some: 'data' }] })
    ).toMatchSnapshot();
  });
});

describe('experiments/actions', () => {
  describe('load', () => {
    it('success', () => {
      const store = mockStore({
        experiment: initialState,
      });

      const expList = [
        {
          id: 1,
          name: 'experiment name',
        },
        {
          id: 2,
          name: 'experiment name',
        },
      ];

      fetch.mockResponse(
        JSON.stringify({
          data: expList,
        })
      );

      store.dispatch(load(1)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
          },
          {
            type: SET,
            data: expList,
          },
          {
            type: LOAD_SUCCESS,
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({
        experiment: initialState,
      });

      const errText = 'some error';
      fetch.mockReject(errText);

      store.dispatch(load(1)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
          },
          {
            type: ERROR_ADD,
            error: errText,
          },
          {
            type: LOAD_ERROR,
            error: [errText],
          },
        ]);
      });
    });
  });
});
