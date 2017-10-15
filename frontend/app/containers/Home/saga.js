import { ARCHIVE_ACTION } from './constants';

import { delay } from 'redux-saga';
import { takeLatest, all, cps, call, put, select, race } from 'redux-saga/effects';

import IpfsArchive from 'ipfs_archive';

import { makeSelectUrl } from './selectors';
import { archiveResponse, archiveError } from './actions';

export function* archive() {
  const url = yield select(makeSelectUrl());

  const api = new IpfsArchive.IpfsApi();

  const body = new IpfsArchive.ArchiveRequest();
  body.url = url;

  try {
    const { archiveResp, timeout } = yield race({
      archiveResp: cps([api, api.archiveUrl], body),
      timeout: call(delay, 30 * 1000),
    });

    if (archiveResp) { yield put(archiveResponse({ archived_url: archiveResp.archived_url })); } else { yield put(archiveError("Server didn't respond in time.")); }
  } catch (error) {
    yield put(archiveError(error));
  }
}

export function* archiveUrlAsync() {
  yield takeLatest(ARCHIVE_ACTION, archive);
}

export default function* defaultSaga() {
  yield archiveUrlAsync();
}
