FROM gcr.io/distroless/static-debian11:debug-nonroot AS production

COPY athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.37.0-uclibc@sha256:cc57e0ff4b6d3138931ff5c7180d18078813300e2508a25fb767a4d36df30d4d /bin/wget /usr/bin/wget

EXPOSE 8080
ENV ATHENAEUM_DB_ROOT=/home/nonroot/athenaeum

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
