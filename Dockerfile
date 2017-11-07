FROM ubuntu:16.04

RUN apt update && apt upgrade -y && apt install cron ca-certificates awscli -y -qq

ADD dist/web-app-sitemap-generator-linux /var/web-app-sitemap-generator

RUN echo '0 */12 * * * root /var/cron.sh' >> /etc/crontab
ADD cron.sh /var/cron.sh
RUN chmod +x /var/cron.sh

# @todo - figure out how to get these environment variables into the ECS Docker.
CM_API_ENDPOINT="https://prod-api.citymunchapp.com"
CM_API_KEY="2v4c6s2iolua50k3o1cn6b87c8"
CM_WEB_APP_URL="https://web.citymunchapp.com"

# ADD crontab /etc/cron.d/sitemap-crontab
# RUN chmod 0644 /etc/cron.d/sitemap-crontab
# RUN touch /var/log/cron.log

ENTRYPOINT ["cron", "-f"]
