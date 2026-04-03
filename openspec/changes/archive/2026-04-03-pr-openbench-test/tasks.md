## 1. Secrets

- [x] 1.1 Set OPENBENCH_USER and OPENBENCH_PASS secrets on tofiks-infra repo

## 2. Tofiks Repo (lightweight trigger)

- [x] 2.1 Create .github/workflows/openbench-test.yml in the tofiks repo: dispatch to tofiks-infra on PR open/synchronize/ready_for_review, skip drafts, pass PR branch and number

## 3. Tofiks-Infra Repo (test creation workflow)

- [x] 3.1 Create .github/workflows/openbench-test.yml in tofiks-infra: trigger on repository_dispatch type openbench-test
- [x] 3.2 Check out tofiks repo, build both branches (PR and master), extract bench values
- [x] 3.3 POST to bench.likeawizard.dev/scripts/ with test parameters
- [x] 3.4 Post a comment on the tofiks PR with the OpenBench test link

## 4. Verification

- [x] 4.1 Open a test PR on tofiks and verify the workflow creates an OpenBench test (pending merge + test PR)
- [x] 4.2 Verify PR comment appears with test link (pending merge + test PR)
