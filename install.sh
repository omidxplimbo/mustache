#!/bin/bash

# global stuff
PROJECT_PATH="$HOME/.mustache-project"
BASE_PATH="$HOME/mustache"
WORDLISTS_PATH="$BASE_PATH/wordlists"
RESOLVERS_PATH="$BASE_PATH/resolvers"
CONFIGS_PATH="$BASE_PATH/configs"
LOGS_PATH="$BASE_PATH/logs"
DEFAULT_SHELL="$HOME/.bashrc"
CWD=$(pwd)
PACKGE_MANAGER="apt-get"

SUDO="sudo"
if [ "$(whoami)" == "root" ]; then
    SUDO=""
fi
[ -x "$(command -v apt)" ] && PACKGE_MANAGER="apt"
if [ -f "$HOME/.zshrc" ]; then
    DEFAULT_SHELL="$HOME/.zshrc"
fi

announce() {
    echo -e "\033[1;37m[\033[1;31m+\033[1;37m]\033[1;32m $1 \033[0m"
}

announce2() {
    echo -e "\033[1;37m[\033[1;32m+\033[1;37m]\033[1;36m $1 \033[0m"
}

install_banner() {
    echo -e "\033[1;37m[\033[1;34m+\033[1;37m]\033[1;32m Installing $1 \033[0m"
}

download() {
    wget --no-check-certificate -q -O $1 $2
    if [ ! -f "$1" ]; then
        wget --no-check-certificate -q -O $1 $2
    fi
}

extractZip() {
	unzip -q -o -j $1 -d $BINARIES_PATH/
	rm -rf $1
}

extractGz() {
	tar -xf $1 -C $BINARIES_PATH/
	rm -rf $1
}

announce "Downloading and Setting up Mustache package. This might take couple minutes if you're using weak connection"
announce2 "This suppose to be private. But if you see this, Fuck You <3 \033[0m"
announce2 "Please don't share this link with other people since this installation is private \033[0m"

announce "NOTE that this installation only works on\033[0m Linux 64-bit intel based machine machine."

$SUDO $PACKGE_MANAGER update -qq > /dev/null 2>&1

announce "Checking Essential tools"

# reinstall all essioontials tools just to double check
[ -x "$(command -v wget)" ] || $SUDO $PACKGE_MANAGER -qq install wget -y >/dev/null 2>&1
[ -x "$(command -v curl)" ] || $SUDO $PACKGE_MANAGER -qq install curl -y >/dev/null 2>&1
[ -x "$(command -v tmux)" ] || $SUDO $PACKGE_MANAGER -qq install tmux -y >/dev/null 2>&1
[ -x "$(command -v git)" ] || $SUDO $PACKGE_MANAGER -qq install git -y >/dev/null 2>&1
[ -x "$(command -v nmap)" ] || $SUDO $PACKGE_MANAGER -qq install nmap -y >/dev/null 2>&1
[ -x "$(command -v masscan)" ] || $SUDO $PACKGE_MANAGER -qq install masscan -y >/dev/null 2>&1
[ -x "$(command -v make)" ] || $SUDO $PACKGE_MANAGER -qq install build-essential -y >/dev/null 2>&1
[ -x "$(command -v unzip)" ] || $SUDO $PACKGE_MANAGER -qq install unzip -y >/dev/null 2>&1
[ -x "$(command -v chromium)" ] || $SUDO $PACKGE_MANAGER -qq install chromium -y >/dev/null 2>&1
[ -x "$(command -v chromium-browser)" ] || $SUDO $PACKGE_MANAGER -qq install chromium-browser -y >/dev/null 2>&1
[ -x "$(command -v jq)" ] || $SUDO $PACKGE_MANAGER -qq install jq -y >/dev/null 2>&1
[ -x "$(command -v make)" ] || $SUDO $PACKGE_MANAGER -qq install build-essential -y >/dev/null 2>&1
[ -x "$(command -v rsync)" ] || $SUDO $PACKGE_MANAGER -qq install rsync -y >/dev/null 2>&1
[ -x "$(command -v netstat)" ] || $SUDO $PACKGE_MANAGER -qq install coreutils net-tools -y >/dev/null 2>&1
[ -x "$(command -v htop)" ] || $SUDO $PACKGE_MANAGER -qq install htop -y >/dev/null 2>&1
[ -x "$(command -v timeout)" ] || $SUDO $PACKGE_MANAGER install timeout -y >/dev/null 2>&1
[ -x "$(command -v pip)" ] || $SUDO $PACKGE_MANAGER install python3 python3-pip -y >/dev/null 2>&1


announce "Create essential folders"
mkdir -p $BASE_PATH
mkdir -p $BASE_PATH/$BINARIES_PATH
mkdir -p $BASE_PATH/$CONFIGS_PATH
mkdir -p $BASE_PATH/$WORDLISTS_PATH
mkdir -p $BASE_PATH/$RESOLVERS_PATH
mkdir -p $BASE_PATH/$LOGS_PATH


announce "Install last version of golang"
wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
sudo rm https://go.dev/dl/go1.20.3.linux-amd64.tar.gz

announce "Install binaries"
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
go install -v github.com/owasp-amass/amass/v3/...@master
go install github.com/projectdiscovery/alterx/cmd/alterx@latest
go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest
GO111MODULE=on go install github.com/jaeles-project/gospider@latest
go install github.com/ffuf/ffuf/v2@latest
go install -v github.com/tomnomnom/anew@latest
go install github.com/lc/gau/v2/cmd/gau@latest
go install -v github.com/OJ/gobuster/v3@latest
go install github.com/gwen001/github-endpoints@latest
go install github.com/projectdiscovery/katana/cmd/katana@latest
go install github.com/tomnomnom/waybackurls@latest
go install -v github.com/tomnomnom/anew@latest
pip3 install py-altdns==1.0.2
go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
go install -v github.com/d3mondev/puredns/v2@latest
go install -v github.com/projectdiscovery/mapcidr/cmd/mapcidr@latest
git clone https://github.com/ProjectAnte/dnsgen && cd dnsgen && pip3 install -r requirements.txt && python3 setup.py install
git clone https://github.com/devanshbatham/ParamSpider && cd ParamSpider && sudo pip3 install -r requirements.txt
git clone https://github.com/Sh1Yo/x8 && cd x8 && cargo build --release && cd target/release && sudo cp x8 /usr/local/bin && cd
sudo pip3 install arjun
go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest
go install -v github.com/hahwul/dalfox/v2@latest


announce "Install last version of mustache package"
go install -v https://github.com/omidxplimbo/mustache@latest
cd $HOME/go/bin/ && sudo cp * /usr/local/bin/ && cd



echo "---->>>"
mustache -h
echo "---->>>"

announce "The Mustache installation is done..."