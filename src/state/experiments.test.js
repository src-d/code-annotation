import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  LOAD_ERROR,
  SET,
  CREATE,
  CREATE_SUCCESS,
  CREATE_ERROR,
  UPDATE,
  UPDATE_SUCCESS,
  UPDATE_ERROR,
  UPLOAD,
  UPLOAD_SUCCESS,
  UPLOAD_ERROR,
  UPLOAD_RESULT_RESET,
  UPLOAD_RES_NONE,
  UPLOAD_RES_SUCCESS,
  load,
  create,
  update,
  uploadFilePairs,
  uploadResultReset,
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

  it('CREATE', () => {
    expect(reducer(initialState, { type: CREATE })).toMatchSnapshot();
  });

  it('CREATE_SUCCESS', () => {
    expect(
      reducer(
        { ...initialState, createInProgress: true },
        { type: CREATE_SUCCESS }
      )
    ).toMatchSnapshot();
  });

  it('CREATE_ERROR', () => {
    expect(
      reducer(
        { ...initialState, createInProgress: true },
        { type: CREATE_ERROR, error: 'test error' }
      )
    ).toMatchSnapshot();
  });

  it('UPDATE', () => {
    expect(reducer(initialState, { type: UPDATE })).toMatchSnapshot();
  });

  it('UPDATE_SUCCESS', () => {
    expect(
      reducer(
        { ...initialState, updateInProgress: true },
        { type: UPDATE_SUCCESS }
      )
    ).toMatchSnapshot();
  });

  it('UPDATE_ERROR', () => {
    expect(
      reducer(
        { ...initialState, updateInProgress: true },
        { type: UPDATE_ERROR, error: 'test error' }
      )
    ).toMatchSnapshot();
  });

  it('UPLOAD', () => {
    expect(reducer(initialState, { type: UPLOAD })).toMatchSnapshot();
  });

  it('UPLOAD_SUCCESS', () => {
    expect(
      reducer(
        {
          ...initialState,
          uploadInProgress: true,
          uploadResult: UPLOAD_RES_NONE,
        },
        { type: UPLOAD_SUCCESS }
      )
    ).toMatchSnapshot();
  });

  it('UPLOAD_ERROR', () => {
    expect(
      reducer(
        {
          ...initialState,
          uploadInProgress: true,
          uploadResult: UPLOAD_RES_NONE,
        },
        { type: UPLOAD_ERROR, error: 'test error' }
      )
    ).toMatchSnapshot();
  });

  it('UPLOAD_RESULT_RESET', () => {
    expect(
      reducer(
        {
          ...initialState,
          uploadInProgress: false,
          uploadResult: UPLOAD_RES_SUCCESS,
        },
        { type: UPLOAD_RESULT_RESET }
      )
    ).toMatchSnapshot();
  });
});

const mockRes = {
  id: 3,
  name: 'new_name',
  description: 'new_desc',
  progress: 0,
};

describe('experiments/actions', () => {
  describe('load', () => {
    it('success', () => {
      const store = mockStore({
        experiments: initialState,
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

      store.dispatch(load()).then(() => {
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

      store.dispatch(load()).then(() => {
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

  describe('create', () => {
    it('success', () => {
      const store = mockStore({ experiments: initialState });

      fetch.mockResponseOnce(JSON.stringify({ data: mockRes }));

      store.dispatch(create('new_name', 'new_desc')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: CREATE,
          },
          {
            type: CREATE_SUCCESS,
          },
          {
            type: LOAD,
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({ experiments: initialState });

      const errText = 'some error';
      fetch.mockReject(errText);

      expect.hasAssertions();

      store.dispatch(create('new_name', 'new_desc')).catch(() => {
        expect(store.getActions()).toEqual([
          {
            type: CREATE,
          },
          {
            type: ERROR_ADD,
            error: errText,
          },
          {
            type: CREATE_ERROR,
            error: [errText],
          },
        ]);
      });
    });
  });

  describe('update', () => {
    it('success', () => {
      const store = mockStore({ experiments: initialState });

      fetch.mockResponseOnce(JSON.stringify({ data: mockRes }));

      store.dispatch(update(3, 'new_name', 'new_desc')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: UPDATE,
          },
          {
            type: UPDATE_SUCCESS,
          },
          {
            type: LOAD,
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({ experiments: initialState });

      const errText = 'some error';
      fetch.mockReject(errText);

      store.dispatch(update(3, 'new_name', 'new_desc')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: UPDATE,
          },
          {
            type: ERROR_ADD,
            error: errText,
          },
          {
            type: UPDATE_ERROR,
            error: [errText],
          },
        ]);
      });
    });
  });

  describe('upload', () => {
    it('success', () => {
      const store = mockStore({ experiments: initialState });

      fetch.mockResponseOnce(JSON.stringify({ data: mockRes }));

      store.dispatch(uploadFilePairs(3, 'file')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: UPLOAD,
          },
          {
            type: UPLOAD_SUCCESS,
          },
        ]);
      });
    });

    it('error', () => {
      const store = mockStore({ experiments: initialState });

      const errText = 'some error';
      fetch.mockReject(errText);

      store.dispatch(uploadFilePairs(3, 'file')).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: UPLOAD,
          },
          {
            type: ERROR_ADD,
            error: errText,
          },
          {
            type: UPLOAD_ERROR,
            error: [errText],
          },
        ]);
      });
    });

    it('reset', () => {
      const store = mockStore({ experiments: initialState });

      store.dispatch(uploadResultReset());

      expect(store.getActions()).toEqual([{ type: UPLOAD_RESULT_RESET }]);
    });
  });
});
