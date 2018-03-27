import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import reducer, {
  initialState,
  LOAD,
  LOAD_SUCCESS,
  load,
  mostSimilar,
  leastSimilar,
} from './features';
import { ADD as ERROR_ADD } from './errors';

const mockStore = configureMockStore([thunk]);

describe('features/reducer', () => {
  it('LOAD', () => {
    expect(reducer(initialState, { type: LOAD })).toMatchSnapshot();
  });

  it('LOAD_SUCCESS', () => {
    expect(
      reducer(initialState, { type: LOAD_SUCCESS, features: [1, 2] })
    ).toMatchSnapshot();
  });
});

describe('features/actions', () => {
  describe('load', () => {
    it('success', () => {
      const filePairId = 1;
      const store = mockStore({
        features: {
          ...initialState,
        },
      });

      fetch.mockResponses([
        JSON.stringify({
          data: {
            featuresA: [
              { name: 'feature1', weight: 0.9 },
              { name: 'feature2', weight: 0.8 },
            ],
            featuresB: [
              { name: 'feature1', weight: 0.8 },
              { name: 'feature3', weight: 0.1 },
            ],
            score: { name: 'score', weight: 0.5 },
          },
        }),
      ]);

      return store.dispatch(load(filePairId)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
          },
          {
            type: LOAD_SUCCESS,
            features: {
              features: [
                { name: 'feature1', weightA: 0.9, weightB: 0.8 },
                { name: 'feature2', weightA: 0.8, weightB: 0 },
                { name: 'feature3', weightA: 0, weightB: 0.1 },
              ],
              score: 0.5,
            },
          },
        ]);
      });
    });

    it('error', () => {
      const blobIdA = 1;
      const blobIdB = 2;
      const store = mockStore({
        features: {
          ...initialState,
        },
      });

      const errText = 'some error';
      fetch.mockReject(errText);

      return store.dispatch(load(blobIdA, blobIdB)).then(() => {
        expect(store.getActions()).toEqual([
          {
            type: LOAD,
          },
          {
            type: ERROR_ADD,
            error: errText,
          },
        ]);
      });
    });
  });
});

describe('features/selectors', () => {
  const stateForSort = {
    features: {
      features: {
        features: [
          { name: 'a', weightA: 5, weightB: 10 },
          { name: 'b', weightA: 10, weightB: 10 },
          { name: 'c', weightA: 10, weightB: 1 },
          { name: 'b', weightA: 2, weightB: 2 },
          { name: 'd', weightA: 10, weightB: 10 },
        ],
        score: 0.5,
      },
    },
  };

  it('mostSimilar', () => {
    expect(mostSimilar(stateForSort)).toEqual([
      { name: 'b', weightA: 10, weightB: 10 },
      { name: 'd', weightA: 10, weightB: 10 },
      { name: 'b', weightA: 2, weightB: 2 },
      { name: 'a', weightA: 5, weightB: 10 },
      { name: 'c', weightA: 10, weightB: 1 },
    ]);
  });

  it('leastSimilar', () => {
    expect(leastSimilar(stateForSort)).toEqual([
      { name: 'c', weightA: 10, weightB: 1 },
      { name: 'a', weightA: 5, weightB: 10 },
      { name: 'b', weightA: 2, weightB: 2 },
      { name: 'd', weightA: 10, weightB: 10 },
      { name: 'b', weightA: 10, weightB: 10 },
    ]);
  });
});
