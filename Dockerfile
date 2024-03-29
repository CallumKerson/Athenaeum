FROM gcr.io/distroless/static-debian11:debug-nonroot AS production

COPY athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.35.0-uclibc /bin/wget /usr/bin/wget

EXPOSE 8080
ENV ATHENAEUM_DB_ROOT=/home/nonroot/athenaeum

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
