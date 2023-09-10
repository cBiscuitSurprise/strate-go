export REDIS_CONN_DIR="test/sg_local/etc/redis-conn"
HOST='0.0.0.0'
PORT='12345'

./bazel-bin/strate-go_/strate-go serve --origin $HOST --port $PORT
