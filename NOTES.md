You need GO verison 1.20 or higher

You need pack cli: brew install buildpacks/tap/pack 

pack config default-builder paketobuildpacks/builder-jammy-tiny 

export DOCKER_ORG=registry.harbor.learn.tapsme.org/convention-service

make image

