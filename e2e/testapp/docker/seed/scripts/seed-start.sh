if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi
if [ -z "$LOGLEVEL" ]; then
    LOGLEVEL="info"
fi

./bin/polard start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR" --polaris.execution-client.jwt-secret-path "$JWTSECRETPATH" --polaris.execution-client.rpc-dial-url "$DIALURL"
