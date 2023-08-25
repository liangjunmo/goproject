FROM nginx:1.23.4-bullseye

COPY ./deploy/k8s/goproject/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./deploy/k8s/goproject/nginx/goproject.conf /etc/nginx/conf.d/goproject.conf
COPY ./deploy/k8s/goproject/nginx/*.html /etc/nginx/
