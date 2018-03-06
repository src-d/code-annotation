import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import { LOCATION_CHANGED, replace } from 'redux-little-router';
import { namedRoutes, makeUrl } from './routes';

import reducer, { middlewares } from './index';

// All those tests expect correct state because they run real actions

describe('routers', () => {
  const MiddlewaredMockStore = configureMockStore([...middlewares, thunk]);

  it('should skip unknown actions', () => {
    const store = MiddlewaredMockStore(reducer(undefined, {}));
    const action = { type: 'UNKNOWN' };
    store.dispatch(action);
    expect(store.getActions()).toEqual([action]);
  });

  it('not found', () => {
    const store = MiddlewaredMockStore(reducer(undefined, {}));
    const action = {
      type: LOCATION_CHANGED,
      payload: {},
    };
    store.dispatch(action);
    expect(store.getActions()).toEqual([
      { payload: {}, type: 'ROUTER_LOCATION_CHANGED' },
    ]);
  });

  describe('experiment', () => {
    it('should load experiment and redirect to question page', () => {
      const store = MiddlewaredMockStore(
        reducer(
          {
            assignments: {
              list: [{ id: 1, answer: null }],
            },
          },
          {}
        )
      );

      // mocked store doesn't run reducers
      // we need to mock only calls not data to keep actions happy
      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // assigments
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.experiment,
          params: { experiment: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions[actions.length - 1]).toEqual(
          replace(makeUrl('question', { experiment: 1, question: 1 }))
        );
      });
    });

    it('should load experiment and redirect to finish page', () => {
      const store = MiddlewaredMockStore(reducer(undefined, {}));

      // mocked store doesn't run reducers
      // we need to mock only calls not data to keep actions happy
      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // assigments
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.experiment,
          params: { experiment: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions[actions.length - 1]).toEqual(
          replace(makeUrl('finish', { experiment: 1 }))
        );
      });
    });
  });

  describe('assignments', () => {
    it('should select assigment on question page', () => {
      const store = MiddlewaredMockStore(
        reducer(
          {
            assignments: {
              list: [{ id: 1, answer: null }],
            },
          },
          {}
        )
      );

      // mocked store doesn't run reducers
      // we need to mock only calls not data to keep actions happy
      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // assigments
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // file pair
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.question,
          params: { experiment: 1, question: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions).toContainEqual({
          type: 'ca/assignments/SET_CURRENT_ASSIGMENT',
          assigment: { id: 1, answer: null },
        });
      });
    });

    it('should load experiments and assigments on finish page', () => {
      const store = MiddlewaredMockStore(reducer(undefined, {}));

      // mocked store doesn't run reducers
      // we need to mock only calls not data to keep actions happy
      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // assigments
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.finish,
          params: { experiment: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions).toContainEqual({ type: 'ca/experiment/LOAD' });
        expect(actions).toContainEqual({ type: 'ca/assignments/LOAD' });
      });
    });
  });

  describe('review', () => {
    it('should redirect to review assigment page', () => {
      const store = MiddlewaredMockStore(
        reducer(
          {
            filePairs: {
              list: [{ id: 1 }],
            },
          },
          {}
        )
      );

      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // file pairs
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.review,
          params: { experiment: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions[actions.length - 1]).toEqual(
          replace(makeUrl('reviewPair', { experiment: 1, pair: 1 }))
        );
      });
    });

    it('should load all necessary information for review page', () => {
      const store = MiddlewaredMockStore(
        reducer(
          {
            filePairs: {
              currentPairId: 1,
              pairs: {
                1: { leftBlobId: 'left-id', rightBlobId: 'right-id' },
              },
            },
          },
          {}
        )
      );

      fetch.mockResponses(
        // experiment
        [
          JSON.stringify({
            data: {},
          }),
          { status: 200 },
        ],
        // file pairs
        [
          JSON.stringify({
            data: [],
          }),
          { status: 200 },
        ],
        // annotations
        [
          JSON.stringify({
            data: [],
          }),
          { status: 200 },
        ],
        // features left
        [
          JSON.stringify({
            data: [],
          }),
          { status: 200 },
        ],
        // features right
        [
          JSON.stringify({
            data: [],
          }),
          { status: 200 },
        ]
      );

      const action = {
        type: LOCATION_CHANGED,
        payload: {
          route: namedRoutes.reviewPair,
          params: { experiment: 1, pair: 1 },
        },
      };
      return store.dispatch(action).then(() => {
        const actions = store.getActions();
        expect(actions).toContainEqual({ type: 'ca/experiment/LOAD' });
        expect(actions).toContainEqual({ type: 'ca/filePairs/LOAD' });
        expect(actions).toContainEqual({
          type: 'ca/filePairs/LOAD_ANNOTATIONS',
        });
      });
    });
  });
});
