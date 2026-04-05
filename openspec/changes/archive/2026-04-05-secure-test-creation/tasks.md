## 1. CI Account Setup

- [x] 1.1 Create `github-actions` user on bench.likeawizard.dev (manual — document steps)
- [x] 1.2 Add `OPENBENCH_CI_USER` and `OPENBENCH_CI_PASS` as GitHub secrets
- [x] 1.3 Update `openbench-test.yml` to use CI credentials instead of admin credentials

## 2. Direct Test Link

- [x] 2.1 Parse test ID from Location header and export as step output
- [x] 2.2 Update PR comment to include direct link to `https://bench.likeawizard.dev/test/{id}/`

## 3. PGN Uploads

- [x] 3.1 Change `upload_pgns=FALSE` to `upload_pgns=TRUE` in the API call
