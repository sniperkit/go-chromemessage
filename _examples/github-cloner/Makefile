.PHONY: clean

SRC			= 	$(CURDIR)/extension
DIST 		= 	$(CURDIR)/extension/dist
ASSEMBLY	= dist

.DEFAULT_GOAL: $(ASSEMBLY)

# @$(GOPATH)/bin/gopherjs build app.go -o ${DIST}/js/app.js
$(ASSEMBLY):
	@mkdir -p ${DIST}/js ${DIST}/css ${DIST}/fonts ${DIST}/icons ${DIST}/static
	@cp $(SRC)/manifest.json ${DIST}
	@cp $(SRC)/index.html ${DIST}/static
	@cp $(SRC)/css/*.css ${DIST}/css
	@cp $(SRC)/js/*.js ${DIST}/js
	@cp $(SRC)/fonts/*.* ${DIST}/fonts

clean:
	@if [ -d ${DIST} ] ; then rm -rf ${DIST} ; fi