FROM redis:8.6.3-alpine3.23

LABEL authors="Redis DB special configuration by Cactus Maestro (On Kube)"

# Create required directories
RUN mkdir -p /var/lib/redis /var/log/redis && \
    chown -R redis:redis /var/lib/redis /var/log/redis

# Copy config
COPY redis.conf /usr/local/etc/redis/redis.conf

# Use Redis user for security
USER redis

# Expose Redis port
EXPOSE 6379

# Start Redis with custom config
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]