CityMunch web app sitemap generator
===================================

Generates a sitemap.xml file for the web app.

3 environment variables must be set: `CM_API_ENDPOINT`, `CM_API_KEY`, `CM_WEB_APP_URL`.

Example usage:

```
CM_API_ENDPOINT="https://prod-api.citymunchapp.com" \
	CM_API_KEY="..." \
	CM_WEB_APP_URL="https://web.citymunchapp.com" \
	./sitemap-generator sitemap.xml
```

Development
===========

Build with `./build`
