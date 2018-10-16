import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, { initialState, SET_OPTION, toggleInvisible } from './user';

const mockStore = configureMockStore([thunk]);

describe('user/reducer', () => {
  afterEach(() => {
    localStorage.clear();
  });

  it('SET_OPTION', () => {
    expect(
      reducer(
        { ...initialState, loading: true },
        {
          type: SET_OPTION,
          name: 'showInvisible',
          value: true,
        }
      )
    ).toMatchSnapshot();

    expect(
      window.localStorage.getItem('user_option_showInvisible')
    ).toBeTruthy();
  });
});

describe('user/actions', () => {
  describe('toggleInvisible', () => {
    it('enable', () => {
      const store = mockStore({
        user: initialState,
      });

      store.dispatch(toggleInvisible());
      expect(store.getActions()).toEqual([
        {
          type: SET_OPTION,
          name: 'showInvisible',
          value: true,
        },
      ]);
    });

    it('disable', () => {
      const store = mockStore({
        user: {
          ...initialState,
          showInvisible: true,
        },
      });

      store.dispatch(toggleInvisible());
      expect(store.getActions()).toEqual([
        {
          type: SET_OPTION,
          name: 'showInvisible',
          value: false,
        },
      ]);
    });
  });
});
