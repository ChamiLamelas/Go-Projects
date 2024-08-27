curl -isSX GET http://localhost:8000/urlshortener/expand/000000000000 > test15.out 2>&1
diff test15.out test15.ref