/var/web-app-sitemap-generator /var/sitemap.xml
aws s3 cp /var/sitemap.xml s3://web.citymunchapp.com/sitemap.xml --acl public-read --cache-control 'max-age=86400'
