FROM nginx:1.23.4-bullseye

COPY ./k8s/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./k8s/nginx/goproject.conf /etc/nginx/conf.d/goproject.conf
COPY ./k8s/nginx/*.html /etc/nginx/
