## Build and Push Image

```bash
# build sponge image
./build-sponge-image.sh v1.8.3

# copy tag
docker tag zhufuyi/sponge:v1.8.3 zhufuyi/sponge:latest

# login docker
docker login -u zhufuyi -p

# push image
docker push zhufuyi/sponge:v1.8.3
docker push zhufuyi/sponge:latest
```