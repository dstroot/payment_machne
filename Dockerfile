##############################################
# Micro Docker image (3mb!!)
##############################################

# Docker base image which builds off of the "scratch" image
# and adds the root certificates for all of the standard
# certificate authorities.

# Certs add 258 kB to the 0 byte scratch image which serves
# as the root layer for all Docker images.

# The only other thing in the resulting image will be the
# copied binary so the total image size will be roughly
# the same as the binary itself.

FROM centurylink/ca-certs

# Bring the directory over (config and SQL)
COPY . /

# Run the program when the container starts
ENTRYPOINT ["/payment_machine"]

# Document that the service listens on port 8000
EXPOSE 8000
