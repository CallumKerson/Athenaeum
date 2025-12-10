FROM gcr.io/distroless/static-debian11:debug-nonroot@sha256:55716e80a7d4320ce9bc2dc8636fc193b418638041b817cf3306696bd0f975d1 AS production

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.37.0-uclibc@sha256:e58014df10240c35c7b1df7ff8e859ad6a54d061bde77c249f96880e15d83049 /bin/wget /usr/bin/wget

EXPOSE 8080
ENV ATHENAEUM_DB_ROOT=/home/nonroot/athenaeum

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
