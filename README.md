# Payment Machine

This program is designed to create NACHA format "6" and "7" records for the collections process.  These records are further processed to add the appropriate header/footer, etc. before being sent out.

Processing outline
------------------
* Select records from COLLECTIONS_PAYMENT to process
* Create a new file and process the records
    * Format the fields properly
    * Write the records to the file
    * Update the records as processed in SQL Server
    * Log any errors

Testing
-------
    $ go test
    $ go test -bench=.
    $ go test -cover

Compiling
---------

    $ go build && ./payment_machine

(Will produce an executable `payment_machine`)

Running
-------

    $ ./payment_machine

Installing
---------

    $ go install payment_machine.go

then run `payment_machine`.  

Docker
------

Our goal here is to generate literally the world's smallest container. One of the (many) benefits of developing with Go is that you have the option of compiling your application into a self-contained, statically-linked binary. A statically-linked binary can be run in a container with NO other dependencies which means you can create incredibly small images.

With a statically-linked binary, you could have a Dockerfile that looks something like this:

    FROM scratch
    COPY hello /
    ENTRYPOINT ["/hello"]

Note that the base image here is the **0 byte** scratch image which serves as the root layer for all Docker images. The only thing in the resulting image will be the copied binary so the total image size will be the same as the binary itself.

One of the downsides to using the scratch image is that you no longer have access to the root CA certificates which come pre-installed in most base images. There are a few different options for dealing with this:

* Disable SSL verification. This is not recommended for obvious reasons.
* Bundle the necessary root CA certificates as part of your application.
* Use a different base image which already contains the root CA certificates.

Option 3 seems simplest and safest. The centurylink/ca-certs image is simply the scratch image with the most common root CA certificates pre-installed. The resulting image is **only 258 kB** which is the smallest possible starting point for creating fully functional docker images.

Next, we will cross-complie our application to run on Linux and while we are at it we will stamp it with the build time and Git hash.  This is much better than a version number. This way all our golang applications are always explicit with the exact code that is running and when it was built. The build system itself automatically stamps it when it builds.

Golang has an elegant way of doing so. Using the linker option -X we can set a value for a symbol that can be accessed from within the binary. So, we use the below ldflags options for our applications:

    go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`" myapp.go

Next, we configure a `.dockerignore` file to ignore "everything" and then explicitly add in our executable and any additional configuration/ancillary files.

Next, we build our docker image with one of the simplest `Dockerfiles` **ever**.  

Finally, we push this up to Docker Hub.

All you have to do is run `./docker_build.sh` and boom!

To run the docker image locally:

    $ docker run --rm -it -p 8000:8000 dstroot/payment_machine

#### To learn more about small docker images:
* https://github.com/CenturyLinkLabs/golang-builder
* https://medium.com/@kelseyhightower/optimizing-docker-images-for-static-binaries-b5696e26eb07#.1rxjbsq1a
* http://blog.dimroc.com/2015/08/20/cross-compiled-go-with-docker/
