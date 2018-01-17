import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, { initialState, ADD, REMOVE, add, remove } from './errors';

const mockStore = configureMockStore([thunk]);

describe('errors/reducer', () => {
  it('ADD', () => {
    expect(
      reducer(initialState, { type: ADD, error: 'test error' })
    ).toMatchSnapshot();
  });

  it('REMOVE', () => {
    expect(
      reducer(['error 1', 'error 2', 'error 3'], {
        type: REMOVE,
        index: 1,
      })
    ).toMatchSnapshot();
  });
});

describe('errors/actions', () => {
  it('add', () => {
    const store = mockStore({
      errors: initialState,
    });

    store.dispatch(add('test error'));
    expect(store.getActions()).toEqual([
      {
        type: ADD,
        error: 'test error',
      },
    ]);
  });

  it('add array', () => {
    const store = mockStore({
      errors: initialState,
    });

    store.dispatch(add(['error1', 'error2']));
    expect(store.getActions()).toEqual([
      {
        type: ADD,
        error: 'error1',
      },
      {
        type: ADD,
        error: 'error2',
      },
    ]);
  });

  it('remove', () => {
    const store = mockStore({
      errors: initialState,
    });

    store.dispatch(remove(1));
    expect(store.getActions()).toEqual([
      {
        type: REMOVE,
        index: 1,
      },
    ]);
  });
});
