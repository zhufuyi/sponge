### Start the elasticsearch service

Switch to the elasticsearch directory and start the service

> docker-compose up -d

<br>

### Start jaeger service

Switch to the jaeger directory, open the. env environment variable, and fill in the url, login account, and password of elastic search

Start jaeger

> docker-compose up -d

Check whether it is normal

> docker-compose ps
