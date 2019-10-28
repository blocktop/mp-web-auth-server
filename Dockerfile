ARG SERVICE_NAME=web-auth-service
FROM blocktop/golang:build@#4bc090cd0598e663633bba35aa8dedea28a14fcb2907a12883bee75734ba2748 as build

EXPOSE 3000

FROM build as test
CMD ["make", "test"]

FROM blocktop/golang:prod@#aff706ed49ea8ad4d7e04646fa25bfd51d2358796c82d8ca134e168009c29fdc

# hadolint ignore=DL3022
ONBUILD COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# hadolint ignore=DL3022
ONBUILD COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=build /go/src/github.com/neighborly/demand-api/demand-api-server /bin/demand-api-server
COPY --from=build /go/src/github.com/neighborly/demand-api/demand-api /bin/demand-api
COPY --from=build /go/src/github.com/neighborly/demand-api/config/config.json /go/src/github.com/neighborly/demand-api/config/config.json
COPY --from=build /go/src/github.com/neighborly/demand-api/db/migrations /db/migrations

EXPOSE 6002
EXPOSE 9060

ENTRYPOINT []
CMD ["demand-api-server"]
