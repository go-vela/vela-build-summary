# SPDX-License-Identifier: Apache-2.0

################################################################################
##    docker build --no-cache --target certs -t vela-build-summary:certs .    ##
################################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

#################################################################
##    docker build --no-cache -t vela-build-summary:local .    ##
#################################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-build-summary /bin/vela-build-summary

ENTRYPOINT [ "/bin/vela-build-summary" ]
