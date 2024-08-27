curl -isSX GET http://localhost:8000/urlshortener/expand/google >> test19.out 2>&1
diff test19.out test19.ref