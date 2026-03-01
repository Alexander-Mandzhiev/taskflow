FROM traefik:v3.2

# Копируем конфигурационные файлы в образ
COPY traefik.yml /etc/traefik/traefik.yml
COPY dynamic.yml /etc/traefik/dynamic.yml
COPY traefik.prod.yml /etc/traefik/traefik.prod.yml
COPY dynamic.prod.yml /etc/traefik/dynamic.prod.yml

