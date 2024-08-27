curl -isSX GET http://localhost:8000/urlshortener/analytics/000000000000 > test16.out 2>&1
diff test16.out test16.ref