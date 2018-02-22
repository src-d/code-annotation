import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_ERROR,
  SET_EXPERIMENT,
  load,
} from './experiment';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('experiment/reducer', () => {
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

  it('SET_EXPERIMENT', () => {
    expect(
      reducer(initialState, { type: SET_EXPERIMENT, id: 1, name: 'test name' })
    ).toMatchSnapshot();
  });
});

describe('experiment/actions', () => {
  describe('load', () => {
    it('success', () => {
      const store = mockStore({
        experiment: initialState,
      });

      fetch.mockResponse(
        JSON.stringify({
          data: {
            id: 1,
            name: 'experiment name',
          },
        })
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
            type: 'ca/experiment/LOAD_SUCCESS',
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
