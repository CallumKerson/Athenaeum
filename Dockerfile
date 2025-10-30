FROM gcr.io/distroless/static-debian11:debug-nonroot@sha256:55716e80a7d4320ce9bc2dc8636fc193b418638041b817cf3306696bd0f975d1 AS production

COPY athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.37.0-uclibc@sha256:7b5ade1b11f52517f5cdf3b6194588587b9c4b14d3244b6c96ed7f5b0ff31458 /bin/wget /usr/bin/wget

EXPOSE 8080
ENV ATHENAEUM_DB_ROOT=/home/nonroot/athenaeum

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
