FROM gcr.io/distroless/static-debian11:debug-nonroot AS production

COPY athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.36.0-uclibc@sha256:206803b272290a34c539a70802747413cb56a9fedb24b99879b1e4ee45c1e203 /bin/wget /usr/bin/wget

EXPOSE 8080
ENV ATHENAEUM_DB_ROOT=/home/nonroot/athenaeum

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
