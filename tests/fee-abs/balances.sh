
BALANCES="$(docker compose exec osmosis osmosisd query bank balances osmo12a47mwn4sf3v7qsnn7l65dvjr47pmx6k8gsrnz --output json 2>&1)"

echo "aaa $BALANCES"
