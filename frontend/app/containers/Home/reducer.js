/*
 *
 * Home reducer
 *
 */

import { fromJS } from 'immutable';
import {
  DEFAULT_ACTION,
  ARCHIVE_ACTION,
  ARCHIVE_RESP,
  ARCHIVE_ERROR,
} from './constants';

const initialState = fromJS({
  url: '',
  archivedUrl: '',
  loading: false,
  error: '',
});

function homeReducer(state = initialState, action) {
  switch (action.type) {
    case DEFAULT_ACTION:
      return state;

    case ARCHIVE_ACTION:
      return initialState.set('url', action.payload.url).set('error', '').set('loading', true);

    case ARCHIVE_RESP:
      return state.set('archivedUrl', action.payload.archivedUrl).set('loading', false);

    case ARCHIVE_ERROR:
      return state.set('error', action.payload.error).set('loading', false);

    default:
      return state;
  }
}

export default homeReducer;
