
scp -P 21000 -r crawler rehiar@138.68.109.221:~/crawler/src
ssh -p 21000 rehiar@138.68.109.221 "cd ~/crawler/src/crawler/ && go build && mv crawler ../.."


echo "success."
say "done"
