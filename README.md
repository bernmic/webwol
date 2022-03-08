# webwol
simple Wake-On-Clan client configurable via a web page

## Environment variables

### WEBWOL_PORT 
The port to listen. Default is ```8080```

### WEBWOL_ASSETS
The directory where the assets (js, css, image) are. Default is ```assets```

### WEBWOL_TEMPLATES
The directory where the templates are. Default is ```templates```
### WEBWOL_CONFIG
The directory where the configuration is stored. Default is ```config```

### WEBWOL_BASEURL
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
      WEBWOL_CONFIG: /config
    restart: always
    volumes:
      - ./config:/config
```