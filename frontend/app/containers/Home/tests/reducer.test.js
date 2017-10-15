
import { fromJS } from 'immutable';
import homeReducer from '../reducer';

describe('homeReducer', () => {
  it('returns the initial state', () => {
    expect(homeReducer(undefined, {})).toEqual(fromJS({}));
  });
});
