# webwol
simple Wake-On-Clan client configurable via a web page

## Environment variables

### ENV_PORT 
The port to listen. Default is ```8080```

### ENV_ASSETS
The directory where the assets (js, css, image) are. Default is ```assets```

### ENV_TEMPLATES
The directory where the templates are. Default is ```templates```
### ENV_CONFIG
The directory where the configuration is stored. Default is ```config```

### ENV_BASEURL
The base url of the app. Default is ```http://localhost:8080```. This is needed if QR-Codes are used.

## docker compose

```
version: '2'

services:
  webwol:
    image: 'darthbermel/webwol'
    ports:
      - "8080:8080"
    environment:
      ENV_CONFIG: /config
    restart: always
    volumes:
      - ./config:/config
```