#!usr/bin/bash
socat tcp4-listen:5870,fork socks4a:127.0.0.1:smtp.dizum.com:2525,socksport=9050
