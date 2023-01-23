LATEST_VERSION=$(git ls-remote --tags git@github.com:StirlingMarketingGroup/ts-smg.git | cut -f 2 | cut -d / -f 3 | php -r '$a=explode("\n",trim(file_get_contents("php://stdin")));usort($a,"version_compare");echo array_reverse($a)[0];')

npm i --save-dev "ssh://git@github.com/StirlingMarketingGroup/ts-smg.git#${LATEST_VERSION}"