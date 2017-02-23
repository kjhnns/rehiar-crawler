#!/bin/bash

# ssh -p 21000 rehiar@138.68.109.221 "cd ~/crawler/src/ && go get && go build && mv models.csv crawler .."
# ssh -p 21000 rehiar@138.68.109.221 "cd ~/crawler/src/ && go get -u all && go build && mv models.csv crawler .."


ssh -p 21000 rehiar@138.68.109.221 "export GOPATH=~/workspace/ && rm -rf ~/workspace/src/github.com/kjhnns/rehiar-crawler && go get github.com/kjhnns/rehiar-crawler && go build github.com/kjhnns/rehiar-crawler && cd ~ && rm ~/workspace/bin/rehiar-crawler && mv rehiar-crawler ~/workspace/src/github.com/kjhnns/rehiar-crawler/models.csv ~/crawler/ "

# cd ~/workspace/src/github.com/kjhnns/rehiar-crawler &&
#  git pull &&
#  go get
#  go build &&
#  cd ~ &&
#  rm ~/workspace/bin/rehiar-crawler &&
#  mv ~/workspace/bin/rehiar-crawler ~/workspace/src/github.com/kjhnns/rehiar-crawler/models.csv ~/crawler/ "


# go get -u all &&
# go build github.com/kjhnns/rehiar-crawler &&

echo "success."
