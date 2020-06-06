FROM alpine:3.7

COPY api/api /api
COPY ui/build /ui/build

RUN chmod +x /api

CMD ["/api"]