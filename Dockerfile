# This is here for deploying this bot on Raspberry Pi (ARMv7)
# Use images with different architectures if you want it on other platforms

FROM cellofellow/rpi-arch

MAINTAINER Kanedias (kanedias@xaker.ru)

# this will also create these dirs if nonexistent
ENV GOPATH=/tldr-bot-git
WORKDIR $GOPATH

# prerequisites - golang and git
RUN pacman --noconfirm -Sy go git

# main build script
RUN git clone https://github.com/Adonai/TldrBot.git $GOPATH \
 && go get -v "github.com/JesusIslam/tldr" \
 && go get -v "github.com/advancedlogic/GoOse"

# bot requires token to work properly
COPY bot-token.txt $GOPATH

ENTRYPOINT ["go", "run", "src/com/tldr/bot/bot.go"]
