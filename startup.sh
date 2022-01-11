#!/bin/bash

apk add curl

curl http://sg.wildfire.paloaltonetworks.com/publicapi/test/elf -o malware-file-sg
curl http://jp.wildfire.paloaltonetworks.com/publicapi/test/elf -o malware-file-jp
curl http://eu.wildfire.paloaltonetworks.com/publicapi/test/elf -o malware-file-eu
curl http://wildfire.paloaltonetworks.com/publicapi/test/elf -o malware-file

./web-req-info
