## 1. CI Account Setup

- [ ] 1.1 Create `github-actions` user on bench.likeawizard.dev (manual — document steps)
- [ ] 1.2 Add `OPENBENCH_CI_USER` and `OPENBENCH_CI_PASS` as GitHub secrets
- [ ] 1.3 Update `openbench-test.yml` to use CI credentials instead of admin credentials

## 2. Direct Test Link

- [ ] 2.1 Parse test ID from Location header and export as step output
- [ ] 2.2 Update PR comment to include direct link to `https://bench.likeawizard.dev/test/{id}/`

## 3. PGN Uploads

- [ ] 3.1 Change `upload_pgns=FALSE` to `upload_pgns=TRUE` in the API call
