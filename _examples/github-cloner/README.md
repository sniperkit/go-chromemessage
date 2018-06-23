# Quick Repo Cloner

Instantly `git clone` a local copy of the GitHub/Bitbucket repository you're looking at in Chrome.

Requires very recent Chrome (45ish) at the moment because of gratuitous ES6 that I haven't setup a transpilation pipeline for yet.


## How it works

Clicking button runs `extension/main.js` which injects `extension/contentscript.js` into the page. The contentscript looks at the tab's URL/location and generates the path to its `.git` repository, if possible. Contentscript sends a message (`git clone <repo.git> repo`) back to `main.js`, which relays it to `snk.golang.chromemsg`, a Python script, using Chrome's "native host" messaging. All that the python script does is execute the command sent from the contentscript.


## Installation

**Dependencies:** Git, Golang, (Optional: Python).

1. Install the Chrome native host to relay commands to git:
    ```shell
    $ cd $GOPATH/src/github.com/sniperkit/snk.golang.chromemsg
    $ ./shared/scripts/install-native-host.sh
    ```

2. Install the Chrome extension.

3. Specify the directory to clone repos into (eg. `$HOME/dev/projects` or `$GOPATH/src`) using the Options panel of the extension, accessible from Chrome's Extensions page.


## Todo

- The Options window is slow af to load despite being absurdly simple. This seems to be a known thing with Chrome's options v2 api. I should rewrite with the old options api.

- You may need to edit the `allowed_origins` in `snk.golang.chromemsg.json` with the extension ID shown after you install the extension. This ID can change during development until the extension gets packaged for the Chrome Web Store.

- Success/error of the actual `git clone` is not visible right now


## Licence (MIT)
