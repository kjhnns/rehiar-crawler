#!/bin/bash

rm -rf data logs
scp -P 21000 -r . rehiar@138.68.109.221:~/crawler/src
ssh -p 21000 rehiar@138.68.109.221 "cd ~/crawler/src/ && go build && mv crawler .."


echo "success."
