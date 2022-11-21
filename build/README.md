
- `Dockerfile`: Build the image by directly copying the compiled binaries.
  - Pros: fast build speed
  - Disadvantage: image size is twice as big as two-stage build.
- `Dockerfile_build`: two-stage build of the image.
  - Pros: minimal image size
  - Disadvantages: slower build speed, each build produces a larger intermediate image, which needs to be cleaned up regularly by executing the command `docker rmi $(docker images | grep "<none>" | awk '{print $3}')`.
