docker build -t apptimus-static-generator -f Dockerfile.static .

docker run --network apptimus_default --rm `
    -e DB_HOST=mysql_db `
    -e DB_USER=apptimus `
    -e DB_PASS=1q2w3e `
    -e DB_NAME=apptimus_db `
    apptimus-static-generator

docker build -t apptimus-static-site -f Dockerfile.static .

docker run -d --name apptimus-static-site -p 8081:80 apptimus-static-site