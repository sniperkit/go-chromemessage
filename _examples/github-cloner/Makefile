.PHONY: $(ASSEMBLY) install install_firefox install_chromium install_chrome uninstall clean 

ASSEMBLY			= dist

EXTENSION_DOMAIN 	= snk.golang.chromemsg
EXTENSION_SRC_DIR	= $(CURDIR)/extension
EXTENSION_DIST_DIR 	= $(CURDIR)/extension/$(ASSEMBLY)

TARGET_FIREFOX		= ~/.mozilla/native-messaging-hosts/$(EXTENSION_DOMAIN).json
TARGET_CHROMIUM		= ~/.config/chromium/NativeMessagingHosts/$(EXTENSION_DOMAIN).json
TARGET_CHROME		= ~/.config/google-chrome/NativeMessagingHosts/$(EXTENSION_DOMAIN).json

.DEFAULT_GOAL: $(ASSEMBLY)

# @$(GOPATH)/bin/gopherjs build app.go -o ${EXTENSION_DIST_DIR}/js/app.js
$(ASSEMBLY):
	@mkdir -p ${EXTENSION_DIST_DIR}/js ${EXTENSION_DIST_DIR}/css ${EXTENSION_DIST_DIR}/fonts ${EXTENSION_DIST_DIR}/icons ${EXTENSION_DIST_DIR}/static
	@cp $(EXTENSION_SRC_DIR)/manifest.json ${EXTENSION_DIST_DIR}
	@cp $(EXTENSION_SRC_DIR)/index.html ${EXTENSION_DIST_DIR}/static
	@cp $(EXTENSION_SRC_DIR)/css/*.css ${EXTENSION_DIST_DIR}/css
	@cp $(EXTENSION_SRC_DIR)/js/*.js ${EXTENSION_DIST_DIR}/js
	@cp $(EXTENSION_SRC_DIR)/fonts/*.* ${EXTENSION_DIST_DIR}/fonts

# helpers
install_firefox:
	@mkdir -p $(dir $(TARGET_FIREFOX))
	@install $(EXTENSION_DOMAIN).json $(TARGET_FIREFOX)
	@sed -i 's|$$PWD|$(PWD)|' $(TARGET_FIREFOX)

install_chromium:
	@mkdir -p $(dir $(TARGET_CHROMIUM))
	@install $(EXTENSION_DOMAIN).json $(TARGET_CHROMIUM)
	@sed -i 's|$$PWD|$(PWD)|' $(TARGET_CHROMIUM)

install_chrome:
	@mkdir -p $(dir $(TARGET_CHROME))
	@install $(CURDIR)/$(EXTENSION_DOMAIN).json $(TARGET_CHROME)
	@sed -i 's|$$PWD|$(PWD)|' $(TARGET_CHROME)

install: install_firefox install_chromium install_chrome

clean:
	@if [ -d ${EXTENSION_DIST_DIR} ] ; then rm -rf ${EXTENSION_DIST_DIR} ; fi
# @git clean -dfx

# install:
# 	@mkdir -p $(NATIVE_PATH) $(FIREFOX_PATH) $(FFEXTENSION_PATH) $(EXTERNAL_PATH)
# 	@cp -v $(EXTENSION_DOMAIN).json $(NATIVE_PATH)
# 	@cp -v ../ckjefchnfjhjfedoccjbhjpbncimppeg.json $(EXTERNAL_PATH)
# cp -rv build/Release/snk.golang.chromemsg.app $(NATIVE_PATH)
# cp -v snk.golang.chromemsg.firefox.json $(FIREFOX_PATH)/snk.golang.chromemsg.json
# cp -v ../{443830f0-1fff-4f9a-aa1e-444bafbc7319}.xpi $(FFEXTENSION_PATH)

uninstall:
	@rm -f $(EXTERNAL_PATH)/ckjefchnfjhjfedoccjbhjpbncimppeg.json
	@rm -f $(NATIVE_PATH)/snk.golang.chromemsg.json
	@rm -rf $(NATIVE_PATH)/snk.golang.chromemsg.app

# rm -f $(FIREFOX_PATH)/snk.golang.chromemsg.json
# rm -f $(FFEXTENSION_PATH)/{443830f0-1fff-4f9a-aa1e-444bafbc7319}.xpi