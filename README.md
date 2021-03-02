# Twedit

Post image from randomly-selected sub-reddit

# Build:
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -v -i -o twedit cmd/main.go
```


# Config

copy `CREDENTIALS.sample` to `CREDENTIALS` and replace each line with coresponding value from your twitter developer page.

# License

    Copyright © 2000 Widnyana Putra <wid (a) widnyana.web.id>
    This work is free. You can redistribute it and/or modify it under the
    terms of the Do What The Fuck You Want To Public License, Version 2,
    as published by Sam Hocevar. See the COPYING file for more details.