ARG SERVICE_NAME=web-auth-server
FROM blocktop/golang:build@sha256:52ad9d664c9bb3f0af24142b1ed6cf35f72de4beb0f3e88763958329ddb3f054 as build

EXPOSE 3000

FROM build as test
CMD ["make", "test"]

FROM blocktop/golang:prod@sha256:866785f278ed6740ab0e0e34d0d9a75a3b41181001a1bef73a95388dd44f4f92

# hadolint ignore=DL3022
ONBUILD COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# hadolint ignore=DL3022
ONBUILD COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

ARG TLS_CERT
ARG TLS_KEY

EXPOSE 3000
