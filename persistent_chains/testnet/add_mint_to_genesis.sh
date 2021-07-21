#!/bin/bash

echo "##### [Validator operators] Add them to the genesis" 

cheqd-noded add-genesis-account ColdStorage 510000000cheq
cheqd-noded add-genesis-account Treasury 290000000cheq
cheqd-noded add-genesis-account operator1 50000000cheq
cheqd-noded add-genesis-account operator2 20000000cheq
cheqd-noded add-genesis-account operator3 40000000cheq
cheqd-noded add-genesis-account operator4 90000000cheq

