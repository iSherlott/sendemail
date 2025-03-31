docker build -t sendemail .
docker run -d --name send-email --env-file .env -p 8080:8080 sendemail