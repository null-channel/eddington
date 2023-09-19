# eddington

The services behind the null cloud.

What are we doing? Building a cloud. How are we doing it? carefully. Want to join in? best way to do that is to join my discord and join my live streams. Want to just use it in production? wait a few months :)

checkout and sign up for updates on our website: nullcloud.io

## How to run locally?

So you want to give it a try? Here is how you can run a lot of it locally.

#### Requirements

First install:
 - kind
 - ctlptl
 - tilt

then to set up and configure your local kind k8s cluster:
```shell
./hack/start_dev.sh
```

To start running your dev environment simply run:
```shell
tilt up
```

Tilt will give you an address see the UI and you can access it via there. If you are looking to use the container building aspect you will need ssh into the container and then run:
```shell
docker login
```
you will need to login to whatever repo you will be using and configure the service to use it. (this is a wip)

At this point you should be able to access the API and your app should be running

