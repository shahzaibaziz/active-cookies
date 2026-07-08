FROM scratch

COPY bin/main /activecookie

ENTRYPOINT ["/activecookie"]
