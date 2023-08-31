
- `Dockerfile`: build the image by directly copying the compiled binaries, fast build speed.
- `Dockerfile_build`: two-stage build of the image, slower build speed, you can specify the golang version.
- `Dockerfile_test`: container for testing rpc services.
