/*
 *
 * Home actions
 *
 */

import {
  ARCHIVE_ACTION,
  ARCHIVE_ERROR,
  ARCHIVE_RESP,
  DEFAULT_ACTION,
} from './constants';

export function defaultAction() {
  return {
    type: DEFAULT_ACTION,
  };
}

export function archiveAction(url) {
  return {
    type: ARCHIVE_ACTION,
    payload: {
      url,
    },
  };
}

export function archiveResponse(resp) {
  return {
    type: ARCHIVE_RESP,
    payload: {
      id: resp.id,
      hash: resp.hash,
      archivedUrl: resp.archived_url,
    },
  };
}

export function archiveError(err) {
  return {
    type: ARCHIVE_ERROR,
    payload: {
      error: err,
    },
  };
}
