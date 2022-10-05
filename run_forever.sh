RETRY_TIMEOUT=10

while true; do
  ./carmack
  echo "stopped, restarting in ${RETRY_TIMEOUT} seconds.."
  sleep ${RETRY_TIMEOUT}
done
