
	FROM  gcc:latest
	WORKDIR /app
	COPY ../check.c .
	CMD ["bash", "-c", "gcc -g check.c -o check && ./check"]
	