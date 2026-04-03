## 1. Docker Build

- [ ] 1.1 Create Dockerfile.worker with Python, Go, g++, make, git, and OpenBench client
- [ ] 1.2 Add worker service to docker-compose.yml under `worker` profile with configurable threads
- [ ] 1.3 Test worker image builds successfully

## 2. Tofiks Makefile

- [ ] 2.1 Add default Makefile target in tofiks repo that accepts EXE= variable for OpenBench compatibility

## 3. Verification

- [ ] 3.1 Start worker locally and verify it connects to bench.likeawizard.dev
- [ ] 3.2 Create a test in OpenBench UI and verify the worker picks it up and runs games
- [ ] 3.3 Verify test results appear in the OpenBench UI
